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

		reason := ""
		if result.Reason != nil {
			reason = *result.Reason
		}

		_, err := db.db.Exec(
			"INSERT INTO requests (service, ip, path, method, status, reason, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
			result.Service,
			result.IP,
			path,
			method,
			status,
			reason,
			time.Now().Format(time.RFC3339),
		)
		if err != nil {
			db.logger.Error("Failed to write to database", "error", err)
		}
	}
}
