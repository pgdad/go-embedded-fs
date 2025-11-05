# Go Embedded Encrypted Filesystem

A Go application that embeds encrypted files in the binary and serves them via HTTP, with on-the-fly decryption.

## Features

- **Encrypted Embedded Files**: Files are encrypted, stored in the repository, and embedded in the binary using Go's embed package
- **Symmetric Encryption**: Uses AES-256-GCM for secure encryption/decryption
- **HTTP Server**: Serves decrypted files on-demand via HTTP
- **File Extraction**: Extract all embedded files to a directory
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
- Files are encrypted, base64-encoded, and stored as JSON in the `encrypted-files/` directory
- The `add` command automatically stages changes in git
- After adding files, you must rebuild the binary: `go build -o go-embedded-fs`
- Commit the changes to git to keep encrypted files in version control

### Extracting Files

Extract all embedded files to a directory (decrypted):

```bash
./go-embedded-fs extract
```

By default, files are extracted to your platform-specific Downloads folder:
- **Linux**: `~/Downloads`
- **macOS**: `~/Downloads`
- **Windows**: `%USERPROFILE%\Downloads`

Extract to a custom directory:

```bash
./go-embedded-fs extract -output /path/to/directory
```

Or with a relative path:

```bash
./go-embedded-fs extract -output ./my-files
```

Example:
```bash
ENCRYPT_KEY="MySecretKey123" ./go-embedded-fs extract -output ./decrypted
```

**Note:** The output directory will be created automatically if it doesn't exist.

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
1. Removes old encrypted files from `encrypted-files/` directory
2. Reads files from disk
3. Encrypts each file using AES-256-GCM with the key from `ENCRYPT_KEY`
4. Base64-encodes the encrypted data
5. Detects and stores MIME type for each file
6. Writes metadata as JSON files to `encrypted-files/` directory
7. Automatically stages changes in git repository

### Serving Files (Decryption)
1. Loads encrypted files from embedded filesystem (using Go's embed package)
2. Receives HTTP request for a file
3. Retrieves encrypted data from embedded files
4. Base64-decodes the data
5. Decrypts using the key from `ENCRYPT_KEY`
6. Serves decrypted content with proper Content-Type header

### Extracting Files
1. Loads encrypted files from embedded filesystem (using Go's embed package)
2. Determines output directory (Downloads folder or custom path)
3. Creates output directory if it doesn't exist
4. For each embedded file:
   - Retrieves encrypted data from embedded files
   - Base64-decodes the data
   - Decrypts using the key from `ENCRYPT_KEY`
   - Writes decrypted content to disk

## Security Considerations

- **Key Management**: The encryption key should be kept secret and managed securely
- **Key Length**: Use a strong, random key (32 bytes for AES-256)
- **Environment Variables**: Never commit the encryption key to version control
- **Repository Security**: Encrypted files in `encrypted-files/` are stored in git and are only as secure as the encryption key
- **Version Control**: The encrypted files can be safely committed to public repositories as long as the encryption key remains secret
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

# Extract all files to Downloads folder
./go-embedded-fs extract

# Or extract to a specific directory
./go-embedded-fs extract -output ./my-files
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
