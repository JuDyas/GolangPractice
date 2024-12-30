package repository

import (
	"context"
	"time"

	"github.com/JuDyas/GolangPractice/pastebin/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PasteRepository interface {
	Create(ctx context.Context, paste *models.Paste) error
	FindPaste(ctx context.Context, id string) (*models.Paste, error)
}

type pasteRepository struct {
	collection *mongo.Collection
}

func NewPasteRepository(db *mongo.Database) PasteRepository {
	return &pasteRepository{
		collection: db.Collection("pastes"),
	}
}

func (pr *pasteRepository) Create(ctx context.Context, paste *models.Paste) error {
	paste.ID = primitive.NewObjectID().Hex()
	paste.CreatedAt = time.Now()
	_, err := pr.collection.InsertOne(ctx, paste)
	return err
}

func (pr *pasteRepository) FindPaste(ctx context.Context, id string) (*models.Paste, error) {
	var paste models.Paste
	err := pr.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&paste)
	if err != nil {
		return nil, err
	}
	return &paste, nil
}
