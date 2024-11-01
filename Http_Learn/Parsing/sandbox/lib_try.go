package main

import (
	"GolangPractice/Http_Learn/Parsing/vyzhenercipher"
	"log"
)

func main() {
	var (
		key  = "БаNaн"
		text = "!Привет 123 World!"
	)
	encryptTry(key, text)
}

// encryptTry - Try lib for encrypting any text
func encryptTry(key, text string) {
	chipher := &vyzhenercipher.VizhenerCipher{
		Key:  key,
		Text: text,
	}

	chipher.Encrypt()
	log.Println(chipher.ChangedText)
}
