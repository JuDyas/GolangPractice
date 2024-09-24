package main

import (
	"fmt"
	"strings"
)

func main() {
	inputText := "Fjj(l):L, J-j/f"
	fmt.Println(palindrom(inputText))
}

func palindrom(text string) bool {
	text = strings.ToLower(text)
	cleanText := ""

	for _, char := range text {
		if char >= 'a' && char <= 'z' {
			cleanText += string(char)
		}
	}

	length := len(cleanText)
	for i := 0; i < length/2; i++ {
		if cleanText[i] != cleanText[length-i-1] {
			return false
		}
	}
	return true
}
