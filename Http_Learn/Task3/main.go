package main

import (
	"fmt"
	"log"
	"net/http"
)

//сервер, который принимает и выводит заголовки запроса.

func main() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	usAg := r.Header.Get("User-Agent")
	content := r.Header.Get("Content-Type")
	lang := r.Header.Get("Accept-Language")

	fmt.Fprintf(w, "User-Agent: %s", usAg)
	fmt.Fprintf(w, "Content-Type: %s", content)
	fmt.Fprintf(w, "Accept-Language: %s", lang)
}
