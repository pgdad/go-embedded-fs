package main

import (
	"encoding/json"
	"fmt"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const encryptedFilesDir = "encrypted-files"

// AddFiles reads files, encrypts them, and saves them to the encrypted-files directory
func AddFiles(filePaths []string, key []byte) error {
	// Remove old encrypted files first
	if err := removeOldEncryptedFiles(); err != nil {
		return fmt.Errorf("failed to remove old encrypted files: %w", err)
	}

	files := make(map[string]EmbeddedFile)

	// Process each file
	for _, filePath := range filePaths {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", filePath, err)
		}

		encrypted, err := Encrypt(data, key)
		if err != nil {
			return fmt.Errorf("failed to encrypt file %s: %w", filePath, err)
		}

		// Get filename and content type
		fileName := filepath.Base(filePath)
		contentType := getContentType(fileName)

		embeddedFile := EmbeddedFile{
			Name:          fileName,
			EncryptedData: encrypted,
			ContentType:   contentType,
		}

		files[fileName] = embeddedFile

		// Write metadata file
		if err := writeEncryptedFile(embeddedFile); err != nil {
			return fmt.Errorf("failed to write encrypted file %s: %w", fileName, err)
		}

		fmt.Printf("Added: %s (%s)\n", fileName, contentType)
	}

	// Update git repository
	if err := updateGitRepository(); err != nil {
		return fmt.Errorf("failed to update git repository: %w", err)
	}

	fmt.Printf("\nSuccessfully embedded %d file(s) in encrypted-files directory.\n", len(files))
	fmt.Println("Run 'go build' to rebuild the binary with the new files.")

	return nil
}

// getContentType determines the MIME type based on file extension
func getContentType(fileName string) string {
	ext := filepath.Ext(fileName)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return contentType
}

// removeOldEncryptedFiles removes all encrypted metadata files from the directory
func removeOldEncryptedFiles() error {
	entries, err := os.ReadDir(encryptedFilesDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == ".gitkeep" {
			continue
		}

		if strings.HasSuffix(entry.Name(), ".meta.json") {
			filePath := filepath.Join(encryptedFilesDir, entry.Name())
			if err := os.Remove(filePath); err != nil {
				return fmt.Errorf("failed to remove %s: %w", filePath, err)
			}
		}
	}

	return nil
}

// writeEncryptedFile writes the encrypted file metadata to the encrypted-files directory
func writeEncryptedFile(file EmbeddedFile) error {
	// Create metadata file
	metaFileName := sanitizeFileName(file.Name) + ".meta.json"
	metaFilePath := filepath.Join(encryptedFilesDir, metaFileName)

	jsonData, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(metaFilePath, jsonData, 0644); err != nil {
		return err
	}

	return nil
}

// sanitizeFileName creates a safe filename for the metadata file
func sanitizeFileName(name string) string {
	// Replace special characters with underscores
	safe := strings.Map(func(r rune) rune {
		if r == '/' || r == '\\' || r == ':' || r == '*' || r == '?' || r == '"' || r == '<' || r == '>' || r == '|' {
			return '_'
		}
		return r
	}, name)
	return safe
}

// updateGitRepository stages changes to the encrypted-files directory
func updateGitRepository() error {
	// Check if we're in a git repository
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	if err := cmd.Run(); err != nil {
		// Not in a git repository, skip git operations
		fmt.Println("\nNote: Not in a git repository. Skipping git operations.")
		return nil
	}

	// Stage all changes in the encrypted-files directory
	cmd = exec.Command("git", "add", encryptedFilesDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git add failed: %w\nOutput: %s", err, string(output))
	}

	fmt.Println("\nGit repository updated. Encrypted files staged for commit.")
	fmt.Println("Run 'git commit' to commit the changes.")

	return nil
}
