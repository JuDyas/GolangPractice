package main

import (
	"GolangPractice/Http_Learn/Parsing/ParsingWebSiteNew/pkg/vyzhenercipher"
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

var (
	rdb        *redis.Client
	chipherKey *string
	ctx        = context.Background()
)

type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Specs string `json:"specs"`
}

func main() {
	chipherKey = flag.String("key", "BanaNa", "Key for vygener cipher")
	flag.Parse()

	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)

	http.HandleFunc("/parse", parseHtml)
	http.HandleFunc("/getdata", getDataFromRedis)
	http.HandleFunc("/product/", getProductDetails)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Ошибка при запуске сервера: ", err)
	}
}

// parseHtml - parse data from ulr
func parseHtml(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	sortOptions := []string{
		"sort-position", "sort-price", "sort-price_desc", "sort-name", "sort-name_desc", "sort-rating", "sort-rating_desc",
	}

	for _, sortParam := range sortOptions {
		url := "https://uastore.com.ua/catalog/noutbuki/" + sortParam + "/page-all"
		var products []Product

		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.WaitVisible(`.product_preview__name_link`, chromedp.ByQueryAll),
			chromedp.Evaluate(`Array.from(document.querySelectorAll('.product_preview')).map(prod => {
                let name = prod.querySelector('.product_preview__name_link')?.innerText.trim() || '';
                let specs = prod.querySelector('.product_preview__annotation p')?.innerText.trim() || '';
                return {name, specs};
            })`, &products),
		)
		if err != nil {
			log.Printf("error parsng with sort: %s: %v", sortParam, err)
			continue
		}

		for i, product := range products {
			products[i].ID = uuid.New().String()
			products[i].Specs = encodeData(*chipherKey, product.Specs)
		}

		key := "products:" + sortParam
		jsonData, err := json.Marshal(products)
		if err != nil {
			log.Printf("error convert data to json %s: %v", sortParam, err)
			continue
		}
		err = rdb.Set(ctx, key, jsonData, 0).Err()
		if err != nil {
			log.Printf("error save data to redis %s: %v", sortParam, err)
		}
	}
}

// getDataFromRedis - get data from redis
func getDataFromRedis(w http.ResponseWriter, r *http.Request) {
	sortParam := r.URL.Query().Get("sort")
	if sortParam == "" {
		sortParam = "sort-position"
	}

	key := "products:" + sortParam
	jsonData, err := rdb.Get(ctx, key).Result()
	if err != nil {
		log.Printf("error get data from redis %s: %v", key, err)
		return
	}

	var products []Product
	err = json.Unmarshal([]byte(jsonData), &products)
	if err != nil {
		log.Printf("error for convert data %s: %v", key, err)
		return
	}

	for i, product := range products {
		products[i].Specs = decodeData(*chipherKey, product.Specs)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// getProductDetails - get data for current product
func getProductDetails(w http.ResponseWriter, r *http.Request) {
	productID := strings.TrimPrefix(r.URL.Path, "/product/")
	if productID == "" {
		return
	}
	sortOptions := []string{
		"sort-position", "sort-price", "sort-price_desc", "sort-name", "sort-name_desc", "sort-rating", "sort-rating_desc",
	}
	var product Product
	ok := false

	for _, sortParam := range sortOptions {
		key := "products:" + sortParam
		jsonData, err := rdb.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var products []Product
		err = json.Unmarshal([]byte(jsonData), &products)
		if err != nil {
			continue
		}

		for _, p := range products {
			if p.ID == productID {
				p.Specs = decodeData(*chipherKey, p.Specs)
				product = p
				ok = true
				break
			}
		}
		if ok {
			break
		}
	}

	if !ok {
		http.Error(w, "product not exist ", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(product)
	if err != nil {
		return
	}
}

// encodeData - encode data
func encodeData(key, val string) string {
	chipher := &vyzhenercipher.VizhenerCipher{
		Key:  key,
		Text: val,
	}

	chipher.Encrypt()
	return chipher.ChangedText
}

// decodeData - decode data
func decodeData(key, val string) string {
	chipher := &vyzhenercipher.VizhenerCipher{
		Key:  key,
		Text: val,
	}

	chipher.Decrypt()
	return chipher.ChangedText
}
