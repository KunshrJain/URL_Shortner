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
	urlStore          = make(map[int]string) // ID -> Original URL
)

type myURL struct {
	ID           int       `json:"id"`
	OriginalUrl  string    `json:"original_url"`
	ShortUrl     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

// Create a short URL object and store it in memory
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

// Handler to generate a short URL (expects ?url=...)
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

// Handler to redirect from short URL to original
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL path
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

	// Redirect karne ke liye
	http.Redirect(w, r, original, http.StatusFound)
}

func main() {
	//bhaii yeh sab kuch aur sunega
	http.HandleFunc("/create", createHandler) // /create ko monitor karega
	http.HandleFunc("/", redirectHandler)     // / ko monitor karega

	fmt.Println("Starting URL Shortener...")
	fmt.Println("Go server running at http://localhost:3000")

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
