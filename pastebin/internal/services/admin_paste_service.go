package services

import (
	"context"

	"github.com/JuDyas/GolangPractice/pastebin_new/internal/repositories"
	"github.com/JuDyas/GolangPractice/pastebin_new/models"
)

type AdminPasteService interface {
	SoftDeletePaste(ctx context.Context, id string) error
	GetAllPastes(ctx context.Context, filters map[string]interface{}, sortBy string, page, limit int) ([]models.Paste, error)
}

type adminPasteServiceImpl struct {
	repo repositories.PasteRepository
}

func NewAdminPasteService(repo repositories.PasteRepository) AdminPasteService {
	return &adminPasteServiceImpl{repo: repo}
}

func (ps *adminPasteServiceImpl) SoftDeletePaste(ctx context.Context, id string) error {
	return ps.repo.SoftDelete(ctx, id)
}

func (ps *adminPasteServiceImpl) GetAllPastes(ctx context.Context, filters map[string]interface{}, sortBy string, page, limit int) ([]models.Paste, error) {
	skip := (page - 1) * limit
	return ps.repo.FindAllPastes(ctx, filters, sortBy, skip, limit)
}
