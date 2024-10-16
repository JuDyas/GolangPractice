package main

import (
	"fmt"
	"log"
	"net/http"
)

//сервер, который принимает данные формы через POST и возвращает их обратно в ответе.
// подключение html для ввода данных в форму

func main() {
	http.HandleFunc("/", htmlServe)
	http.HandleFunc("/submit", handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func htmlServe(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed (405)", http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Parameter name is empty (400)", http.StatusBadRequest)
	}

	fmt.Fprintf(w, "Name: %s", name)
}
