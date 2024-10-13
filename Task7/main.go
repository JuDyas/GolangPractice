package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"
)

func main() {
	jsonData := flag.String("json", `{"name": "Denys", "age": "123", "password": "UpbQ2X*&FQ$L", "info": {"weight": "60", "ip": "192.168.1.1"}}`, "Input JSON data")
	hashType := flag.String("hash", "sha1", "Hash Type")
	whatHash := flag.String("hashElement", "password,ip", "Elements to hash")
	flag.Parse()

	dataProcessed, whatHashProcessed, err := processFlags(jsonData, whatHash)
	if err != nil {
		log.SetPrefix("ERROR: ")
		log.Fatalf("Processing data: %v", err)
	}
	resultJson, err := hashJson(dataProcessed, whatHashProcessed, hashType)
	if err != nil {
		log.SetPrefix("ERROR: ")
		log.Fatalf("During hashing: %v", err)
	}
	log.Println(resultJson)
}

// processFlags - processing input data to map and slice
func processFlags(jsonData, hashElements *string) (map[string]interface{}, []string, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(*jsonData), &data)
	if err != nil {
		return nil, nil, err
	}
	whatHash := strings.Split(*hashElements, ",")
	log.SetPrefix("INFO: ")
	log.Println("Processing frags complete", whatHash, data)
	return data, whatHash, nil
}

// sha1Hash - sha1 hashing func
func sha1Hash(data string) string {
	hash := sha1.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

// md5Hash - md5 hashing func
func md5Hash(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

// hashJson - hash data to sh1 or md5
func hashJson(data map[string]interface{}, whatHash []string, hashType *string) (string, error) {
	for k, v := range data {
		log.SetPrefix("INFO: ")
		strV, ok := v.(string)

		for _, hashKey := range whatHash {
			if ok && k == hashKey {
				if *hashType == "sha1" {
					data[k] = sha1Hash(strV)
					log.Printf("Hashed sha1: %v", data[k])
				} else if *hashType == "md5" {
					data[k] = md5Hash(strV)
					log.Printf("Hashed md5: %v", data[k])
				} else {
					return "", fmt.Errorf("unsupported hash type: %v", *hashType)
				}
				break
			}
		}
		if nestedMap, ok := v.(map[string]interface{}); ok {
			_, err := hashJson(nestedMap, whatHash, hashType)
			if err != nil {
				return "", fmt.Errorf("run recursion failed: %v", err)
			}
		}
	}
	resultJson, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON marshalling failed: %v", err)
	}
	return string(resultJson), nil
}
