package services

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/mail"

	"github.com/JuDyas/GolangPractice/pastebin_new/internal/repositories"
	"github.com/JuDyas/GolangPractice/pastebin_new/models"
	"golang.org/x/crypto/bcrypt"
)

type PasteService interface {
	CreatePaste(ctx context.Context, paste *models.Paste) error
}

type pasteServiceImpl struct {
	repo repositories.PasteRepository
}

func NewPasteService(repo repositories.PasteRepository) PasteService {
	return pasteServiceImpl{repo: repo}
}

func (ps pasteServiceImpl) CreatePaste(ctx context.Context, paste *models.Paste) error {
	if paste.TTL < 0 {
		return errors.New("paste TTL is negative")
	}

	if paste.Password != "" {
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(paste.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("could not hash password: %w", err)
		}
		paste.Password = string(hashPassword)
	}

	if paste.AllowedEmail != "" && !isValidEmail(paste.AllowedEmail) {
		return errors.New("email is invalid")
	}

	if paste.AllowedIp != "" && !isValidIP(paste.AllowedIp) {
		return errors.New("ip is invalid")
	}

	err := ps.repo.CreatePaste(ctx, paste)
	return err
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}
