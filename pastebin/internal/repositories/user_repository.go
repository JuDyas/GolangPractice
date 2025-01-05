package repositories

import (
	"context"
	"errors"
	"log"

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
	db *mongo.Collection
}

func NewUserRepository(db *mongo.Collection) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (u *userRepositoryImpl) Create(ctx context.Context, user *models.User) error {
	_, err := u.db.InsertOne(ctx, user)
	if err != nil {
		//TODO: add zap logger
		log.Println(err)
		return err
	}

	return nil
}

func (u *userRepositoryImpl) FindByID(ctx context.Context, userId string) (*models.User, error) {
	var user models.User
	err := u.db.FindOne(ctx, bson.M{"_id": userId}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	return &user, err
}

func (u *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := u.db.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	return &user, err
}
