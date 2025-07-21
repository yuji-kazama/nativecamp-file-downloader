package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	
	"ncfiledownloader/internal/models"
	ncerrors "ncfiledownloader/pkg/errors"
)

type Downloader interface {
	Download(fileURL string, outputFolder string) (*models.DownloadResult, error)
}

type HTTPDownloader struct {
	client *http.Client
}

func New(timeout time.Duration) Downloader {
	return &HTTPDownloader{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (d *HTTPDownloader) Download(fileURL string, outputFolder string) (*models.DownloadResult, error) {
	resp, err := d.client.Get(fileURL)
	if err != nil {
		return nil, ncerrors.NewDownloadError(fileURL, 0, err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, ncerrors.NewDownloadError(fileURL, resp.StatusCode, ncerrors.ErrInvalidResponse)
	}
	
	fileName := extractFileName(fileURL)
	filePath := filepath.Join(outputFolder, fileName)
	
	out, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()
	
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}
	
	return &models.DownloadResult{
		URL:      fileURL,
		FilePath: filePath,
	}, nil
}

func extractFileName(fileURL string) string {
	lastSlash := strings.LastIndex(fileURL, "/")
	if lastSlash < 0 {
		return fileURL
	}
	return fileURL[lastSlash+1:]
}

func CreateDirectoryIfNotExist(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	return nil
}