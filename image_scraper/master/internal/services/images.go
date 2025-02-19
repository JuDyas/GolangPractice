package services

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"

	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/master/internal/repositories"

	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/master/internal/models"
)

type ImageService interface {
	ProcessImage(imgUrl string) (string, error)
}

type ImageServiceImpl struct {
	repo       *repositories.ImageRepositoryImpl
	uploadsDir string
}

func NewImageService(repo *repositories.ImageRepositoryImpl) *ImageServiceImpl {
	var uploadsDir = os.Getenv("UPLOADS_DIR")
	if uploadsDir == "" {
		fmt.Println("env var UPLOADS_DIR is empty")
		uploadsDir = "/app/uploads"
	}

	return &ImageServiceImpl{
		repo:       repo,
		uploadsDir: uploadsDir,
	}
}

func (s *ImageServiceImpl) ProcessImage(imgUrl string, maxWidth int, maxHeight int) (string, error) {
	var (
		parts           = strings.Split(imgUrl, "/")
		fileName        = parts[len(parts)-1]
		filePathForSave = filepath.Join(s.uploadsDir, fileName)
		finalFilePath   = strings.TrimPrefix(filePathForSave, "/app")
	)

	existImg, exist, err := s.repo.CheckExist(fileName)
	if err != nil {
		return "", err
	}

	if !exist {
		img, err := s.downloadImage(imgUrl, fileName, filePathForSave, finalFilePath)
		if err != nil {
			fmt.Println("Error downloading image", imgUrl, err)
			return "", err
		}

		err = s.repo.SaveImage(img)
		if err != nil {
			return "", err
		}

		if ok := s.CheckImageSize(maxWidth, maxHeight, img); !ok {
			fmt.Println("Image size is too big", imgUrl)
			return "", models.MaxSizeErr
		}

		return img.Filepath, nil
	}

	ok := s.CheckImageSize(maxWidth, maxHeight, existImg)
	if !ok {
		fmt.Println("Image size is too big", imgUrl)
		return "", models.MaxSizeErr
	}

	return finalFilePath, nil
}

func (s *ImageServiceImpl) CheckImageSize(maxWidth int, maxHeight int, img models.Image) bool {
	var (
		minWidth  = int(float64(maxWidth) * models.ResizeFactor)
		minHeight = int(float64(maxHeight) * models.ResizeFactor)
	)

	if img.Width == 0 && img.Height == 0 || img.Width >= minWidth && img.Width <= maxWidth && img.Height >= minHeight && img.Height <= maxHeight {
		return true
	}
	return false
}

var insecureHttpClient = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func (s *ImageServiceImpl) downloadImage(imgUrl, fileName, filePath, finalFilePath string) (models.Image, error) {
	var images models.Image

	resp, err := insecureHttpClient.Get(imgUrl)
	if err != nil {
		return images, fmt.Errorf("could not download image: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return images, fmt.Errorf("could not download image, status code: %v", resp.StatusCode)
	}

	imgBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return images, fmt.Errorf("could not read image data: %v", err)
	}

	img, err := imaging.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return images, fmt.Errorf("could not decode image: %v", err)
	}
	size := img.Bounds()

	out, err := os.Create(filePath)
	if err != nil {
		return images, fmt.Errorf("could not create file: %v", err)
	}
	defer out.Close()

	_, err = out.Write(imgBytes)
	if err != nil {
		return images, fmt.Errorf("could not save file: %v", err)
	}

	images.Filename = fileName
	images.Format = filepath.Ext(fileName)
	images.Width = size.Dx()
	images.Height = size.Dy()
	images.Filepath = finalFilePath

	return images, nil
}
