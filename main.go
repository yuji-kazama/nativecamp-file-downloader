package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

const (
	usageNotice = "Usage: go run main.go <NativeCamp_DailyNews_Page_URLs>"
	// US Prononciation
	audioXPath = "/html/body/div[4]/div/div/div/div/article/div[2]/div[8]/div/div[2]/div/div[2]/div[1]/div/button"
	// British Pronounciation
	// audioXPath = "/html/body/div[4]/div/div/div/div/article/div[2]/div[8]/div/div[2]/div/div[2]/div[2]/div/button"
	folderPath = "./out"
	// Timeout settings
	pageLoadTimeout    = 30 * time.Second
	elementWaitTimeout = 10 * time.Second
)

func main() {
	totalArg := len(os.Args)
	if totalArg < 2 {
		fmt.Println(usageNotice)
		return
	}

	// Create output directory
	err := createDirectoryIfNotExist(folderPath)
	if err != nil {
		fmt.Println("Failed to create output directory:", err)
		return
	}

	successCount := 0
	failCount := 0

	for i := 1; i < totalArg; i++ {
		pageURL := os.Args[i]
		fmt.Printf("[%d/%d] Processing: %s\n", i, totalArg-1, pageURL)

		audioFileURL, err := getAudioFileURL(pageURL)
		if err != nil {
			fmt.Printf("  ❌ Failed to get audio file URL: %v\n", err)
			failCount++
			continue
		}

		err = downloadFile(audioFileURL)
		if err != nil {
			fmt.Printf("  ❌ Failed to download file: %v\n", err)
			failCount++
			continue
		}

		fmt.Printf("  ✅ Successfully downloaded audio from: %s\n", pageURL)
		successCount++
	}

	fmt.Printf("\n📊 Summary: %d successful, %d failed, total: %d\n",
		successCount, failCount, totalArg-1)
}

func getAudioFileURL(pageURL string) (string, error) {
	// ChromeDP options configuration
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Add timeout settings
	ctx, cancel := context.WithTimeout(allocCtx, pageLoadTimeout)
	defer cancel()

	// Create browser context
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// Channel for error handling
	errChan := make(chan error, 1)
	go func() {
		errChan <- chromedp.Run(ctx,
			chromedp.Navigate(pageURL),
			chromedp.WaitVisible(audioXPath, chromedp.BySearch),
		)
	}()

	// Timeout handling
	select {
	case err := <-errChan:
		if err != nil {
			return "", fmt.Errorf("failed to navigate page or detect element: %w", err)
		}
	case <-time.After(elementWaitTimeout):
		return "", fmt.Errorf("audio element detection timed out: %s", pageURL)
	}

	var nodes []*cdp.Node
	err := chromedp.Run(ctx, chromedp.Nodes(audioXPath, &nodes, chromedp.BySearch))
	if err != nil {
		return "", fmt.Errorf("failed to get audio element: %w", err)
	}

	if len(nodes) == 0 {
		return "", fmt.Errorf("audio element not found")
	}

	dataSrc := nodes[0].AttributeValue("data-src")
	if dataSrc == "" {
		return "", fmt.Errorf("audio URL not found (data-src attribute is empty)")
	}

	return dataSrc, nil
}

func downloadFile(fileURL string) error {
	// Add timeout setting to HTTP client
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(fileURL)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error response from server: %s", resp.Status)
	}

	fileName := getDownloadableFileName(fileURL)
	filePath := folderPath + "/" + fileName

	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	fmt.Printf("  📁 Saved to: %s\n", filePath)
	return nil
}

func createDirectoryIfNotExist(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.Mkdir(dirPath, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	return nil
}

func getDownloadableFileName(fileURL string) string {
	keyword := "/"
	index := strings.LastIndex(fileURL, keyword)
	if index < 0 {
		return fileURL
	}
	return fileURL[index+len(keyword):]
}
