package config

import (
	"encoding/json"
	"fmt"
	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/parser/internal/models"
	"os"
	"path/filepath"
)

func LoadHeaders() ([]models.UserAgent, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config"
		fmt.Println("CONFIG_PATH environment variable not set, use:", configPath)
	}

	filePath := filepath.Join(configPath, "headers.json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var userAgents []models.UserAgent
	err = json.Unmarshal(data, &userAgents)
	if err != nil {
		return nil, err
	}

	return userAgents, nil
}
