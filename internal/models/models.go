package models

type WorkItem struct {
	Index   int
	PageURL string
}

type WorkResult struct {
	Index    int
	PageURL  string
	AudioURL string
	Error    error
}

type DownloadResult struct {
	URL      string
	FilePath string
	Error    error
}