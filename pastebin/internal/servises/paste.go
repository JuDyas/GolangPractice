package servises

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/JuDyas/GolangPractice/pastebin/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PasteService struct {
	DB *mongo.Collection
}

func NewPasteService(c *mongo.Client) *PasteService {
	return &PasteService{
		DB: c.Database("pastebin").Collection("pastes"),
	}
}

func (ps *PasteService) CreatePaste(paste *models.Paste) error {
	paste.ID = primitive.NewObjectID().Hex()
	paste.CreatedAt = time.Now()
	_, err := ps.DB.InsertOne(context.Background(), paste)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PasteService) GetPaste(id string) (*models.Paste, error) {
	var paste models.Paste
	err := ps.DB.FindOne(context.Background(), bson.M{"_id": id}).Decode(&paste)
	if err != nil {
		return nil, err
	}
	return &paste, nil
}
