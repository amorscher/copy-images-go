package model

import "time"

//FileInfo represents a File to copy
type FileInfo struct {
	Path         string
	CreationDate time.Time
}

type FileOperations struct {
	FileOperations []FileOperation `json:"operations,omitempty"`
}

type OpType string

const (
	MoveOp OpType = "MOVE"
	CopyOp OpType = "COPY"
)

type FileOperation struct {
	From   string `json:"from,omitempty"`
	To     string `json:"to,omitempty"`
	OpType OpType `json:"type,omitempty"`
}
