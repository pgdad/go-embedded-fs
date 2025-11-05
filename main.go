package main

import (
	"fmt"
	"log"
	"os"
)

const (
	envKeyName = "ENCRYPT_KEY"
	serverPort = "8989"
)

func main() {
	// Get encryption key from environment variable
	keyString := os.Getenv(envKeyName)
	if keyString == "" {
		log.Fatalf("Error: %s environment variable is not set\n", envKeyName)
	}

	key := DeriveKey(keyString)

	// Parse command-line arguments
	if len(os.Args) < 2 {
		// Default to server mode
		runServer(key)
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: program add <file1> [file2] [file3] ...")
			fmt.Println("Example: program add document.pdf image.png presentation.pptx")
			os.Exit(1)
		}
		filePaths := os.Args[2:]
		if err := AddFiles(filePaths, key); err != nil {
			log.Fatalf("Error adding files: %v\n", err)
		}

	case "serve", "server":
		runServer(key)

	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("\nUsage:")
		fmt.Println("  program              Start web server (default)")
		fmt.Println("  program serve        Start web server")
		fmt.Println("  program add <files>  Add and encrypt files")
		os.Exit(1)
	}
}

func runServer(key []byte) {
	if len(EmbeddedFiles) == 0 {
		fmt.Println("Warning: No files are currently embedded.")
		fmt.Println("Use 'program add <files>' to add files first.")
		fmt.Println()
	}

	server := NewServer(key)
	if err := server.Start(serverPort); err != nil {
		log.Fatalf("Server error: %v\n", err)
	}
}
