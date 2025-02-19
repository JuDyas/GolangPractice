package services

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"sync"
	"time"

	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/parser/internal/repositories"

	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/parser/internal/models"

	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/parser/config"
	"github.com/gocolly/colly"
)

const (
	linksSelector = "a[href]"
	imgSelector   = "img[src]"
)

type Parser interface {
	Start(url string)
}

type WebParser struct {
	url  string
	repo repositories.LinksRepository
	//TODO: Сделать проверку на паузу, чтобы нельзя было запустить парсинг прежде, чем остановится прошлый.
	//paused      bool
	mutex        sync.Mutex
	controlChan  <-chan models.CommandType
	imageUrlChan chan<- string
	domain       string
}

func NewWebParser(repo repositories.LinksRepository, controlChan <-chan models.CommandType, imageUrlChan chan<- string) *WebParser {
	return &WebParser{
		controlChan:  controlChan,
		imageUrlChan: imageUrlChan,
		repo:         repo,
	}
}

// Start - start parser process
func (wp *WebParser) Start(url string) {
	wp.url = url
	wp.mutex.Lock()
	defer wp.mutex.Unlock()

	log.Println("parser start")
	wp.parse(url)
}

func (wp *WebParser) parse(webURL string) {
	domain, err := url.Parse(webURL)
	if err != nil {
		log.Fatal("Ошибка парсинга начального URL:", err)
	}
	wp.domain = domain.Host
	wp.parseProcess(webURL)
}

// parseProcess - check exist link in redis and start parsing page with parameters
func (wp *WebParser) parseProcess(url string) {
	ctx := context.Background()
	exist := wp.repo.CheckExists(ctx, url)
	if exist {
		return
	} else {
		wp.repo.SaveLink(ctx, url)
	}

	select {
	case cmd := <-wp.controlChan:
		if cmd == models.CmdPause {
			fmt.Println(models.CmdPause)
			for {
				if cmd := <-wp.controlChan; cmd == models.CmdResume {
					fmt.Println(models.CmdResume)
					break
				} else if cmd == models.CmdStop {
					fmt.Println(models.CmdStop)
					return
				}
			}
		}
	default:
		err := wp.parsePage(
			url,
			linksSelector,
			imgSelector,
			func(e *colly.HTMLElement) {
				link := e.Request.AbsoluteURL(e.Attr("href"))
				fmt.Println("link found: ", link)
				wp.parseProcess(link)
			},
			func(e *colly.HTMLElement) {
				imgSrc := e.Request.AbsoluteURL(e.Attr("src"))
				fmt.Println("imgSrc found: ", imgSrc)
				exist = wp.repo.CheckExists(ctx, imgSrc)
				if exist {
					fmt.Println("img already exists")
				} else {
					wp.repo.SaveLink(ctx, imgSrc)
					wp.imageUrlChan <- imgSrc
				}
			})

		if err != nil {
			fmt.Println(err)
		}
	}
}

// parsePage - Parse page with selectors and process functions
func (wp *WebParser) parsePage(url string, selectorLinks string, selectorImg string, processFuncLinks, processFuncImg func(e *colly.HTMLElement)) error {
	c := colly.NewCollector(colly.AllowedDomains(wp.domain))
	rand.Seed(time.Now().UnixNano())
	err := c.Limit(&colly.LimitRule{
		DomainGlob:  wp.domain,
		Parallelism: 1,
		RandomDelay: 3 * time.Second,
		Delay:       1 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("set limit: %v", err)
	}

	headers, err := config.LoadHeaders()
	if err != nil {
		return fmt.Errorf("load user agents: %v", err)
	}

	c.OnRequest(func(r *colly.Request) {
		header := headers[rand.Intn(len(headers))]
		fmt.Println(header)
		r.Headers.Set("User-Agent", header.UserAgent)
		r.Headers.Set("Accept", header.Accept)
		r.Headers.Set("Accept-Encoding", header.AcceptEncoding)
		r.Headers.Set("Accept-Language", header.AcceptLanguage)
		r.Headers.Set("Referer", header.Referer)
	})

	c.OnHTML(selectorImg, processFuncImg)
	c.OnHTML(selectorLinks, processFuncLinks)

	err = c.Visit(url)
	if err != nil {
		return fmt.Errorf("visit: %v", err)
	}

	c.Wait()
	return nil
}
