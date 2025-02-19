package models

import (
	"errors"
	"time"
)

const ResizeFactor = 0.8

var MaxSizeErr = errors.New("size of image is so large")

type Image struct {
	ID         int       `json:"id"`
	Filename   string    `json:"filename"`
	Format     string    `json:"format"`
	Width      int       `json:"width"`
	Height     int       `json:"height"`
	UploadDate time.Time `json:"upload_date"`
	Filepath   string    `json:"filepath"`
}
