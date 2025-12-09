package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// 1. The Upload Form
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `
		<html>
		<body>
			<h2>Upload a File</h2>
			<form action="/upload" method="post" enctype="multipart/form-data">
				<input type="file" name="myFile" />
				<input type="submit" value="Upload" />
			</form>
		</body>
		</html>`
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, html)
	})

	// 2. The Upload Handler
	http.HandleFunc("/upload", uploadHandler)

	// 3. The File Server (The part that was likely broken)
	// We point this directly to the "uploads" folder on your disk
	fs := http.FileServer(http.Dir("./uploads"))
	
	// We strip "/files/" from the URL so the server looks for just "image.jpg"
	http.Handle("/files/", http.StripPrefix("/files/", fs))

	fmt.Println("Server starting on http://localhost:8080 ...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File upload started...") // Debug print

	// Limit upload to 10MB
	r.ParseMultipartForm(10 << 20) 

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error retrieving file:", err)
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Ensure the uploads directory exists
	os.MkdirAll("./uploads", os.ModePerm)

	// Create the file
	dstPath := filepath.Join("./uploads", handler.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		http.Error(w, "Error creating file on server", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the bytes
	if _, err := io.Copy(dst, file); err != nil {
		fmt.Println("Error saving file:", err)
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Uploaded successfully: %s\n", handler.Filename)

	// Send the link back to the user
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h2>Upload Successful!</h2>")
	fmt.Fprintf(w, `<p>Saved to: %s</p>`, dstPath)
	fmt.Fprintf(w, `<a href="/files/%s">Click here to download %s</a>`, handler.Filename, handler.Filename)
}