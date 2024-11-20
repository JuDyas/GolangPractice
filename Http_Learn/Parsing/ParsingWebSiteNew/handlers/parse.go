package handlers

import (
	"GolangPractice/Http_Learn/Parsing/ParsingWebSiteNew/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/chromedp/chromedp"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Specs string `json:"specs"`
}

func ParseHtml(rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		sortOptions := []string{
			"sort-position", "sort-price", "sort-price_desc",
			"sort-name", "sort-name_desc", "sort-rating",
			"sort-rating_desc",
		}

		for _, sortParam := range sortOptions {
			url := "https://uastore.com.ua/catalog/noutbuki/" + sortParam + "/page-all"
			prod, err := goParse(ctx, url)
			if err != nil {
				log.Println(err)
			}

			err = saveToDB(prod, rdb, sortParam, utils.Key)
			if err != nil {
				log.Println(err)
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}

func goParse(ctx context.Context, url string) ([]Product, error) {
	var products []Product
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`.product_preview__name_link`, chromedp.ByQueryAll), // Ждем видимости элементов
		chromedp.Evaluate(`Array.from(document.querySelectorAll('.product_preview')).map(prod => {
					let name = prod.querySelector('.product_preview__name_link')?.innerText.trim() || '';
					let specs = prod.querySelector('.product_preview__annotation p')?.innerText.trim() || '';
					return {name, specs};
				})`, &products),
	)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func saveToDB(prod []Product, rdb *redis.Client, sp, key string) error {
	for i, product := range prod {
		prod[i].ID = uuid.New().String()
		prod[i].Specs = utils.EncodeData(key, product.Specs)
	}

	keyDb := "products:" + sp
	jsonData, err := json.Marshal(prod)
	if err != nil {
		return fmt.Errorf("error converting data to json %s: %v", sp, err)
	}

	err = rdb.Set(context.Background(), keyDb, jsonData, 0).Err()
	if err != nil {
		return fmt.Errorf("error saving data to redis %s: %v", sp, err)
	}
	return nil
}
