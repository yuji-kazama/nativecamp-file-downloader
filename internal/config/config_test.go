package config

import (
	"flag"
	"os"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	cfg := New()
	
	if cfg.PronunciationType != PronunciationUS {
		t.Errorf("Expected default pronunciation type to be %s, got %s", PronunciationUS, cfg.PronunciationType)
	}
	
	if cfg.Concurrency != 3 {
		t.Errorf("Expected default concurrency to be 3, got %d", cfg.Concurrency)
	}
	
	if cfg.OutputFolder != "./out" {
		t.Errorf("Expected default output folder to be ./out, got %s", cfg.OutputFolder)
	}
	
	if cfg.PageLoadTimeout != 30*time.Second {
		t.Errorf("Expected page load timeout to be 30s, got %v", cfg.PageLoadTimeout)
	}
}

func TestConfig_GetAudioXPath(t *testing.T) {
	cfg := New()
	
	tests := []struct {
		pronType PronunciationType
		expected string
	}{
		{PronunciationUS, cfg.USXPath},
		{PronunciationUK, cfg.UKXPath},
		{PronunciationCA, cfg.CAXPath},
	}
	
	for _, tt := range tests {
		t.Run(string(tt.pronType), func(t *testing.T) {
			cfg.PronunciationType = tt.pronType
			result := cfg.GetAudioXPath()
			if result != tt.expected {
				t.Errorf("GetAudioXPath() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestConfig_ParseFlags(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		expected PronunciationType
	}{
		{
			name:     "default pronunciation",
			args:     []string{"cmd", "https://example.com"},
			wantErr:  false,
			expected: PronunciationUS,
		},
		{
			name:     "uk pronunciation",
			args:     []string{"cmd", "-p", "uk", "https://example.com"},
			wantErr:  false,
			expected: PronunciationUK,
		},
		{
			name:     "ca pronunciation",
			args:     []string{"cmd", "-p", "ca", "https://example.com"},
			wantErr:  false,
			expected: PronunciationCA,
		},
		{
			name:    "no URLs provided",
			args:    []string{"cmd"},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			os.Args = tt.args
			
			cfg := New()
			err := cfg.ParseFlags()
			
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr && cfg.PronunciationType != tt.expected {
				t.Errorf("Expected pronunciation type %s, got %s", tt.expected, cfg.PronunciationType)
			}
		})
	}
}