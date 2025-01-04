package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/JuDyas/GolangPractice/pastebin_new/models"
)

func LoadConfig() (*models.Config, error) {
	var config models.Config
	file, err := os.Open("./config/config.json")
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}

	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("error decoding file: %v", err)
	}

	return &config, nil
}
