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
	enc := vyzhenercipher.Encode(text, key)
	dec := vyzhenercipher.Decode(enc, key)
	fmt.Println(enc)
	fmt.Println(dec)
}
