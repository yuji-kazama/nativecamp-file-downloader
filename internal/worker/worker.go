package worker

import (
	"fmt"
	"log/slog"
	"sync"
	
	"ncfiledownloader/internal/config"
	"ncfiledownloader/internal/downloader"
	"ncfiledownloader/internal/models"
	"ncfiledownloader/internal/scraper"
)

type Pool struct {
	config     *config.Config
	scraper    scraper.Scraper
	downloader downloader.Downloader
	logger     *slog.Logger
}

func NewPool(cfg *config.Config, s scraper.Scraper, d downloader.Downloader, logger *slog.Logger) *Pool {
	return &Pool{
		config:     cfg,
		scraper:    s,
		downloader: d,
		logger:     logger,
	}
}


func (p *Pool) ProcessConcurrent(urls []string) (int, int) {
	workChan := make(chan models.WorkItem, len(urls))
	resultChan := make(chan models.WorkResult, len(urls))
	
	var wg sync.WaitGroup
	for i := 0; i < p.config.Concurrency; i++ {
		wg.Add(1)
		go p.worker(i+1, workChan, resultChan, &wg, len(urls))
	}
	
	for i, pageURL := range urls {
		workChan <- models.WorkItem{Index: i, PageURL: pageURL}
	}
	close(workChan)
	
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	
	results := make([]models.WorkResult, 0, len(urls))
	for result := range resultChan {
		results = append(results, result)
	}
	
	successCount := 0
	failCount := 0
	
	for i := 0; i < len(results); i++ {
		for _, result := range results {
			if result.Index == i {
				if result.Error != nil {
					p.logger.Error("Failed", 
						"index", result.Index+1,
						"total", len(urls),
						"url", result.PageURL,
						"error", result.Error)
					failCount++
				} else {
					p.logger.Info("Success",
						"index", result.Index+1,
						"total", len(urls),
						"url", result.PageURL)
					successCount++
				}
				break
			}
		}
	}
	
	return successCount, failCount
}

func (p *Pool) worker(id int, workChan <-chan models.WorkItem, resultChan chan<- models.WorkResult, wg *sync.WaitGroup, totalCount int) {
	defer wg.Done()
	
	for work := range workChan {
		p.logger.Info("Worker processing",
			"workerID", id,
			"index", work.Index+1,
			"total", totalCount,
			"url", work.PageURL)
		
		audioFileURL, err := p.scraper.GetAudioURL(work.PageURL, p.config.GetAudioXPath())
		if err != nil {
			resultChan <- models.WorkResult{
				Index:   work.Index,
				PageURL: work.PageURL,
				Error:   fmt.Errorf("failed to get audio file URL: %w", err),
			}
			continue
		}
		
		result, err := p.downloader.Download(audioFileURL, p.config.OutputFolder)
		if err != nil {
			resultChan <- models.WorkResult{
				Index:   work.Index,
				PageURL: work.PageURL,
				Error:   fmt.Errorf("failed to download file: %w", err),
			}
			continue
		}
		
		p.logger.Debug("Worker completed download",
			"workerID", id,
			"filePath", result.FilePath)
		
		resultChan <- models.WorkResult{
			Index:    work.Index,
			PageURL:  work.PageURL,
			AudioURL: audioFileURL,
			Error:    nil,
		}
	}
}