package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-redis/redis/v8"
)

type Product struct {
	Name  string
	Specs string
}

func ParseHtml(rdb *redis.Client, productChannel chan []Product) http.HandlerFunc {
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
		}

		skip := r.URL.Query().Get("skip")
		if skip == "" {
			http.Error(w, "skip parameter is required (only numbers)", http.StatusBadRequest)
		}

		url := "https://uastore.com.ua/catalog/noutbuki/" + sortOpt + "/page-all"
		err, products := goParse(url, skip, limit)
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
		"sort-position", "sort-price", "sort-price_desc",
		"sort-name", "sort-name_desc", "sort-rating",
		"sort-rating_desc",
	}
	for _, validSortOpt := range validSortOpts {
		if sortOpt == validSortOpt {
			return true
		}
	}
	return false
}

func goParse(url, skipStr, limitStr string) (error, []Product) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("http get error: %d", err)
		return err, nil
	}
	defer resp.Body.Close()

	skip, err := strconv.Atoi(skipStr)
	if err != nil {
		log.Printf("skip parameter is invalid")
		return err, nil
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Printf("limit parameter is invalid")
		return err, nil
	}

	products, err := extractProduct(resp.Body, skip, limit)
	if err != nil {
		log.Println(err)
		return err, nil
	}

	return nil, products
}

func extractProduct(body io.ReadCloser, skip, limit int) ([]Product, error) {
	var (
		products []Product
		count    int
	)
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, fmt.Errorf("goquery.NewDocumentFromReader: %w", err)
	}

	doc.Find(".fn_transfer.clearfix").Each(func(i int, s *goquery.Selection) {
		if skip > 0 && count < skip {
			count++
			return
		}
		if limit > 0 && len(products) >= limit {
			return
		}
		name := strings.TrimSpace(s.Find(".product_preview__name_link").Contents().Not(".product_preview__sku").Text())
		specs := strings.TrimSpace(s.Find(".product_preview__annotation p").Text())
		products = append(products, Product{
			Name:  name,
			Specs: specs,
		})
	})
	return products, nil

}
