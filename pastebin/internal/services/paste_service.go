package services

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/JuDyas/GolangPractice/pastebin/internal/repository"

	"github.com/JuDyas/GolangPractice/pastebin/models"
)

type PasteService interface {
	CreatePaste(ctx context.Context, paste *models.Paste) error
	GetPaste(ctx context.Context, id, password string) (*models.Paste, error)
}

type pasteService struct {
	repo repository.PasteRepository
}

func NewPasteService(repo repository.PasteRepository) PasteService {
	return &pasteService{repo: repo}
}

func (ps *pasteService) CreatePaste(ctx context.Context, paste *models.Paste) error {
	if paste.TTL < 0 {
		return errors.New("invalid ttl: must be non-negative")
	}

	if paste.Password != "" {
		hashPass, err := hashPassword(paste.Password)
		if err != nil {
			return err
		}
		paste.Password = hashPass
	}

	//TODO: Handle error
	ps.repo.Create(ctx, paste)
	return nil
}

func (ps *pasteService) GetPaste(ctx context.Context, id, password string) (*models.Paste, error) {
	paste, err := ps.repo.FindPaste(ctx, id)
	if err != nil {
		return nil, err
	}

	if paste.TTL > 0 && paste.CreatedAt.Add(time.Duration(paste.TTL)*time.Second).Before(time.Now()) {
		return nil, errors.New("paste expired")
	}

	if paste.Password != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(paste.Password), []byte(password)); err != nil {
		}
	}

	return paste, nil
}

func hashPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", errors.New("password is empty")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
