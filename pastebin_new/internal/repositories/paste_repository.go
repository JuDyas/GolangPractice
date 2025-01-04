package repositories

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/JuDyas/GolangPractice/pastebin_new/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PasteRepository interface {
	CreatePaste(ctx context.Context, paste *models.Paste) error
}

type pasteRepositoryImpl struct {
	collection *mongo.Collection
}

func NewPasteRepository(collection *mongo.Collection) PasteRepository {
	return pasteRepositoryImpl{collection: collection}
}

func (r pasteRepositoryImpl) CreatePaste(ctx context.Context, paste *models.Paste) error {
	paste.ID = primitive.NewObjectID().Hex()
	paste.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, paste)
	if err != nil {
		return fmt.Errorf("save paste error: %w", err)
	}
	return nil
}
