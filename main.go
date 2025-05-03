package main

import (
  "context"
  "fmt"
  "io"
  "net/http"
  "os"
  "strings"

  "github.com/chromedp/cdproto/cdp"
  "github.com/chromedp/chromedp"
)

  const (
    usageNotice = "Usage: go run main.go <NativeCamp_DailyNews_Page_URLs>"
    // US Prononciation
	audioXPath = "/html/body/div[4]/div/div/div/div/article/div[2]/div[8]/div/div[2]/div/div[2]/div[1]/div/button"
    // British Pronounciation
    // audioXPath = "/html/body/div[4]/div/div/div/div/article/div[2]/div[8]/div/div[2]/div/div[2]/div[2]/div/button
    folderPath = "./out"
  )

func main() {
  totalArg := len(os.Args)
  if totalArg < 2 {
    fmt.Println(usageNotice)
    return
  }
  for i := 1; i < totalArg; i++ {
    pageURL := os.Args[i]
    fmt.Printf("[%d/%d] %s\n", i, totalArg-1, pageURL)
    audioFileURL, err := getAudioFileURL(pageURL)
    if err != nil {
      fmt.Println("Failed to get audio file URL:", err)
      return
    }
    err = downloadFile(audioFileURL)
    if err != nil {
      fmt.Println("Failed to downlod file:", err)
      return
    }
  }
}

func getAudioFileURL(pageURL string) (string, error) {
  ctx, cancel := chromedp.NewContext(context.Background())
  defer cancel()

  var nodes []*cdp.Node
  err := chromedp.Run(ctx,
    chromedp.Navigate(pageURL),
    chromedp.WaitVisible(audioXPath, chromedp.BySearch),
    chromedp.Nodes(audioXPath, &nodes, chromedp.BySearch),
  )
  if err != nil {
    return "", fmt.Errorf("failed to execute ChromeDP: %w", err)
  }

  return nodes[0].AttributeValue("data-src"), nil
}

func downloadFile(fileURL string) error {
  resp, err := http.Get(fileURL)
  if err != nil {
    return fmt.Errorf("failed to download file: %w", err)
  }
  defer resp.Body.Close()

  err = createDirectoryIfNotExist(folderPath)
  if err != nil {
    return fmt.Errorf("failed to create directory: %w", err)
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
