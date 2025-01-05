package repositories

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/JuDyas/GolangPractice/pastebin_new/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PasteRepository interface {
	CreatePaste(ctx context.Context, paste *models.Paste) error
	GetPasteByID(ctx context.Context, id string) (*models.Paste, error)
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

func (r pasteRepositoryImpl) GetPasteByID(ctx context.Context, id string) (*models.Paste, error) {
	var paste models.Paste
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&paste)
	if err != nil {
		return nil, fmt.Errorf("get paste error: %w", err)
	}

	return &paste, nil
}
