package scraper

import (
	"context"
	"fmt"
	"time"
	
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	
	ncerrors "ncfiledownloader/pkg/errors"
)

type Scraper interface {
	GetAudioURL(pageURL string, xpath string) (string, error)
}

type ChromeScraper struct {
	pageLoadTimeout    time.Duration
	elementWaitTimeout time.Duration
	headless           bool
}

func New(pageLoadTimeout, elementWaitTimeout time.Duration, headless bool) Scraper {
	return &ChromeScraper{
		pageLoadTimeout:    pageLoadTimeout,
		elementWaitTimeout: elementWaitTimeout,
		headless:           headless,
	}
}

func (s *ChromeScraper) GetAudioURL(pageURL string, xpath string) (string, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", s.headless),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)
	
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	
	ctx, cancel := context.WithTimeout(allocCtx, s.pageLoadTimeout)
	defer cancel()
	
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()
	
	errChan := make(chan error, 1)
	go func() {
		errChan <- chromedp.Run(ctx,
			chromedp.Navigate(pageURL),
			chromedp.WaitVisible(xpath, chromedp.BySearch),
		)
	}()
	
	select {
	case err := <-errChan:
		if err != nil {
			return "", ncerrors.NewScraperError(pageURL, fmt.Errorf("failed to navigate page or detect element: %w", err))
		}
	case <-time.After(s.elementWaitTimeout):
		return "", ncerrors.NewScraperError(pageURL, ncerrors.ErrTimeout)
	}
	
	var nodes []*cdp.Node
	err := chromedp.Run(ctx, chromedp.Nodes(xpath, &nodes, chromedp.BySearch))
	if err != nil {
		return "", ncerrors.NewScraperError(pageURL, fmt.Errorf("failed to get audio element: %w", err))
	}
	
	if len(nodes) == 0 {
		return "", ncerrors.NewScraperError(pageURL, ncerrors.ErrElementNotFound)
	}
	
	dataSrc := nodes[0].AttributeValue("data-src")
	if dataSrc == "" {
		return "", ncerrors.NewScraperError(pageURL, ncerrors.ErrEmptyDataSrc)
	}
	
	return dataSrc, nil
}