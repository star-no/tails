package main

import (
	"fmt"
	"flag"
	"os"
	"time"
	"bufio"
	. "github.com/logrusorgru/aurora"
	"strings"
	"path/filepath"
	"log"
	"regexp"
)

var file = flag.String("f", "", "file-name")
var debug = flag.Bool("d", false, "debug")

func main() {
	flag.Parse()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	filePath := *file
	if !strings.Contains(filePath, dir) {
		filePath = dir + "/" + filePath
	}

	if *debug {
		fmt.Println("[DEBUG]", dir)
		fmt.Println("[DEBUG]", *file)
		fmt.Println("[DEBUG]", filePath)
	}

	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var currentSeek int64 = 0
	r := regexp.MustCompile("\\d{4}-\\d{2}-\\d{2}.*")

	for {
		f.Seek(currentSeek, 0)

		reader := bufio.NewReader(f)
		bytes, _ := reader.ReadBytes('\n')

		if len(bytes) != 0 {
			currentSeek += int64(len(bytes))
			if !r.Match(bytes) {
				fmt.Println(string(bytes))
				continue
			}

			buf := []interface{}{}
			line := string(bytes)

			words := strings.Fields(line)

			for idx, s := range words {
				switch idx {
				case 2:
					buf = append(buf, Green(s))
				case 3:
					buf = append(buf, Magenta(s))
				case 6:
					buf = append(buf, Cyan(s))
				default:
					buf = append(buf, s)
				}
			}
			fmt.Println(strings.Trim(fmt.Sprint(buf), "[]"))

		} else {
			time.Sleep(time.Millisecond * 50)
		}
	}
}
