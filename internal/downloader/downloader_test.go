package downloader

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestHTTPDownloader_Download(t *testing.T) {
	testData := []byte("test audio file content")
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/test.mp3" {
			w.Header().Set("Content-Type", "audio/mpeg")
			w.Write(testData)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	
	tempDir, err := os.MkdirTemp("", "downloader_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	
	d := New(30 * time.Second)
	
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "successful download",
			url:     server.URL + "/test.mp3",
			wantErr: false,
		},
		{
			name:    "404 not found",
			url:     server.URL + "/notfound.mp3",
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := d.Download(tt.url, tempDir)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Download() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				if result == nil {
					t.Error("Expected result, got nil")
					return
				}
				
				content, err := os.ReadFile(result.FilePath)
				if err != nil {
					t.Errorf("Failed to read downloaded file: %v", err)
					return
				}
				
				if string(content) != string(testData) {
					t.Errorf("Downloaded content mismatch, got %s, want %s", content, testData)
				}
			}
		})
	}
}

func TestExtractFileName(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{"https://example.com/path/to/file.mp3", "file.mp3"},
		{"https://example.com/file.mp3", "file.mp3"},
		{"file.mp3", "file.mp3"},
		{"https://example.com/", ""},
	}
	
	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			result := extractFileName(tt.url)
			if result != tt.expected {
				t.Errorf("extractFileName(%s) = %s, want %s", tt.url, result, tt.expected)
			}
		})
	}
}

func TestCreateDirectoryIfNotExist(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "create_dir_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	
	newDir := filepath.Join(tempDir, "new_directory")
	
	err = CreateDirectoryIfNotExist(newDir)
	if err != nil {
		t.Errorf("CreateDirectoryIfNotExist() error = %v", err)
		return
	}
	
	if _, err := os.Stat(newDir); os.IsNotExist(err) {
		t.Error("Directory was not created")
	}
	
	err = CreateDirectoryIfNotExist(newDir)
	if err != nil {
		t.Error("CreateDirectoryIfNotExist() should not error on existing directory")
	}
}