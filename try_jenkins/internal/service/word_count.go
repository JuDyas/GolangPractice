package service

import (
	"regexp"
	"strings"
)

type WordCountService struct{}

func NewWordCountService() *WordCountService {
	return &WordCountService{}
}

func (svc *WordCountService) CountWords(text string) int {
	text = strings.ToLower(strings.ReplaceAll(text, "\n", " "))

	re := regexp.MustCompile(`[a-zA-Z0-9]*[a-zA-Z]+[a-zA-Z0-9]*`)
	words := re.FindAllString(text, -1)

	return len(words)
}
