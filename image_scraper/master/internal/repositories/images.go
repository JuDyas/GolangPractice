package repositories

import (
	"database/sql"
	"fmt"

	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/master/internal/models"
)

type ImageRepository interface {
	SaveImage(url string) error
	CheckExist(url string) (bool, error)
}

type ImageRepositoryImpl struct {
	db *sql.DB
}

func NewImageRepository(database *sql.DB) *ImageRepositoryImpl {
	return &ImageRepositoryImpl{
		db: database,
	}
}

func (repo *ImageRepositoryImpl) SaveImage(image models.Image) error {
	query := `INSERT INTO images (filename, format, size, filepath)
			VALUES ($1, $2, $3, $4)
			RETURNING id`
	err := repo.db.QueryRow(query, image.Filename, image.Format, image.Size, image.Filepath).Scan(&image.ID)
	if err != nil {
		return fmt.Errorf("save image data in postgres: %w", err)
	}

	return nil
}

func (repo *ImageRepositoryImpl) CheckExist(filename string) (bool, error) {
	var count int
	err := repo.db.QueryRow("SELECT COUNT(*) FROM images WHERE filename = $1", filename).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("Error checking existence of image " + filename + ": " + err.Error())
	}
	return count > 0, nil
}
