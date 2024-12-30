package services

import (
	"context"
	"errors"
	"time"

	"github.com/JuDyas/GolangPractice/pastebin/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"github.com/JuDyas/GolangPractice/pastebin/models"
)

type UserService interface {
	CreateUser(ctx context.Context, email string, password string) (string, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	AuthenticateUser(ctx context.Context, email string, password string) (*models.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (us *userService) CreateUser(ctx context.Context, email string, password string) (string, error) {
	_, err := us.repo.GetByEmail(ctx, email)
	if err == nil {
		return "", errors.New("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	newUser := &models.User{
		ID:        primitive.NewObjectID().Hex(),
		Email:     email,
		Password:  string(hash),
		CreatedAt: time.Now(),
		Role:      "standard",
	}

	return us.repo.Create(ctx, newUser)
}

func (us *userService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return us.repo.GetByEmail(ctx, email)
}

func (us *userService) AuthenticateUser(ctx context.Context, email string, password string) (*models.User, error) {
	user, err := us.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil

}
