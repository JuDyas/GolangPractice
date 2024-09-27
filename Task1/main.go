package main

import (
	"fmt"
	"unicode"
)

func main() {
	inputText := "ыjj(線l):L線, J-j/Ы"
	fmt.Println(palindrom(inputText))
}

func palindrom(text string) bool {
	var cleanText []rune

	for _, char := range text {
		if unicode.IsLetter(char) {
			lowerChar := unicode.ToLower(char)
			cleanText = append(cleanText, lowerChar)
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
