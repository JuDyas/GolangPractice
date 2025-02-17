package services

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/master/internal/repositories"

	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/master/internal/models"
)

type ImageService interface {
	ProcessImage(imgUrl string) (string, error)
}

type ImageServiceImpl struct {
	repo *repositories.ImageRepositoryImpl
}

func NewImageService(repo *repositories.ImageRepositoryImpl) *ImageServiceImpl {
	return &ImageServiceImpl{
		repo: repo,
	}
}

func (s *ImageServiceImpl) ProcessImage(imgUrl string) (string, error) {
	var (
		parts    = strings.Split(imgUrl, "/")
		fileName = parts[len(parts)-1]
		filePath = filepath.Join(models.UploadsDir, fileName)
	)

	exist, err := s.repo.CheckExist(fileName)
	if err != nil {
		return "", err
	}

	if !exist {
		image, err := s.downloadImage(imgUrl, fileName, filePath)
		if err != nil {
			fmt.Println("Error downloading image", imgUrl, err)
			return "", err
		}

		s.repo.SaveImage(image)
		return image.Filepath, nil
	}

	return filePath, nil
}

var insecureHttpClient = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func (s *ImageServiceImpl) downloadImage(imgUrl, fileName, filePath string) (models.Image, error) {
	var (
		image models.Image
	)
	resp, err := insecureHttpClient.Get(imgUrl)
	if err != nil {
		return image, fmt.Errorf("could not download image: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return image, fmt.Errorf("could not download image, status code: %v", resp.StatusCode)
	}

	out, err := os.Create(filePath)
	if err != nil {
		return image, fmt.Errorf("could not create file: %v", err)
	}
	defer out.Close()

	size, err := io.Copy(out, resp.Body)
	if err != nil {
		return image, fmt.Errorf("could not save file: %v", err)

	}

	image.Filename = fileName
	image.Format = filepath.Ext(fileName)
	image.Size = int(size)
	image.Filepath = filePath
	return image, nil
}
