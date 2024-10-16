package main

import (
	"fmt"
	"net/http"
)

// сервер, который будет выводить HTTP-метод, путь и заголовки запроса.

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	path := r.URL.Path
	userAgent := r.Header.Get("User-Agent")
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "User Agent: %s\n", userAgent)
	fmt.Fprintf(w, "Path: %s\n", path)
	fmt.Fprintf(w, "Method: %s\n", method)
}
