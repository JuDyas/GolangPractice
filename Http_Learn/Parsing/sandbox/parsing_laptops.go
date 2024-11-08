package main

import (
	"context"
	"fmt"
	"log"

	"github.com/chromedp/chromedp"
)

func main() {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var (
		err       error
		url       = "https://uastore.com.ua/catalog/noutbuki"
		prodName  []string
		prodSpecs []string
	)
	err = chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`.product_preview__name_link`, chromedp.ByQueryAll),
		chromedp.Evaluate(`Array.from(document.querySelectorAll('.product_preview__name_link')).map(el => el.innerText.trim())`, &prodName),
		chromedp.Evaluate(`Array.from(document.querySelectorAll('.product_preview__annotation p')).map(el => el.innerText.trim())`, &prodSpecs),
	)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(prodName); i++ {
		fmt.Printf("Name: %s\n ", prodName[i])
		if i < len(prodSpecs) {
			fmt.Printf("Spec: %s\n", prodSpecs[i])
		}
		fmt.Println("-----------------------------------------")
	}
}
