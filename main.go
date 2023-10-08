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
		audioXPath = "/html/body/div[4]/div/div/div/div/article/div[1]/div[8]/div/div[2]/div/div[2]/p/a" 
		folderPath = "./out"
	)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <pageURL>")
		return
	}
	pageURL := os.Args[1]

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var nodes []*cdp.Node
	if err := chromedp.Run(ctx,
		chromedp.Navigate(pageURL),
		chromedp.WaitVisible(audioXPath, chromedp.BySearch),
		chromedp.Nodes(audioXPath, &nodes, chromedp.BySearch),
	); err != nil {
		panic(err)
	}
	var audioFileURL = nodes[0].AttributeValue("data-src")
	fmt.Println("Audio file URL: " + audioFileURL)

	download(audioFileURL)
	fmt.Println("Finish file download")
}

func download(fileURL string) {
	resp, err := http.Get(fileURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	_, err = os.Stat(folderPath)
	if os.IsNotExist(err) {
		err := os.Mkdir(folderPath, 0755)
		if err != nil {
			panic(err)
		}
	}
	out, err := os.Create(folderPath + "/" + getFileName(fileURL))
	if err != nil {
		panic(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
}

func getFileName(fileURL string) string {
	keyword := "/"
	index := strings.LastIndex(fileURL, keyword)
	if index < 0 {
		return fileURL
	}
	return fileURL[index+len(keyword):]
}