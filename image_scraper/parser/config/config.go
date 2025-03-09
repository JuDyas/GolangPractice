package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/JuDyas/GolangPractice/image_scraper/parser/internal/models"
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
