package config

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

type PronunciationType string

const (
	PronunciationUS PronunciationType = "us"
	PronunciationUK PronunciationType = "uk"
	PronunciationCA PronunciationType = "ca"
)

type Config struct {
	PronunciationType PronunciationType
	Concurrency       int
	OutputFolder      string
	URLs              []string
	
	PageLoadTimeout    time.Duration
	ElementWaitTimeout time.Duration
	DownloadTimeout    time.Duration
	
	USXPath string
	UKXPath string
	CAXPath string
	
	HeadlessMode bool
}

func New() *Config {
	return &Config{
		PronunciationType:  PronunciationUS,
		Concurrency:        3,
		OutputFolder:       "./out",
		PageLoadTimeout:    30 * time.Second,
		ElementWaitTimeout: 10 * time.Second,
		DownloadTimeout:    30 * time.Second,
		USXPath:            "/html/body/div[4]/div/div/div/div/article/div[2]/div[8]/div/div[2]/div/div[2]/div[1]/div/button",
		UKXPath:            "/html/body/div[4]/div/div/div/div/article/div[2]/div[8]/div/div[2]/div/div[2]/div[2]/div/button",
		CAXPath:            "/html/body/div[4]/div/div/div/div/article/div[2]/div[8]/div/div[2]/div/div[2]/div[3]/div/button",
		HeadlessMode:       true,
	}
}

func (c *Config) ParseFlags() error {
	pronType := flag.String("p", "us", "Pronunciation type (us/uk/ca)")
	concurrency := flag.Int("c", 3, "Number of concurrent downloads (default: 3)")
	flag.Parse()
	
	switch strings.ToLower(*pronType) {
	case "ca":
		c.PronunciationType = PronunciationCA
	case "uk":
		c.PronunciationType = PronunciationUK
	default:
		c.PronunciationType = PronunciationUS
	}
	
	c.Concurrency = *concurrency
	c.URLs = flag.Args()
	
	if len(c.URLs) < 1 {
		return fmt.Errorf("no URLs provided")
	}
	
	return nil
}

func (c *Config) GetAudioXPath() string {
	switch c.PronunciationType {
	case PronunciationCA:
		return c.CAXPath
	case PronunciationUK:
		return c.UKXPath
	default:
		return c.USXPath
	}
}

func (c *Config) UsageNotice() string {
	return "Usage: ncfiledownloader [-p us/uk/ca] [-c concurrency] <NativeCamp_DailyNews_Page_URLs>"
}