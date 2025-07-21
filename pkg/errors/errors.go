package errors

import (
	"errors"
	"fmt"
)

var (
	ErrNoURLsProvided    = errors.New("no URLs provided")
	ErrElementNotFound   = errors.New("audio element not found")
	ErrEmptyDataSrc      = errors.New("audio URL not found (data-src attribute is empty)")
	ErrTimeout           = errors.New("operation timed out")
	ErrInvalidResponse   = errors.New("invalid server response")
	ErrDirectoryCreation = errors.New("failed to create directory")
)

type ScraperError struct {
	URL string
	Err error
}

func (e *ScraperError) Error() string {
	return fmt.Sprintf("scraper error for URL %s: %v", e.URL, e.Err)
}

func (e *ScraperError) Unwrap() error {
	return e.Err
}

type DownloadError struct {
	URL        string
	StatusCode int
	Err        error
}

func (e *DownloadError) Error() string {
	if e.StatusCode != 0 {
		return fmt.Sprintf("download error for URL %s (status %d): %v", e.URL, e.StatusCode, e.Err)
	}
	return fmt.Sprintf("download error for URL %s: %v", e.URL, e.Err)
}

func (e *DownloadError) Unwrap() error {
	return e.Err
}

func NewScraperError(url string, err error) error {
	return &ScraperError{URL: url, Err: err}
}

func NewDownloadError(url string, statusCode int, err error) error {
	return &DownloadError{URL: url, StatusCode: statusCode, Err: err}
}