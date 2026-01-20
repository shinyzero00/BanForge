package parser

import (
	"os"
	"testing"
	"time"
)

func TestNewScannerTail(t *testing.T) {

	file, err := os.CreateTemp("", "test-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	file.Close()

	scanner, err := NewScannerTail(file.Name())
	if err != nil {
		t.Fatalf("NewScannerTail() error = %v", err)
	}

	if scanner == nil {
		t.Fatal("Scanner is nil")
	}

	if scanner.cmd == nil {
		t.Fatal("cmd is nil")
	}

	if scanner.cmd.Process == nil {
		t.Fatal("process is nil")
	}

	scanner.Stop()
}

func TestScannerTailEvents(t *testing.T) {
	tests := []struct {
		name      string
		lines     []string
		wantLines int
	}{
		{
			name: "multiple lines",
			lines: []string{
				"Failed password for root from 192.168.1.1",
				"Invalid user admin from 192.168.1.2",
				"Accepted publickey for user from 192.168.1.3",
			},
			wantLines: 3,
		},
		{
			name: "single line",
			lines: []string{
				"Failed password for root",
			},
			wantLines: 1,
		},
		{
			name: "many lines",
			lines: []string{
				"line 1",
				"line 2",
				"line 3",
				"line 4",
				"line 5",
			},
			wantLines: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			file, err := os.CreateTemp("", "test-*.log")
			if err != nil {
				t.Fatal(err)
			}
			filePath := file.Name()
			file.Close()
			defer os.Remove(filePath)

			scanner, err := NewScannerTail(filePath)
			if err != nil {
				t.Fatalf("NewScannerTail() error = %v", err)
			}
			defer scanner.Stop()

			scanner.Start()

			time.Sleep(200 * time.Millisecond)

			file, err = os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				t.Fatal(err)
			}

			for _, line := range tt.lines {
				if _, err := file.WriteString(line + "\n"); err != nil {
					t.Fatal(err)
				}
			}

			if err := file.Sync(); err != nil {
				t.Fatal(err)
			}
			file.Close()

			// 5. Собираем события
			timeout := time.After(1 * time.Second)
			var events []Event

		eventLoop:
			for {
				select {
				case event := <-scanner.Events():
					events = append(events, event)
					t.Logf("Read: %s", event.Data)

					if len(events) == tt.wantLines {
						break eventLoop
					}

				case <-timeout:
					break eventLoop
				}
			}

			if len(events) != tt.wantLines {
				t.Errorf("got %d lines, want %d", len(events), tt.wantLines)
			}

			for i, event := range events {
				if event.Data != tt.lines[i] {
					t.Errorf("line %d: got %q, want %q", i, event.Data, tt.lines[i])
				}
			}
		})
	}
}

func TestScannerStop(t *testing.T) {

	file, err := os.CreateTemp("", "test-*.log")
	if err != nil {
		t.Fatal(err)
	}
	filePath := file.Name()
	file.Close()
	defer os.Remove(filePath)

	scanner, err := NewScannerTail(filePath)
	if err != nil {
		t.Fatal(err)
	}

	scanner.Start()
	time.Sleep(100 * time.Millisecond)

	scanner.Stop()

	err = scanner.cmd.Process.Signal(os.Signal(nil))
	if err == nil {
		t.Error("Process still alive after Stop()")
	}

	select {
	case _, ok := <-scanner.Events():
		if ok {
			t.Error("Channel still open after Stop()")
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Channel not closed after Stop()")
	}
}
func TestMultipleScanners(t *testing.T) {

	file1, err := os.CreateTemp("", "test1-*.log")
	if err != nil {
		t.Fatal(err)
	}
	path1 := file1.Name()
	file1.Close()
	defer os.Remove(path1)

	file2, err := os.CreateTemp("", "test2-*.log")
	if err != nil {
		t.Fatal(err)
	}
	path2 := file2.Name()
	file2.Close()
	defer os.Remove(path2)

	scanner1, err := NewScannerTail(path1)
	if err != nil {
		t.Fatal(err)
	}
	defer scanner1.Stop()

	scanner2, err := NewScannerTail(path2)
	if err != nil {
		t.Fatal(err)
	}
	defer scanner2.Stop()

	scanner1.Start()
	scanner2.Start()

	time.Sleep(200 * time.Millisecond)

	f1, _ := os.OpenFile(path1, os.O_APPEND|os.O_WRONLY, 0644)
	f1.WriteString("scanner1 line\n")
	f1.Sync()
	f1.Close()

	f2, _ := os.OpenFile(path2, os.O_APPEND|os.O_WRONLY, 0644)
	f2.WriteString("scanner2 line\n")
	f2.Sync()
	f2.Close()

	timeout := time.After(1 * time.Second)

	var event1, event2 Event
	got1, got2 := false, false

	for !got1 || !got2 {
		select {
		case event1 = <-scanner1.Events():
			got1 = true
			t.Logf("Scanner1: %s", event1.Data)

		case event2 = <-scanner2.Events():
			got2 = true
			t.Logf("Scanner2: %s", event2.Data)

		case <-timeout:
			if !got1 {
				t.Error("Scanner1 did not receive event")
			}
			if !got2 {
				t.Error("Scanner2 did not receive event")
			}
			return
		}
	}

	if event1.Data != "scanner1 line" {
		t.Errorf("Scanner1 got wrong data: %q", event1.Data)
	}

	if event2.Data != "scanner2 line" {
		t.Errorf("Scanner2 got wrong data: %q", event2.Data)
	}
}
func BenchmarkScanner(b *testing.B) {
	file, err := os.CreateTemp("", "bench-*.log")
	if err != nil {
		b.Fatal(err)
	}
	filePath := file.Name()
	file.Close()
	defer os.Remove(filePath)

	scanner, err := NewScannerTail(filePath)
	if err != nil {
		b.Fatal(err)
	}
	defer scanner.Stop()

	scanner.Start()
	time.Sleep(200 * time.Millisecond)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		f, _ := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
		f.WriteString("benchmark line\n")
		f.Sync()
		f.Close()
		<-scanner.Events()
	}
}
