package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Server handles HTTP requests for serving encrypted files
type Server struct {
	key   []byte
	files map[string]EmbeddedFile
}

// NewServer creates a new server instance
func NewServer(key []byte) (*Server, error) {
	files, err := GetEmbeddedFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to load embedded files: %w", err)
	}

	return &Server{
		key:   key,
		files: files,
	}, nil
}

// Start starts the HTTP server on the specified port
func (s *Server) Start(port string) error {
	http.HandleFunc("/", s.handleRequest)

	addr := ":" + port
	fmt.Printf("Starting server on http://localhost%s\n", addr)
	fmt.Printf("Available files: %d\n", len(s.files))

	return http.ListenAndServe(addr, nil)
}

// handleRequest routes requests to either list files or serve a specific file
func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")

	// Root path - return list of files
	if path == "" {
		s.handleListFiles(w, r)
		return
	}

	// Serve specific file
	s.handleServeFile(w, r, path)
}

// handleListFiles returns a JSON list of available file names
func (s *Server) handleListFiles(w http.ResponseWriter, r *http.Request) {
	fileNames := make([]string, 0, len(s.files))
	for name := range s.files {
		fileNames = append(fileNames, name)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(fileNames); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// handleServeFile decrypts and serves a specific file
func (s *Server) handleServeFile(w http.ResponseWriter, r *http.Request, fileName string) {
	file, exists := s.files[fileName]
	if !exists {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Decrypt the file
	decrypted, err := Decrypt(file.EncryptedData, s.key)
	if err != nil {
		log.Printf("Error decrypting file %s: %v", fileName, err)
		http.Error(w, "Failed to decrypt file", http.StatusInternalServerError)
		return
	}

	// Set content type and serve the file
	w.Header().Set("Content-Type", file.ContentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%q", fileName))

	if _, err := w.Write(decrypted); err != nil {
		log.Printf("Error writing response for file %s: %v", fileName, err)
		return
	}

	log.Printf("Served file: %s (%d bytes, %s)", fileName, len(decrypted), file.ContentType)
}
