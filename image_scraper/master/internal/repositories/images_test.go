package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/master/internal/models"
)

func TestImageRepositoryImpl_SaveImage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %s", err)
	}
	defer db.Close()

	repo := NewImageRepository(db)

	tests := []struct {
		name      string
		image     models.Image
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Success",
			image: models.Image{
				Filename: "test.jpg", Format: "jpg", Width: 800, Height: 600, Filepath: "/images/test.jpg",
			},
			mockSetup: func() {
				mock.ExpectQuery("^INSERT INTO images").
					WithArgs("test.jpg", "jpg", 800, 600, "/images/test.jpg").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantErr: false,
		},
		{
			name: "DB Error",
			image: models.Image{
				Filename: "test.jpg", Format: "jpg", Width: 800, Height: 600, Filepath: "/images/test.jpg",
			},
			mockSetup: func() {
				mock.ExpectQuery("^INSERT INTO images").
					WithArgs("test.jpg", "jpg", 800, 600, "/images/test.jpg").
					WillReturnError(fmt.Errorf("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repo.SaveImage(tt.image)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveImage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled mock expectations: %v", err)
			}
		})
	}
}

func TestImageRepositoryImpl_CheckExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %s", err)
	}
	defer db.Close()

	repo := NewImageRepository(db)

	tests := []struct {
		name      string
		filename  string
		mockSetup func()
		wantImage models.Image
		wantExist bool
		wantErr   bool
	}{
		{
			name:     "Exists",
			filename: "test.jpg",
			mockSetup: func() {
				mock.ExpectQuery("^SELECT id, filename, format, width, height, filepath").
					WithArgs("test.jpg").
					WillReturnRows(sqlmock.NewRows([]string{"id", "filename", "format", "width", "height", "filepath"}).
						AddRow(1, "test.jpg", "jpg", 800, 600, "/images/test.jpg"))
			},
			wantImage: models.Image{
				ID: 1, Filename: "test.jpg", Format: "jpg", Width: 800, Height: 600, Filepath: "/images/test.jpg",
			},
			wantExist: true,
			wantErr:   false,
		},
		{
			name:     "Does Not Exist",
			filename: "missing.jpg",
			mockSetup: func() {
				mock.ExpectQuery("^SELECT id, filename, format, width, height, filepath").
					WithArgs("missing.jpg").
					WillReturnError(sql.ErrNoRows)
			},
			wantImage: models.Image{},
			wantExist: false,
			wantErr:   false,
		},
		{
			name:     "DB Error",
			filename: "error.jpg",
			mockSetup: func() {
				mock.ExpectQuery("^SELECT id, filename, format, width, height, filepath").
					WithArgs("error.jpg").
					WillReturnError(errors.New("db error"))
			},
			wantImage: models.Image{},
			wantExist: false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			gotImage, gotExist, err := repo.CheckExist(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckExist() error = %v, wantErr %v", err, tt.wantErr)
			}
			if gotExist != tt.wantExist {
				t.Errorf("CheckExist() gotExist = %v, want %v", gotExist, tt.wantExist)
			}
			if gotImage != tt.wantImage {
				t.Errorf("CheckExist() gotImage = %v, want %v", gotImage, tt.wantImage)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled mock expectations: %v", err)
			}
		})
	}
}
