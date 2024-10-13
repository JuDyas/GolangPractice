/*
Написать функцию которая зашифрует (md5, sha1, etc.) строки в JSON-не по ключам, или пути (https://ru.wikipedia.org/wiki/JSON)
*/

package main

import (
	"GolangPractice/Task7/lib"
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"log"
)

type HashSha1 struct{}

func (h *HashSha1) Hash(data string) (string, error) {
	hash := sha1.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil)), nil
}

var (
	jsonData = flag.String("json", `{"name": "Denys", "age": "123", "password": "UpbQ2X*&FQ$L", "info": {"weight": "60", "ip": "192.168.1.1"}}`, "Input JSON data")
	whatHash = flag.String("hashElement", "password,ip", "Elements to hash")
)

func main() {
	flag.Parse()
	dataProcessed, whatHashProcessed, err := lib.ProcessFlags(jsonData, whatHash)
	if err != nil {
		log.SetPrefix("ERROR: ")
		log.Fatalf("Processing data: %v", err)
	}
	resultJson, err := lib.HashJson(dataProcessed, whatHashProcessed, &HashSha1{})

	if err != nil {
		log.SetPrefix("ERROR: ")
		log.Fatalf("During hashing: %v", err)
	}
	log.Println(resultJson)
}
