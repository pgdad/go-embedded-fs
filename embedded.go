package main

import (
	"embed"
	"encoding/json"
	"path/filepath"
	"strings"
)

//go:embed encrypted-files/*
var encryptedFS embed.FS

// EmbeddedFile represents an encrypted file embedded in the binary
type EmbeddedFile struct {
	Name            string `json:"name"`
	EncryptedData   string `json:"encrypted_data"`
	ContentType     string `json:"content_type"`
}

// GetEmbeddedFiles reads all embedded files from the encrypted-files directory
func GetEmbeddedFiles() (map[string]EmbeddedFile, error) {
	files := make(map[string]EmbeddedFile)

	entries, err := encryptedFS.ReadDir("encrypted-files")
	if err != nil {
		return files, err
	}

	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == ".gitkeep" {
			continue
		}

		// Read the metadata file (JSON)
		if !strings.HasSuffix(entry.Name(), ".meta.json") {
			continue
		}

		metaPath := filepath.Join("encrypted-files", entry.Name())
		metaData, err := encryptedFS.ReadFile(metaPath)
		if err != nil {
			continue
		}

		var embeddedFile EmbeddedFile
		if err := json.Unmarshal(metaData, &embeddedFile); err != nil {
			continue
		}

		files[embeddedFile.Name] = embeddedFile
	}

	return files, nil
}
