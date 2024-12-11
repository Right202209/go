package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse uploaded file
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to upload file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save file to the server
	dst, err := os.Create(fmt.Sprintf("./uploads/%s", header.Filename))
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("File uploaded successfully"))
}
