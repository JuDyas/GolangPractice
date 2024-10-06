package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func main() {
	fileNameIn := "./example.json"
	fileNameOut := "exampleSHA1.json"
	whatHash := []string{"password", "ip"}

	err := readWrite(fileNameIn, fileNameOut, whatHash)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

// readWrite - open/read/write json file
func readWrite(fileNameIn, fileNameOut string, whatHash []string) error {
	file, err := os.Open(fileNameIn)
	if err != nil {
		return err
	}
	defer file.Close()

	readerByte, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var data map[string]interface{}
	err = json.Unmarshal(readerByte, &data)
	if err != nil {
		return err
	}
	hashJson(data, whatHash)

	writeJson, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(fileNameOut, writeJson, 0644)
	if err != nil {
		return err
	}

	fmt.Println("Complete")
	return nil
}

// sha1Hash - protocol of hashing
func sha1Hash(data string) string {
	hash := sha1.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

// hashJson - do hashing
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

// exist - check existing item in json file
func exist(whatHash []string, k string) bool {
	for _, v := range whatHash {
		if v == k {
			return true
		}
	}
	return false
}
