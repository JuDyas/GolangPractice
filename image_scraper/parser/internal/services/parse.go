package services

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/parser/config"
	"github.com/gocolly/colly"
)

const (
	linksSelector = "a[href]"
	imgSelector   = "img[src]"
)

type Parser interface {
	Start(url string)
	Stop()
	Pause()
	Continue()
}

type WebParser struct {
	url       string
	paused    bool
	mutex     sync.Mutex
	stopChan  chan struct{}
	pauseChan chan struct{}
	redis     *redis.Client
	domain    string
}

func NewWebParser(redisAddr string) *WebParser {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	return &WebParser{
		stopChan:  make(chan struct{}),
		pauseChan: make(chan struct{}),
		redis:     rdb,
	}
}

func (wp *WebParser) Start(url string) {
	wp.url = url
	wp.paused = false
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

func (wp *WebParser) parseProcess(url string) {
	ctx := context.Background()
	ttl, err := wp.redis.TTL(ctx, url).Result()
	if err != nil {
		log.Println("check TTL:", err)
		return
	}

	if ttl > 0 {
		log.Println("ulr exist, still alive for:", ttl)
		return
	}

	err = wp.redis.Set(ctx, url, "visited", 2*time.Minute).Err()
	if err != nil {
		log.Println("set: ", err)
		return
	}

	err = wp.parsePage(
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
		})

	if err != nil {
		fmt.Println(err)
	}
}

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

func (wp *WebParser) Stop() {

}

func (wp *WebParser) Pause() {

}

func (wp *WebParser) Continue() {

}
