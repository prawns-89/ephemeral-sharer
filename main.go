package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prawns-89/ephemeral-sharer/handlers"
	"github.com/prawns-89/ephemeral-sharer/storage"
)

func main() {
	// 1. Initialize components
	store := storage.New("./uploads")
	uploadHandler := &handlers.UploadHandler{Store: store}

	// 2. Start Background Worker (The Cleaner)
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			fmt.Println("🧹 Pruning old files...")
			store.Prune(5 * time.Minute)
		}
	}()

	// 3. Routes
	http.Handle("/", http.FileServer(http.Dir("./static"))) // Serves index.html directly
	http.HandleFunc("/upload", uploadHandler.HandleHTTP)
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("./uploads"))))

	// 4. Start
	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}