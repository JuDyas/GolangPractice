package servises

import (
	"context"
	"errors"
	"time"

	"github.com/JuDyas/GolangPractice/pastebin/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	DB *mongo.Collection
}

func NewUserService(c *mongo.Client) *UserService {
	return &UserService{
		DB: c.Database("pastebin").Collection("users"),
	}
}

func (us *UserService) CreateUser(email, password string) (string, error) {
	exist, err := us.GetByEmail(email)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return "", err
	}

	if exist != nil {
		return "", errors.New("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := models.User{
		ID:        primitive.NewObjectID().Hex(),
		Email:     email,
		Password:  string(hash),
		CreatedAt: time.Now(),
		Role:      "standard",
	}

	_, err = us.DB.InsertOne(context.Background(), user)
	if err != nil {
		return "", err
	}

	return user.ID, nil
}

func (us *UserService) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := us.DB.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *UserService) CheckPassword(email, password string) (*models.User, error) {
	user, err := us.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}
