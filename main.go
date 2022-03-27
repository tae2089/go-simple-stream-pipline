package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func PrintData(data string) <-chan string {
	oc := make(chan string)
	go func() {
		oc <- data
		close(oc)
	}()
	return oc
}

func ExtractRegexpData(ic <-chan string) <-chan string {
	oc := make(chan string)

	go func() {
		for data := range ic {
			re1, err := regexp.Compile("\\w+/\\w+/\\w+ \\w+:\\w+:\\w+ - ")
			if err != nil {
				panic(err)
			}
			regexpRemoveText := re1.FindString(data)
			//fmt.Println(regexpRemoveText)
			//fmt.Println(data)
			extraText := strings.Replace(data, regexpRemoveText, "", 1)
			//fmt.Println(extraText)
			oc <- extraText
		}
		close(oc)
	}()
	return oc
}

func WriteFile(ic <-chan string) <-chan string {
	oc := make(chan string)
	var f *os.File
	var err error
	f, err = os.OpenFile("data.txt", os.O_RDWR|os.O_CREATE, 0755)
	f.Seek(0, 2)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	for data := range ic {
		dataByte := []byte(data + "\n")
		_, err = f.Write(dataByte)
		if err != nil {
			panic(err)
		}
	}
	return oc
}

func PrintChanData(ic <-chan string) {
	for data := range ic {
		fmt.Println(data)
	}
}

func main() {
	str1 := "2022/01/03 09:53:27 - 123123"
	PrintChanData(ExtractRegexpData(PrintData(str1)))
	WriteFile(ExtractRegexpData(PrintData(str1)))
	context.Context()
}
