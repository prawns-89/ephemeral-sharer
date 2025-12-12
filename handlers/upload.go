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

	// --- START MAGIC BYTES CHECK ---
	
	// A. Read the first 512 bytes
	// Go's standard library only needs the first 512 bytes to guess the type.
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		http.Error(w, "Server error reading file", http.StatusInternalServerError)
		return
	}

	// B. Detect the type
	fileType := http.DetectContentType(buff)
	// Example result: "image/png" or "text/plain; charset=utf-8"

	// C. Whitelist allowed types
	if fileType != "image/jpeg" && fileType != "image/png" && fileType != "application/pdf" {
		http.Error(w, "Invalid file type. Only JPEG, PNG, and PDF allowed.", http.StatusBadRequest)
		return
	}

	// D. CRITICAL STEP: Rewind the file!
	// We just read 512 bytes. If we save the file now, the saved file will be missing
	// the beginning (it will be corrupt). We must "seek" back to byte 0.
	file.Seek(0, 0)

	// --- END MAGIC BYTES CHECK ---

	// 2. Logic (Generate Name)
	ext := filepath.Ext(header.Filename)
	newFilename := uuid.New().String() + ext

	// 3. Save to Storage
	if err := h.Store.Save(newFilename, file); err != nil {
		http.Error(w, "Storage failed", http.StatusInternalServerError)
		return
	}

	// 4. Response (Fixed with "success: true")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"link":    "http://localhost:8080/files/" + newFilename,
		"type":    fileType,
	})
}