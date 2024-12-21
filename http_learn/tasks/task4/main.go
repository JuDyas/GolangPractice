package main

import (
	"GolangPractice/tasks/task7/hasher"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
)

type HashSha1 struct{}

// Hash HashSha1 - method of hashing
func (h *HashSha1) Hash(data string) (string, error) {
	hash := sha1.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func main() {
	http.HandleFunc("/", serveHtml)
	http.HandleFunc("/hash", hashHandler)
	fmt.Println("http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// serveHtml - Serve html for input json data and choice what elements need to hash
func serveHtml(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

// hashHandler - handle input data, hashing and return results
func hashHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed (405)", http.StatusMethodNotAllowed)
		return
	}

	jsonData := r.FormValue("jsonData")
	hashElements := r.FormValue("hashElements")
	data, fields, err := hasher.ProcessFlags(&jsonData, &hashElements)
	if err != nil {
		http.Error(w, "processing JSON error (400)", http.StatusBadRequest)
		return
	}

	hashedResult, err := hasher.HashJson(data, fields, &HashSha1{})
	if err != nil {
		http.Error(w, "hashing data error (500): ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, hashedResult)
}
