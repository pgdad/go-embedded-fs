package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	envKeyName        = "ENCRYPT_KEY"
	defaultServerPort = "9193"
)

func main() {
	// Define flags
	portFlag := flag.String("port", defaultServerPort, "Port to listen on")
	outputFlag := flag.String("output", "", "Output directory for extracted files (default: Downloads folder)")

	// Get encryption key from environment variable
	keyString := os.Getenv(envKeyName)
	if keyString == "" {
		log.Fatalf("Error: %s environment variable is not set\n", envKeyName)
	}

	key := DeriveKey(keyString)

	// Parse command-line arguments
	if len(os.Args) < 2 || (len(os.Args) >= 2 && os.Args[1][0] == '-') {
		// Default to server mode (no args or first arg is a flag)
		flag.Parse()
		runServer(key, *portFlag)
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
		// Parse flags after the serve command
		flag.CommandLine.Parse(os.Args[2:])
		runServer(key, *portFlag)

	case "extract":
		// Parse flags after the extract command
		flag.CommandLine.Parse(os.Args[2:])
		if err := ExtractFiles(key, *outputFlag); err != nil {
			log.Fatalf("Error extracting files: %v\n", err)
		}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("\nUsage:")
		fmt.Println("  program                     Start web server (default)")
		fmt.Println("  program serve [-port]       Start web server")
		fmt.Println("  program add <files>         Add and encrypt files")
		fmt.Println("  program extract [-output]   Extract all files to directory")
		fmt.Println("\nFlags:")
		fmt.Println("  -port string    Port to listen on (default \"9193\")")
		fmt.Println("  -output string  Output directory for extracted files (default: Downloads folder)")
		os.Exit(1)
	}
}

func runServer(key []byte, port string) {
	server, err := NewServer(key)
	if err != nil {
		log.Fatalf("Failed to create server: %v\n", err)
	}

	if err := server.Start(port); err != nil {
		log.Fatalf("Server error: %v\n", err)
	}
}
