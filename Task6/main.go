package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

func main() {
	fileName := "./example.txt"
	count, err := countCharacters(fileName)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Total characters in file: %d\n", count)
	}
}

// Подсчёт символов в файле
func countCharacters(fileName string) (int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	count := 0

	for {
		line, err := reader.ReadString('\n')

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
