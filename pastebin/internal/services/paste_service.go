package services

import (
	"context"
	"errors"
	"net"
	"net/mail"
	"time"

	"github.com/JuDyas/GolangPractice/pastebin_new/dto"

	"github.com/JuDyas/GolangPractice/pastebin_new/internal/repositories"
	"github.com/JuDyas/GolangPractice/pastebin_new/models"
	"golang.org/x/crypto/bcrypt"
)

type PasteService interface {
	CreatePaste(ctx context.Context, input *dto.CreatePaste) (string, error)
	GetPaste(input *dto.GetPasteResponse, paste *dto.GetPaste) error
	GetPasteByID(ctx context.Context, id string) (*dto.GetPaste, error)
}

type pasteServiceImpl struct {
	repo repositories.PasteRepository
}

func NewPasteService(repo repositories.PasteRepository) PasteService {
	return pasteServiceImpl{repo: repo}
}

func (ps pasteServiceImpl) CreatePaste(ctx context.Context, input *dto.CreatePaste) (string, error) {
	paste := models.Paste{
		Content:      input.Content,
		TTL:          input.TTL,
		Password:     input.Password,
		AllowedEmail: input.AllowedEmail,
		AllowedIp:    input.AllowedIp,
		Authorized:   input.Authorized,
	}

	if paste.TTL < 0 {
		return "", errors.New("paste TTL is negative")
	}

	if paste.Password != "" {
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(paste.Password), bcrypt.DefaultCost)
		if err != nil {
			return "", errors.New("could not hash password")
		}
		paste.Password = string(hashPassword)
	}

	if paste.AllowedEmail != "" && !isValidEmail(paste.AllowedEmail) {
		return "", errors.New("email is invalid")
	}

	if paste.AllowedIp != "" && !isValidIP(paste.AllowedIp) {
		return "", errors.New("ip is invalid")
	}

	err := ps.repo.CreatePaste(ctx, &paste)
	if err != nil {
		return "", err
	}
	return paste.ID, nil
}

func (ps pasteServiceImpl) GetPaste(input *dto.GetPasteResponse, paste *dto.GetPaste) error {
	if paste.TTL > 0 && paste.CreatedAt.Add(time.Duration(paste.TTL)*time.Second).Before(time.Now()) {
		return errors.New("ttl has expired")
	}

	if paste.Password != "" && bcrypt.CompareHashAndPassword([]byte(paste.Password), []byte(input.Password)) != nil {
		return errors.New("password is invalid")
	}

	if paste.AllowedEmail != "" && input.Email != paste.AllowedEmail {
		return errors.New("email is invalid")
	}

	if paste.AllowedIp != "" && input.IP != paste.AllowedIp {
		return errors.New("ip is invalid")
	}

	if paste.Authorized != false && input.Email == "" {
		return errors.New("unauthorized user")
	}

	return nil
}

func (ps pasteServiceImpl) GetPasteByID(ctx context.Context, id string) (*dto.GetPaste, error) {
	repoPaste, err := ps.repo.GetPasteByID(ctx, id)
	paste := dto.GetPaste{
		ID:           repoPaste.ID,
		Content:      repoPaste.Content,
		TTL:          repoPaste.TTL,
		Password:     repoPaste.Password,
		AllowedEmail: repoPaste.AllowedEmail,
		AllowedIp:    repoPaste.AllowedIp,
		Authorized:   repoPaste.Authorized,
		Deleted:      repoPaste.Deleted,
	}

	if err != nil {
		return nil, err
	}

	return &paste, nil
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}
