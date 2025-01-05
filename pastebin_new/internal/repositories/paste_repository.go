package repositories

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/JuDyas/GolangPractice/pastebin_new/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PasteRepository interface {
	CreatePaste(ctx context.Context, paste *models.Paste) error
	GetPasteByID(ctx context.Context, id string) (*models.Paste, error)
	SoftDelete(ctx context.Context, id string) error
	FindAllPastes(ctx context.Context, filters map[string]interface{}, sortBy string, skip, limit int) ([]models.Paste, error)
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

func (r pasteRepositoryImpl) SoftDelete(ctx context.Context, id string) error {
	var (
		update = bson.M{"$set": bson.M{"deleted": true}}
		filter = bson.M{"_id": id}
	)
	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("soft delete paste error: %w", err)
	}

	return nil
}

func (r pasteRepositoryImpl) FindAllPastes(ctx context.Context, filters map[string]interface{}, sortBy string, skip, limit int) ([]models.Paste, error) {
	filter := bson.M{"deleted": false}
	for k, v := range filters {
		filter[k] = v
	}

	//sort and pagination options
	opts := options.Find().
		SetSort(bson.D{{Key: sortBy, Value: 1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("find pastes error: %w", err)
	}

	defer cursor.Close(ctx)
	var pastes []models.Paste
	if err := cursor.All(ctx, &pastes); err != nil {
		return nil, fmt.Errorf("find pastes error: %w", err)
	}

	return pastes, nil
}
