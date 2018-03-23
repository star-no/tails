package config

import (
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Input     string            `yaml:"pattern"`
	NumberMap map[int]string    `yaml:"numberMap"`
	WordMap   map[string]string `yaml:"wordMap"`
}

var ConfigList = []Config{}
var RootPath string

func ReadConfigFile() {
	RootPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))

	f, err := os.Open(RootPath + "/config.yaml")
	if err != nil {
		log.Fatal("can't find file 'config.yaml'\nConfig file must exist at " + RootPath + "\n")
	}

	s := readConfigFile(f)

	err = yaml.Unmarshal([]byte(s), &ConfigList)
	if err != nil {
		log.Fatalf("There is error in config.yaml.\nplease check file contents.")
	}
}

func readConfigFile(f *os.File) string {
	s := ""
	n := 0
	b := make([]byte, 1024)

	for {
		c, err := f.ReadAt(b, int64(n))

		if err != io.EOF && err != nil {
			panic(err)
		}
		s += string(b[:c])
		n += c

		if err == io.EOF || c == 0 {
			break
		}
	}

	return s
}
