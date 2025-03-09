package services

import (
	"context"
	"log"
	"math/rand"
	"net/url"
	"sync"

	"github.com/JuDyas/GolangPractice/image_scraper/parser/config"
	"github.com/JuDyas/GolangPractice/image_scraper/parser/internal/models"
	"github.com/JuDyas/GolangPractice/image_scraper/parser/internal/repositories"

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
	url          string
	repo         repositories.LinksRepository
	controlChan  chan models.CommandType
	imageUrlChan chan<- string
	domain       string
	mutex        sync.Mutex
	wg           sync.WaitGroup
	paused       bool
	stopFlag     bool
}

func NewWebParser(repo repositories.LinksRepository, controlChan chan models.CommandType, imageUrlChan chan<- string) *WebParser {
	return &WebParser{
		controlChan:  controlChan,
		imageUrlChan: imageUrlChan,
		repo:         repo,
	}
}

func (wp *WebParser) Start(url string) {
	wp.url = url
	wp.mutex.Lock()
	defer wp.mutex.Unlock()

	log.Println("Парсер запущен")
	wp.processCommands()
	wp.parse(url)
}

func (wp *WebParser) processCommands() {
	go func() {
		for cmd := range wp.controlChan {
			wp.mutex.Lock()
			switch cmd {
			case models.CmdPause:
				wp.paused = true
				log.Println("Парсинг приостановлен")
			case models.CmdResume:
				wp.paused = false
				log.Println("Парсинг продолжен")
			case models.CmdStop:
				wp.stopFlag = true
				log.Println("Парсинг остановлен")
				wp.mutex.Unlock()
				return
			}
			wp.mutex.Unlock()
		}
	}()
}

func (wp *WebParser) parse(webURL string) {
	domain, err := url.Parse(webURL)
	if err != nil {
		log.Fatalf("parse error URL: %v", err)
	}
	wp.domain = domain.Host
	wp.parseProcess(webURL)
}

func (wp *WebParser) parseProcess(webURL string) {
	ctx := context.Background()
	if wp.repo.CheckExists(ctx, webURL) {
		return
	}
	wp.repo.SaveLink(ctx, webURL)

	for {
		wp.mutex.Lock()
		if wp.stopFlag {
			wp.mutex.Unlock()
			return
		}
		if wp.paused {
			wp.mutex.Unlock()
			<-wp.controlChan
			continue
		}
		wp.mutex.Unlock()

		wp.parsePage(webURL, linksSelector, imgSelector)
		return
	}
}

func (wp *WebParser) parsePage(url string, linkSelector, imgSelector string) {
	wp.mutex.Lock()
	if wp.stopFlag {
		wp.mutex.Unlock()
		return
	}
	wp.mutex.Unlock()

	c := colly.NewCollector(
		colly.AllowedDomains(wp.domain),
	)

	c.OnRequest(func(r *colly.Request) {
		headers, err := config.LoadHeaders()
		if err != nil {
			log.Fatalf("loading headers: %v", err)
		}
		header := headers[rand.Intn(len(headers))]
		r.Headers.Set("User-Agent", header.UserAgent)
		r.Headers.Set("Accept", header.Accept)
		r.Headers.Set("Accept-Encoding", header.AcceptEncoding)
		r.Headers.Set("Accept-Language", header.AcceptLanguage)
		r.Headers.Set("Referer", header.Referer)
	})

	c.OnHTML(linkSelector, func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if !wp.isSameDomain(link) {
			return
		}
		log.Printf("found link: %s", link)
		wp.parseProcess(link)
	})

	c.OnHTML(imgSelector, func(e *colly.HTMLElement) {
		imgSrc := e.Request.AbsoluteURL(e.Attr("src"))
		if wp.repo.CheckExists(context.Background(), imgSrc) {
			log.Printf("image not found: %s", imgSrc)
			return
		}
		log.Printf("found image: %s", imgSrc)
		wp.repo.SaveLink(context.Background(), imgSrc)
		wp.imageUrlChan <- imgSrc
	})

	err := c.Visit(url)
	if err != nil {
		log.Printf("wisiting URL: %v", err)
	}
	c.Wait()
}

func (wp *WebParser) isSameDomain(link string) bool {
	parsedURL, err := url.Parse(link)
	if err != nil {
		log.Printf("url parse: %v", err)
		return false
	}
	return parsedURL.Host == wp.domain
}
