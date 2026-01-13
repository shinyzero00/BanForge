package parser

import (
	"fmt"
	"regexp"

	"github.com/d3m0k1d/BanForge/internal/logger"
	"github.com/d3m0k1d/BanForge/internal/storage"
)

type NginxParser struct {
	pattern *regexp.Regexp
	logger  *logger.Logger
}

func NewNginxParser() *NginxParser {
	pattern := regexp.MustCompile(
		`^(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}).*\[(.*?)\]\s+"(\w+)\s+(.*?)\s+HTTP.*"\s+(\d+)`,
	)
	return &NginxParser{
		pattern: pattern,
		logger:  logger.New(false),
	}
}

func (p *NginxParser) Parse(line string) (*storage.LogEntry, error) {
	// Group 1: IP, Group 2: Timestamp, Group 3: Method, Group 4: Path, Group 5: Status
	matches := p.pattern.FindStringSubmatch(line)
	if matches == nil {
		return nil, fmt.Errorf("invalid log format")
	}

	return &storage.LogEntry{
		Service: "nginx",
		IP:      matches[1],
		Path:    &matches[4],
		Status:  &matches[5],
		Method:  &matches[3],
		Reason:  nil,
	}, nil
}
