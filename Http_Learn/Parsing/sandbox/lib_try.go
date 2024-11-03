package main

import (
	"GolangPractice/Http_Learn/Parsing/vyzhenercipher"
	"fmt"
)

func main() {
	var (
		key  = "БаNaн"
		text = "!Привет 123 World!"
	)
	enc := encryptTry(key, text)
	dec := decryptTry(key, enc)
	fmt.Println(enc)
	fmt.Println(dec)
}

// encryptTry - Try lib for encrypting any text
func encryptTry(key, text string) string {
	chipher := &vyzhenercipher.VizhenerCipher{
		Key:  key,
		Text: text,
	}

	chipher.Encrypt()
	return chipher.ChangedText
}

func decryptTry(key, text string) string {
	chipher := &vyzhenercipher.VizhenerCipher{
		Key:  key,
		Text: text,
	}
	chipher.Decrypt()
	return chipher.ChangedText
}
