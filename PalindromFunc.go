package main

import (
	"fmt"
	"strings"
)

func main() {
	inputText := "Fjj(l):L, J-j/f"
	fmt.Println(Palindrom(inputText))
}

func Palindrom(stringIn string) bool {

	stringIn = strings.ToLower(stringIn)
	cleanString := ""

	for _, char := range stringIn {
		if char >= 'a' && char <= 'z' {
			cleanString += string(char)
		}
	}

	length := len(cleanString)
	for i := 0; i < length/2; i++ {
		if cleanString[i] != cleanString[length-i-1] {
			return false
		}
	}
	return true
}
