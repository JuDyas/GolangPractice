package repositories

import (
	"context"
	"errors"
	"log"

	"github.com/JuDyas/GolangPractice/pastebin_new/dto"

	"github.com/JuDyas/GolangPractice/pastebin_new/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, userId string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
}

type userRepositoryImpl struct {
	client *mongo.Client
	dbName string
}

func NewUserRepository(client *mongo.Client, dbName string) UserRepository {
	return &userRepositoryImpl{
		client: client,
		dbName: dbName,
	}
}

func (u *userRepositoryImpl) getCollection(collectionName string) *mongo.Collection {
	return u.client.Database(u.dbName).Collection(collectionName)
}

func (u *userRepositoryImpl) Create(ctx context.Context, user *models.User) error {
	collection := u.getCollection(dto.DBUsers)
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		//TODO: add zap logger
		log.Println(err)
		return err
	}

	return nil
}

func (u *userRepositoryImpl) FindByID(ctx context.Context, userId string) (*models.User, error) {
	var user models.User
	collection := u.getCollection(dto.DBUsers)
	err := collection.FindOne(ctx, bson.M{"_id": userId}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	return &user, err
}

func (u *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	collection := u.getCollection(dto.DBUsers)
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	return &user, err
}
