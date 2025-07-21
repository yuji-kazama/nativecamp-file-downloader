# nativecamp-file-downloader

A Go application for downloading audio files from NativeCamp Daily News pages with support for different pronunciations (US, UK, CA) and concurrent downloads.

## Requirements
* Go 1.21 or higher
* Chrome/Chromium browser (for headless scraping)

## Installation

### From Source
```bash
# Clone the repository
git clone <repository-url>
cd nativecamp-file-downloader

# Download dependencies
make deps

# Build the binary
make build

# Or install to $GOPATH/bin
make install
```

## Usage
```
ncfiledownloader [-p us/uk/ca] [-c concurrency] <NativeCamp_DailyNews_Page_URLs>
```

### Options
- `-p`: Pronunciation type (us/uk/ca), default: us
- `-c`: Number of concurrent downloads, default: 3 (use 1 for sequential downloads)

### Examples
```bash
# Download with default settings (US pronunciation, 3 concurrent workers)
./ncfiledownloader https://nativecamp.net/textbook/page-detail/1/40468

# Download UK pronunciation with 5 concurrent workers
./ncfiledownloader -p uk -c 5 url1 url2 url3 url4 url5

# Sequential download (single worker)
./ncfiledownloader -c 1 url1 url2 url3

# Using make run
make run ARGS="-p uk https://nativecamp.net/textbook/page-detail/1/40468"

# Run example URLs from run.sh
make example
```

## Project Structure
```
.
├── cmd/
│   └── ncfiledownloader/    # Main application entry point
├── internal/
│   ├── config/              # Configuration management
│   ├── downloader/          # HTTP download functionality
│   ├── models/              # Data structures
│   ├── scraper/             # Web scraping logic
│   └── worker/              # Concurrent processing
├── pkg/
│   └── errors/              # Custom error types
├── out/                     # Downloaded files (created automatically)
├── Makefile                 # Build and development tasks
├── go.mod                   # Go module definition
└── run.sh                   # Example script
```

## Development

### Running Tests
```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

### Code Quality
```bash
# Format code
make fmt

# Run go vet
make vet

# Run all linters
make lint
```

### Building
```bash
# Build for current platform
make build

# Build for multiple platforms
make build-all

# Clean build artifacts
make clean
```

## Features

- **Modular Architecture**: Clean separation of concerns with dedicated packages
- **Concurrent Downloads**: Configurable worker pool for parallel processing
- **Multiple Pronunciations**: Support for US, UK, and Canadian pronunciations
- **Structured Logging**: Detailed logging with slog for better debugging
- **Error Handling**: Custom error types with context information
- **Testable Design**: Interfaces and dependency injection for easy testing
- **Cross-platform**: Builds for macOS, Linux, and Windows

## Troubleshooting

If the program is not working due to XPath changes:

1. Open a NativeCamp Daily News page in Chrome
2. Open Developer Tools (F12)
3. Select the audio button element
4. Copy the full XPath
5. Update the XPath constants in `internal/config/config.go`

## License

[Your License Here]
