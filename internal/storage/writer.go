package storage

import (
	"time"
)

func Write(db *DB, resultCh <-chan *LogEntry) {
	for result := range resultCh {
		path := ""
		if result.Path != nil {
			path = *result.Path
		}

		status := ""
		if result.Status != nil {
			status = *result.Status
		}

		method := ""
		if result.Method != nil {
			method = *result.Method
		}

		_, err := db.db.Exec(
			"INSERT INTO requests (service, ip, path, method, status, created_at) VALUES (?, ?, ?, ?, ?, ?)",
			result.Service,
			result.IP,
			path,
			method,
			status,
			time.Now().Format(time.RFC3339),
		)
		if err != nil {
			db.logger.Error("Failed to write to database", "error", err)
		}
	}
}
