package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/JuDyas/GolangPractice/pastebin_new/dto"

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
	client *mongo.Client
	dbName string
}

func NewPasteRepository(client *mongo.Client, dbName string) PasteRepository {
	return pasteRepositoryImpl{
		client: client,
		dbName: dbName,
	}
}

func (r pasteRepositoryImpl) getCollection(collectionName string) *mongo.Collection {
	return r.client.Database(r.dbName).Collection(collectionName)
}

func (r pasteRepositoryImpl) CreatePaste(ctx context.Context, paste *models.Paste) error {
	paste.ID = primitive.NewObjectID().Hex()
	paste.CreatedAt = time.Now()
	collection := r.getCollection(dto.DBPastes)
	_, err := collection.InsertOne(ctx, paste)
	if err != nil {
		return fmt.Errorf("save paste error: %w", err)
	}
	return nil
}

func (r pasteRepositoryImpl) GetPasteByID(ctx context.Context, id string) (*models.Paste, error) {
	var paste models.Paste
	collection := r.getCollection(dto.DBPastes)
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&paste)
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
	collection := r.getCollection(dto.DBPastes)
	_, err := collection.UpdateOne(ctx, filter, update)
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

	collection := r.getCollection(dto.DBPastes)
	cursor, err := collection.Find(ctx, filter, opts)
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
