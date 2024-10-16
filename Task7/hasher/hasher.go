package hasher

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Hasher interface {
	Hash(string) (string, error)
}

// ProcessFlags - processing input data to map and slice
func ProcessFlags(jsonData, hashElements *string) (map[string]interface{}, []string, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(*jsonData), &data)
	if err != nil {
		return nil, nil, err
	}

	whatHash := strings.Split(*hashElements, ",")
	return data, whatHash, nil
}

// HashJson - hash data to any hash
func HashJson(data map[string]interface{}, whatHash []string, hasher Hasher) (string, error) {
	var err error
	for k, v := range data {
		strV, ok := v.(string)
		for _, hashKey := range whatHash {
			if ok && k == hashKey {
				data[k], err = hasher.Hash(strV)
				if err != nil {
					return "", fmt.Errorf("hasher: %w", err)
				}
			}
		}

		if nestedMap, ok := v.(map[string]interface{}); ok {
			_, err := HashJson(nestedMap, whatHash, hasher)
			if err != nil {
				return "", fmt.Errorf("run recursion: %v", err)
			}
		}
	}

	resultJson, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshalling: %v", err)
	}

	return string(resultJson), nil
}
