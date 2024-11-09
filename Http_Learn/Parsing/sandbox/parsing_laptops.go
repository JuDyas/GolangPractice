package main

import (
	"context"
	"fmt"
	"log"

	"github.com/chromedp/chromedp"
)

type Product struct {
	Name string
	Spec string
}

func main() {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var (
		err      error
		url      = "https://uastore.com.ua/catalog/noutbuki/sort-price/page-all"
		products []Product
	)

	err = chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`.product_preview__name_link`, chromedp.ByQueryAll),
		chromedp.Evaluate(`Array.from(document.querySelectorAll('.product_preview')).map(prod => {
			let name = prod.querySelector('.product_preview__name_link')?.innerText.trim() || '';
			let spec = prod.querySelector('.product_preview__annotation p')?.innerText.trim() || '';
			return { name, spec };
		})`, &products),
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, product := range products {
		fmt.Printf("Name: %s\nSpec: %s\n", product.Name, product.Spec)
		fmt.Println("-----------------------------------------")
	}
}
