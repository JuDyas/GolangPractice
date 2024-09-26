package main

import (
	"fmt"
	"unicode"
)

func main() {
	textTry := "I`ll be back"
	shiftTry := 3
	fmt.Println("Оригинал: ", textTry, " Сдвиг: ", cezar(textTry, shiftTry))

}

func cezar(text string, shift int) string {
	var changedText []rune
	shift = shift % 26 // если 27, то будет 1, чтобы не вылезать за алфавит :)

	for _, char := range text {
		if unicode.IsLetter(char) {
			changedText = append(changedText, shiftText(char, shift))
		} else {
			changedText = append(changedText, char) // Если символ, то соотвественно не меняем
		}
	}
	return string(changedText)
}

func shiftText(char rune, shift int) rune { // Фунция для сдвига буковок (увы, только латиница)
	if unicode.IsLower(char) {
		return 'a' + (char-'a'+rune(shift))%26 // Принцып работы тут - https://drive.google.com/file/d/15CVPX2bywOTffOlKTZ_voQqlNN-8X2VR/view?usp=sharing
	} else if unicode.IsUpper(char) {
		return 'A' + (char-'A'+rune(shift))%26
	}
	return char
}
