package services

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/parseserver/config"

	"github.com/gocolly/colly"
)

type Parser interface {
	Parse(url string, output chan<- string) error
}

type CollyParser struct {
}

const (
	//Moyo classes
	classCategory = "div.catalog-submenu-item:not(.child)"
	classProduct  = "a.gtm-link-product"
	classImages   = "div.product_image_item-wrap img"

	//Comfy classes
	//classCategory = "a.menu-desktop-content__top-link"
	//classProduct  = "header.prdl-item__info-head a.prdl-item__name"
	//classImages   = "ul.bc-photos__gallery img"
)

// Random delay (sleep)
func delay() {
	sleepTime := time.Duration(rand.Intn(2)+1) * time.Second
	fmt.Printf("Задержка %v...\n", sleepTime)
	time.Sleep(sleepTime)
}

// Universal parse func with random headers and limits
func parsePage(url string, selector string, processFunc func(e *colly.HTMLElement)) error {
	c := colly.NewCollector()
	rand.Seed(time.Now().UnixNano())
	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
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

	c.OnHTML(selector, processFunc)
	err = c.Visit(url)
	if err != nil {
		return fmt.Errorf("visit: %v", err)
	}

	c.Wait()
	return nil
}

// Parse product links from site
func parseProducts(categoryL []string, output chan<- string) ([]string, error) {
	var productLinks []string
	for _, category := range categoryL {
		delay()
		err := parsePage(category, classProduct, func(e *colly.HTMLElement) {
			link := e.Request.AbsoluteURL(e.Attr("href"))
			fmt.Println("Product found:", link)
			productLinks = append(productLinks, link)
			err := parseImages(link, output)
			if err != nil {
				fmt.Println("images parse:", err)
			}
		})

		if err != nil {
			return nil, fmt.Errorf("parse products: %v", err)
		}
	}

	return productLinks, nil
}

// Parse images from product pages
func parseImages(productL string, output chan<- string) error {
	for _, product := range productL
	fmt.Println("Visiting:", productL)
	delay()
	err := parsePage(productL, classImages, func(e *colly.HTMLElement) {
		imgSrc := e.Request.AbsoluteURL(e.Attr("src"))
		fmt.Println("Image found:", imgSrc)
		output <- imgSrc
	})

	if err != nil {
		return fmt.Errorf("parse images: %v", err)
	}
	return nil
	}
}

// Parse - Start parsing links and images
func (cp CollyParser) Parse(url string, output chan<- string) error {
	defer close(output)
	var categoryLinks []string
	err := parsePage(url, classCategory, func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("data-href"))
		fmt.Println("found category:", link)
		categoryLinks = append(categoryLinks, link)
	})

	if err != nil {
		fmt.Println(err)
		return err
	}

	productLinks, err := parseProducts(categoryLinks, output)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_ = productLinks

	return nil
}

type ParseService struct {
	Parser Parser
}

func (p ParseService) ParseImages(url string, output chan<- string) error {
	return p.Parser.Parse(url, output)
}
