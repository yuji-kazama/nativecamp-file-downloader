package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func main() {
	// get pagelURL from args
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <pageURL>")
	}
	pageURL := os.Args[1]

	// get fileURL from pageURL
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var nodes []*cdp.Node
	if err := chromedp.Run(ctx,
		chromedp.Navigate(pageURL),
		chromedp.Sleep(5000*time.Millisecond),
		chromedp.Nodes(`/html/body/div[4]/div/div/div/div/article/div[1]/div[8]/div/div[2]/div/div[2]/p/a`, &nodes, chromedp.BySearch),

	); err != nil {
		log.Fatal(err)
	}
	var fileURL = nodes[0].AttributeValue("data-src")
	fmt.Println(fileURL)

	// get file
	resp, err := http.Get(fileURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// dowlad file
	folderPath := "./out"
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
	fmt.Println("Finished: downloaded file from " + fileURL)
}

func getFileName(fileURL string) string {
	keyword := "/"
	index := strings.LastIndex(fileURL, keyword)
	if index < 0 {
		return fileURL
	}
	return fileURL[index+len(keyword):]
}