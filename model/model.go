package model

import "time"

//FileInfo represents a File to copy
type FileInfo struct {
	Path         string
	CreatedMonth time.Month
	CreateYear   int
}

type FileCopyDescription struct {
	Copies []FileCopy `json:"copies,omitempty"`
}

type FileCopy struct {
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}