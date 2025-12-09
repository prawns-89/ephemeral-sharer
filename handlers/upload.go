package handlers

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/prawns-89/ephemeral-sharer/storage"
)

// UploadHandler needs access to the storage system
type UploadHandler struct {
	Store *storage.LocalStore
}

func (h *UploadHandler) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. Validation
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.ParseMultipartForm(10 << 20) // 10 MB

	file, header, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 2. Logic
	ext := filepath.Ext(header.Filename)
	newFilename := uuid.New().String() + ext

	// 3. Call Storage Layer
	if err := h.Store.Save(newFilename, file); err != nil {
		http.Error(w, "Storage failed", http.StatusInternalServerError)
		return
	}

	// 4. Response
	json.NewEncoder(w).Encode(map[string]string{
		"link": "http://localhost:8080/files/" + newFilename,
	})
}