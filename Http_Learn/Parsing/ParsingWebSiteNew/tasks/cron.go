package tasks

import (
	"GolangPractice/Http_Learn/Parsing/ParsingWebSiteNew/handlers"
	"log"

	"github.com/robfig/cron/v3"
)

func InitCron() {
	c := cron.New()
	_, err := c.AddFunc("*/5 * * * * *", func() {
		handlers.ParseHtml(nil)
	})
	if err != nil {
		log.Printf("cron addFunc err:%v", err)
	}

}
