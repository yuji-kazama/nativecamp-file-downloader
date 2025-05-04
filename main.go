package main

import (
	"context"
	"flag"
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
	usageNotice = "Usage: go run main.go [-p br] <NativeCamp_DailyNews_Page_URLs>"
	folderPath  = "./out"
	// Timeout settings
	pageLoadTimeout    = 30 * time.Second
	elementWaitTimeout = 10 * time.Second

	// XPath for different pronunciations
	usXPath = "/html/body/div[4]/div/div/div/div/article/div[2]/div[8]/div/div[2]/div/div[2]/div[1]/div/button"
	ukXPath = "/html/body/div[4]/div/div/div/div/article/div[2]/div[8]/div/div[2]/div/div[2]/div[2]/div/button"
	caXPath = "/html/body/div[4]/div/div/div/div/article/div[2]/div[8]/div/div[2]/div/div[2]/div[3]/div/button"
)

var audioXPath string

func main() {
	pronType := flag.String("p", "us", "Pronunciation type (us/uk/ca)")
	flag.Parse()

	switch strings.ToLower(*pronType) {
	case "ca":
		audioXPath = caXPath
	case "uk":
		audioXPath = ukXPath
	default:
		audioXPath = usXPath
	}

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println(usageNotice)
		return
	}

	err := createDirectoryIfNotExist(folderPath)
	if err != nil {
		fmt.Println("Failed to create output directory:", err)
		return
	}

	successCount := 0
	failCount := 0

	for i, pageURL := range args {
		fmt.Printf("[%d/%d] Processing: %s\n", i+1, len(args), pageURL)

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
		successCount, failCount, len(args))
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

	ctx, cancel := context.WithTimeout(allocCtx, pageLoadTimeout)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

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
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(fileURL)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

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
