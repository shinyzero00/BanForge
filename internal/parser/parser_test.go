package parser

import (
	"os"
	"testing"
	"time"
)

func TestNewScanner(t *testing.T) {
	file, err := os.CreateTemp("", "test.log")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	defer os.Remove(file.Name())
	s, err := NewScanner(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	if s == nil {
		t.Fatal("Scanner is nil")
	}
}

func TestScannerStart(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantErr   bool
		wantLines int
	}{
		{
			name: "correct file",
			input: `Failed password for root from 192.168.1.1
Invalid user admin from 192.168.1.1
Accepted publickey for user from 192.168.1.2`,
			wantErr:   false,
			wantLines: 3,
		},
		{
			name:      "empty file",
			input:     "",
			wantErr:   false,
			wantLines: 0,
		},
		{
			name:      "single line",
			input:     `Failed password for root`,
			wantErr:   false,
			wantLines: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			file, err := os.CreateTemp("", "test-*.log")
			if err != nil {
				t.Fatal(err)
			}
			filePath := file.Name()

			if _, err := file.WriteString(tt.input); err != nil {
				t.Fatal(err)
			}
			file.Close()
			defer os.Remove(filePath)

			scanner, err := NewScanner(filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewScanner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}
			defer scanner.Stop()

			scanner.Start()

			timeout := time.After(500 * time.Millisecond)
			linesRead := 0

			for {
				select {
				case event := <-scanner.Events():
					linesRead++
					t.Logf("Read: %s", event.Data)
				case <-timeout:
					if linesRead != tt.wantLines {
						t.Errorf("got %d lines, want %d", linesRead, tt.wantLines)
					}
					return
				}
			}
		})
	}
}
