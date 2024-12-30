package repository

import (
	"context"

	"github.com/JuDyas/GolangPractice/pastebin/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (string, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type userRepository struct {
	db *mongo.Collection
}

func NewUserRepository(db *mongo.Collection) UserRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) Create(ctx context.Context, user *models.User) (string, error) {
	_, err := ur.db.InsertOne(ctx, user)
	return user.ID, err
}

func (ur *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := ur.db.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	return &user, err
}
