package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	pageUrl     = "https://uastore.com.ua"
	catalogPath = "/catalog/noutbuki/sort-"
)

type Specs struct {
	DisplaySize       string
	DisplayResolution string
	Cpu               string
	Ram               string
	HardDrives        string
	Gpu               string
}

type Product struct {
	Name  string
	Specs Specs //details/info/data
}

// add Context

func ParseHtml(productChannel chan []Product) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sortOpt := r.URL.Query().Get("sort")
		if sortOpt == "" {
			http.Error(w, "Sort parameter is required", http.StatusBadRequest)
			return
		}

		if !validSortOpt(sortOpt) {
			http.Error(w, "Sort parameter is invalid", http.StatusBadRequest)
			return
		}

		limit := r.URL.Query().Get("limit")
		if limit == "" {
			http.Error(w, "limit parameter is required (only numbers)", http.StatusBadRequest)
			return
		}

		offset := r.URL.Query().Get("offset")
		if offset == "" {
			http.Error(w, "skip parameter is required (only numbers)", http.StatusBadRequest)
			return
		}

		url := pageUrl + catalogPath + sortOpt + "/page-all"
		err, products := goParse(url, offset, limit)
		if err != nil {
			http.Error(w, "parsing error", http.StatusInternalServerError)
			return
		}
		productChannel <- products
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, "Parse complete")
	}
}

func validSortOpt(sortOpt string) bool {
	var validSortOpts = []string{
		"position", "price", "price_desc",
		"name", "name_desc", "rating",
		"rating_desc",
	}

	for _, validSortOpt := range validSortOpts {
		if sortOpt == validSortOpt {
			return true
		}
	}

	return false
}

func goParse(url, offsetStr, limitStr string) (error, []Product) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("http get error: %d", err)
		return err, nil
	}

	defer resp.Body.Close()

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		log.Printf("offset parameter is invalid")
		return err, nil
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Printf("limit parameter is invalid")
		return err, nil
	}

	products, err := extractProduct(resp.Body, offset, limit)
	if err != nil {
		log.Println(err)
		return err, nil
	}

	return nil, products
}

func extractProduct(r io.Reader, offset, limit int) ([]Product, error) {
	var (
		products []Product
		count    int
	)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("goquery document error: %s", err)
	}

	doc.Find("a.product_preview__name_link").Each(func(i int, s *goquery.Selection) {
		if offset > 0 && count < offset {
			count++
			return
		}

		if limit > 0 && len(products) >= limit {
			return
		}

		href, _ := s.Attr("href")
		name := strings.TrimSpace(s.Contents().Not(".product_preview__sku").Text())
		spec := processProduct(pageUrl + href)
		products = append(products, Product{
			Name:  name,
			Specs: spec,
		})
	})
	return products, nil
}

func processProduct(url string) Specs {
	var (
		spec = []string{
			"Диагональ:", "Разрешение:", "Видеокарта:",
			"Процессор:", "Объем оперативной памяти:", "Объем накопителя:",
		}
	)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("http get (product) error: %d", err)
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("goquery document (product) error: %s", err)
	}

	prod := Specs{}
	doc.Find("ul.d-sm-flex.flex-sm-wrap.features.mobile_tab__content").Each(func(i int, s *goquery.Selection) {
		s.Find("li").Each(func(i int, s *goquery.Selection) {
			name := strings.TrimSpace(s.Find(".features__name").Text())
			val := strings.TrimSpace(s.Find(".features__value").Text())
			for _, s := range spec {
				if strings.Contains(name, s) {
					switch s {
					case "Диагональ:":
						prod.DisplaySize = val
					case "Разрешение:":
						prod.DisplayResolution = val
					case "Видеокарта:":
						prod.Gpu = val
					case "Процессор:":
						prod.Cpu = val
					case "Объем оперативной памяти:":
						prod.Ram = val
					case "Объем накопителя:":
						prod.HardDrives = val
					}
				}
			}
		})
	})
	return prod
}
