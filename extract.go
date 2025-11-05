package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// ExtractFiles decrypts and extracts all embedded files to the specified directory
func ExtractFiles(key []byte, outputDir string) error {
	// Resolve the output directory
	dir, err := resolveOutputDirectory(outputDir)
	if err != nil {
		return fmt.Errorf("failed to resolve output directory: %w", err)
	}

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory %s: %w", dir, err)
	}

	if len(EmbeddedFiles) == 0 {
		return fmt.Errorf("no files are currently embedded")
	}

	fmt.Printf("Extracting %d file(s) to: %s\n\n", len(EmbeddedFiles), dir)

	// Extract each file
	successCount := 0
	for name, file := range EmbeddedFiles {
		outputPath := filepath.Join(dir, name)

		// Decrypt the file
		decrypted, err := Decrypt(file.EncryptedData, key)
		if err != nil {
			fmt.Printf("✗ Failed to decrypt %s: %v\n", name, err)
			continue
		}

		// Write to disk
		if err := os.WriteFile(outputPath, decrypted, 0644); err != nil {
			fmt.Printf("✗ Failed to write %s: %v\n", name, err)
			continue
		}

		fmt.Printf("✓ Extracted: %s (%d bytes)\n", name, len(decrypted))
		successCount++
	}

	fmt.Printf("\nSuccessfully extracted %d of %d file(s)\n", successCount, len(EmbeddedFiles))

	if successCount < len(EmbeddedFiles) {
		return fmt.Errorf("some files failed to extract")
	}

	return nil
}

// resolveOutputDirectory resolves the output directory, using Downloads as default
func resolveOutputDirectory(outputDir string) (string, error) {
	// If user specified a directory, use it
	if outputDir != "" {
		// Convert to absolute path
		absPath, err := filepath.Abs(outputDir)
		if err != nil {
			return "", err
		}
		return absPath, nil
	}

	// Use platform-specific Downloads directory as default
	return getDownloadsDirectory()
}

// getDownloadsDirectory returns the platform-specific Downloads directory
func getDownloadsDirectory() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	var downloadsDir string

	switch runtime.GOOS {
	case "windows":
		downloadsDir = filepath.Join(homeDir, "Downloads")
	case "darwin": // macOS
		downloadsDir = filepath.Join(homeDir, "Downloads")
	case "linux":
		downloadsDir = filepath.Join(homeDir, "Downloads")
	default:
		// Fallback for other Unix-like systems
		downloadsDir = filepath.Join(homeDir, "Downloads")
	}

	return downloadsDir, nil
}
