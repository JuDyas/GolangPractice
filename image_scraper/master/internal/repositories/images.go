package repositories

import (
	"database/sql"
	"errors"
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
	query := `INSERT INTO images (filename, format, width, height, filepath)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id`
	err := repo.db.QueryRow(query, image.Filename, image.Format, image.Width, image.Height, image.Filepath).Scan(&image.ID)
	if err != nil {
		return fmt.Errorf("save image data in postgres: %w", err)
	}

	return nil
}

func (repo *ImageRepositoryImpl) CheckExist(filename string) (models.Image, bool, error) {
	var image models.Image
	err := repo.db.QueryRow(`
		SELECT id, filename, format, width, height, filepath
		FROM images 
		WHERE filename = $1`, filename).Scan(&image.ID, &image.Filename, &image.Format, &image.Width, &image.Height, &image.Filepath)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return image, false, nil
		}
		return image, false, fmt.Errorf("checking existence of image %s: %v", filename, err)
	}

	return image, true, nil
}
