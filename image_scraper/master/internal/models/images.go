package models

import (
	"os"
	"time"
)

var UploadsDir = os.Getenv("UPLOAD_DIR")

type Image struct {
	ID         int       `json:"id"`
	Filename   string    `json:"filename"`
	Format     string    `json:"format"`
	Size       int       `json:"size"`
	UploadDate time.Time `json:"upload_date"`
	Filepath   string    `json:"filepath"`
}
