package logs_postgres_repository

import "time"

type LogModel struct {
	ID         int
	FileName   string
	Status     string
	UploadedAt time.Time
	NodeCount  int
	PortCount  int
}
