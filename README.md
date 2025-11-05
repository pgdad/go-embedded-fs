# Go Embedded Encrypted Filesystem

A Go application that embeds encrypted files in the binary and serves them via HTTP, with on-the-fly decryption.

## Features

- **Encrypted Embedded Files**: Files are encrypted and embedded directly in the Go binary source code
- **Symmetric Encryption**: Uses AES-256-GCM for secure encryption/decryption
- **HTTP Server**: Serves decrypted files on-demand via HTTP
- **Configurable Port**: Default port 9193, customizable via command-line flag
- **Cross-Platform**: Works on Linux, Windows, and macOS
- **Simple API**: JSON file listing and direct file access
- **Content-Type Detection**: Automatically sets correct MIME types based on file extensions

## Requirements

- Go 1.16 or later
- Encryption key (32 bytes recommended for AES-256)

## Installation

```bash
go build -o go-embedded-fs
```

## Usage

### Environment Variable

The encryption key must be provided via the `ENCRYPT_KEY` environment variable:

```bash
export ENCRYPT_KEY="your-secret-key-here"
```

### Adding Files

To encrypt and embed files in the source code:

```bash
./go-embedded-fs add <file1> [file2] [file3] ...
```

Example:
```bash
ENCRYPT_KEY="MySecretKey123" ./go-embedded-fs add document.pdf image.png presentation.pptx
```

**Important Notes:**
- Each `add` operation replaces ALL previously embedded files
- After adding files, you must rebuild the binary: `go build -o go-embedded-fs`
- Files are encrypted, base64-encoded, and written to `embedded.go`

### Running the Server

Start the HTTP server on the default port (9193):

```bash
./go-embedded-fs
```

Or explicitly:

```bash
./go-embedded-fs serve
```

Start the server on a custom port:

```bash
./go-embedded-fs -port 8080
```

Or with the serve command:

```bash
./go-embedded-fs serve -port 8080
```

## API Endpoints

### List Files

Get a JSON list of all available files:

```bash
curl http://localhost:9193/
```

Response:
```json
["file1.pdf", "file2.png", "file3.pptx"]
```

### Download File

Retrieve a specific file (automatically decrypted):

```bash
curl http://localhost:9193/file1.pdf > file1.pdf
```

The response includes the correct `Content-Type` header based on the file extension.

## How It Works

### Adding Files (Encryption)
1. Reads files from disk
2. Encrypts each file using AES-256-GCM with the key from `ENCRYPT_KEY`
3. Base64-encodes the encrypted data
4. Writes the encoded data into `embedded.go` source file
5. Detects and stores MIME type for each file

### Serving Files (Decryption)
1. Receives HTTP request for a file
2. Retrieves encrypted data from embedded map
3. Base64-decodes the data
4. Decrypts using the key from `ENCRYPT_KEY`
5. Serves decrypted content with proper Content-Type header

## Security Considerations

- **Key Management**: The encryption key should be kept secret and managed securely
- **Key Length**: Use a strong, random key (32 bytes for AES-256)
- **Environment Variables**: Never commit the encryption key to version control
- **Source Protection**: The encrypted data in `embedded.go` is only as secure as the encryption key
- **HTTPS**: For production use, consider running behind a reverse proxy with HTTPS

## Example Workflow

```bash
# Set encryption key
export ENCRYPT_KEY="MyVerySecretKey123456789012345"

# Add files to embed
./go-embedded-fs add docs/report.pdf images/logo.png config/settings.json

# Rebuild binary with embedded files
go build -o go-embedded-fs

# Start server
./go-embedded-fs

# In another terminal, test the server
curl http://localhost:9193/                    # List files
curl http://localhost:9193/report.pdf -o report.pdf  # Download file

# Or start server on a custom port
./go-embedded-fs -port 8080
curl http://localhost:8080/                    # List files on custom port
```

## Cross-Platform Compilation

Build for different platforms:

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o go-embedded-fs-linux

# Windows
GOOS=windows GOARCH=amd64 go build -o go-embedded-fs.exe

# macOS
GOOS=darwin GOARCH=amd64 go build -o go-embedded-fs-mac
```

## License

MIT
