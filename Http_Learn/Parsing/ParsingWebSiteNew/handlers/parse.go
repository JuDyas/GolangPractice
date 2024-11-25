package handlers

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"

	"github.com/go-redis/redis/v8"
)

type Product struct {
	Name  string
	Specs string
}

func ParseHtml(rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sortOpt := r.URL.Query().Get("sort")
		if sortOpt == "" {
			http.Error(w, "sort parameter is required", http.StatusBadRequest)
			return
		}

		if !validSortOpt(sortOpt) {
			http.Error(w, "sort parameter is invalid", http.StatusBadRequest)
			return
		}

		url := "https://uastore.com.ua/catalog/noutbuki/" + sortOpt + "/page-all"
		err, products := goParse(url, sortOpt)
		if err != nil {
			http.Error(w, "parsing error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		for i, product := range products {
			fmt.Fprintf(w, "Товар %d:\n", i+1)
			fmt.Fprintf(w, "  Название: %s\n", product.Name)
			fmt.Fprintf(w, "  Характеристики: %s\n\n", product.Specs)
		}

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

func goParse(url, sortOpt string) (error, []Product) {
	flag.Parse()

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("http get error: %w", err), nil
	}

	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("html parse error: %w", err), nil
	}

	products := extractProduct(doc)
	return nil, products
}

func extractProduct(n *html.Node) []Product {
	var (
		products  []Product
		parseNode func(*html.Node)
	)
	parseNode = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, attr := range n.Attr {
				if attr.Key == "class" && strings.Contains(attr.Val, "product_preview") {
					name := extractText(n, "product_preview__name_link")
					specs := extractText(n, "product_preview__annotation")
					if name != "" && specs != "" {
						products = append(products, Product{Name: name, Specs: specs})
					}
				}
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			parseNode(child)
		}
	}
	parseNode(n)
	return products
}

func extractText(n *html.Node, className string) string {
	var (
		res        string
		searchNode func(*html.Node)
	)

	searchNode = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, attr := range n.Attr {
				if attr.Key == "class" && strings.Contains(attr.Val, className) {
					res = getTextFromNode(n)
					return
				}
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			searchNode(child)
		}
	}
	searchNode(n)
	return res
}

func getTextFromNode(n *html.Node) string {
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data)
	}
	var result strings.Builder
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		result.WriteString(getTextFromNode(child))
	}
	return result.String()
}
