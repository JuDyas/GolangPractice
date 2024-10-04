package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func main() {
	fileName := "./example.json"
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	readerByte, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	var data map[string]interface{}
	err = json.Unmarshal(readerByte, &data)
	if err != nil {
		fmt.Println(err)
		return
	}

	whatHash := []string{"password", "ip"}
	hashJson(data, whatHash)

	writeJson, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = os.WriteFile("exampleSHA1.json", writeJson, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Гатова!")
}

// Функции для расчёта хеша и преобразования в строку.
func md5Hash(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}
func sha1Hash(data string) string {
	hash := sha1.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

// Проверка ключа и сздание хеша
func hashJson(data map[string]interface{}, whatHash []string) {
	for k, v := range data {
		strV, ok := v.(string)
		if ok && exist(whatHash, k) {
			data[k] = sha1Hash(strV)
		}
		if nk, ok := v.(map[string]interface{}); ok {
			hashJson(nk, whatHash)
		}
	}
}

// Проверка существования ключа
func exist(whatHash []string, k string) bool {
	for _, v := range whatHash {
		if v == k {
			return true
		}
	}
	return false
}
