package domain

import "time"

type LogStatus string

var (
	LogStatusProcessing LogStatus = "processing"
	LogStatusDone       LogStatus = "done"
	LogStatusFailed     LogStatus = "failed"
)

type Log struct {
	ID         int
	FileName   string
	Status     LogStatus
	UploadedAt time.Time
	NodeCount  int
	PortCount  int
}

func NewLog(id int, fileName string, status LogStatus, uploadedAt time.Time, nodeCount int, portCount int) Log {
	return Log{
		ID: id,
		FileName: fileName,
		Status: status,
		UploadedAt: uploadedAt,
		NodeCount: nodeCount,
		PortCount: portCount,
	}
}

func NewUninitializedLog(fileName string) Log {
	return NewLog(
		UninitializedID,
		fileName,
		LogStatusProcessing,
		time.Now().UTC(),
		0,
		0,
	)
}