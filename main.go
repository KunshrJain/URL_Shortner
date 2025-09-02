package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	globalCounter int = 0
	urlStore          = make(map[int]string) 
)

type myURL struct {
	ID           int       `json:"id"`
	OriginalUrl  string    `json:"original_url"`
	ShortUrl     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

func createShortURL(originalURL string) myURL {
	globalCounter++
	id := globalCounter
	urlStore[id] = originalURL

	return myURL{
		ID:           id,
		OriginalUrl:  originalURL,
		ShortUrl:     fmt.Sprintf("http://localhost:3000/%d", id),
		CreationDate: time.Now(),
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	original := r.URL.Query().Get("url")
	if original == "" {
		http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
		return
	}

	short := createShortURL(original)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(short)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	original, exists := urlStore[id]
	if !exists {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, original, http.StatusFound)
}

func main() {
	http.HandleFunc("/create", createHandler) 
	http.HandleFunc("/", redirectHandler)    

	fmt.Println("Starting URL Shortener...")
	fmt.Println("Go server running at http://localhost:3000")

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

