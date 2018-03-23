package main

import (
	"bufio"
	"flag"
	"fmt"
	. "github.com/logrusorgru/aurora"
	"os"
	"regexp"
	"strings"
	"tails/reader/config"
	"time"
	"unicode/utf8"
)

var file = flag.String("f", "", "file-name")
var debug = flag.Bool("d", false, "debug")

func main() {
	flag.Parse()

	config.ReadConfigFile()

	filepath := parseFilepath()
	if *debug {
		fmt.Println("[DEBUG]", filepath)
	}

	f := openFile(filepath)
	defer f.Close()

	printTargetFile(f)
}

func parseFilepath() string {
	if *debug {
		fmt.Println("[DEBUG]", config.RootPath)
		fmt.Println("[DEBUG]", *file)
	}

	inputFilepath := *file

	if strings.HasPrefix(inputFilepath, "~/") {
		return inputFilepath
	}

	if strings.HasPrefix(inputFilepath, "/") {
		return inputFilepath
	}

	return config.RootPath + "/" + inputFilepath
}

func openFile(fileName string) *os.File {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	return f
}

func printTargetFile(f *os.File) {
	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}

	var currentSeek = fi.Size() - 1024

	for {
		f.Seek(currentSeek, 0)

		reader := bufio.NewReader(f)
		bytes, _ := reader.ReadBytes('\n')

		if len(bytes) != 0 {
			currentSeek += int64(len(bytes))

			for _, c := range config.ConfigList {
				r := regexp.MustCompile(c.Input)
				if r.Match(bytes) {
					printColor(c, bytes)
					goto End
				}
			}

			fmt.Print(string(bytes))
		End:
		} else {
			time.Sleep(time.Millisecond * 50)
		}
	}
}

func printColor(config config.Config, bytes []byte) {
	buf := []interface{}{}
	line := string(bytes)
	words := strings.Fields(line)

	for idx, str := range words {
		var color = "non-color"

		c, prs := config.NumberMap[idx]
		if prs {
			color = c
		}

		c, prs = config.WordMap[str]
		if prs {
			color = c
		}

		var coloredString interface{}
		switch color {
		case "non-color":
			coloredString = str
		case "black":
			coloredString = Black(str)
		case "green":
			coloredString = Green(str)
		case "gray":
			coloredString = Gray(str)
		case "magenta":
			coloredString = Magenta(str)
		case "red":
			coloredString = Red(str)
		case "blue":
			coloredString = Blue(str)
		case "cyan":
			coloredString = Cyan(str)
		}

		buf = append(buf, coloredString)

		if idx < len(words)-1 {
			buf = printSpace(line, str, buf)
		}
	}

	fmt.Println(strings.Trim(fmt.Sprint(buf), "[]"))
}

func printSpace(line string, str string, buf []interface{}) []interface{} {
	start := strings.Index(line, str) + utf8.RuneCountInString(str)
	var space = ""
	for i := start + 1; ' ' == line[i]; i++ {
		space += " "
	}

	if space != "" {
		buf = append(buf, space)
	}
	return buf
}
