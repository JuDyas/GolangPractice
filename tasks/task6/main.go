package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

func main() {
	fileName := flag.String("file", "", "Path to the file")
	inputString := flag.String("string", "", "String to count characters")
	flag.Parse()

	var reader io.Reader

	if *inputString != "" {
		reader = strings.NewReader(*inputString)
	} else if *fileName != "" {
		file, err := os.Open(*fileName)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				fmt.Println("Error closing file:", err)
			}
		}(file)
		reader = file
	} else {
		fmt.Println("You must provide either a file or a string input.")
		return
	}

	count, err := countCharacters(reader)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Total characters: %d\n", count)
	}
}

func countCharacters(reader io.Reader) (int, error) {
	bufReader := bufio.NewReader(reader)
	count := 0

	for {
		line, err := bufReader.ReadString('\n')

		if err != nil {
			if err == io.EOF {
				line = strings.TrimSpace(line)
				count += utf8.RuneCountInString(line)
				break
			}
			return 0, err
		}
		line = strings.TrimSpace(line)
		count += utf8.RuneCountInString(line)
	}
	return count, nil
}
