package config

import (
	"encoding/json"
	"os"

	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/parseserver/models"
)

func LoadHeaders(filePath string) ([]models.UserAgent, error) {
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
