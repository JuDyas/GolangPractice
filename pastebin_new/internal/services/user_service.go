package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/JuDyas/GolangPractice/pastebin_new/internal/repositories"
	"github.com/JuDyas/GolangPractice/pastebin_new/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, email, password string) error
	GetUser(ctx context.Context, emailOrID string) (*models.User, error)
}

type userServiceImpl struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userServiceImpl{repo: repo}
}

func (u *userServiceImpl) CreateUser(ctx context.Context, email, password string) error {
	existUser, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("error checking user existence: %v", err)
	}

	if existUser != nil {
		return fmt.Errorf("user already exists")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}

	newUser := models.User{
		ID:        primitive.NewObjectID().Hex(),
		Email:     email,
		Password:  string(hashPassword),
		CreatedAt: time.Now(),
		Role:      "standard",
	}

	err = u.repo.Create(ctx, &newUser)
	if err != nil {
		return fmt.Errorf("error creating new user: %v", err)
	}

	return nil
}

func (u *userServiceImpl) GetUser(ctx context.Context, emailOrID string) (*models.User, error) {
	if strings.Contains(emailOrID, "@") && strings.Contains(emailOrID, ".") {
		user, err := u.repo.FindByEmail(ctx, emailOrID)
		if err != nil {
			return nil, fmt.Errorf("error getting user: %v", err)
		}

		return user, nil
	} else {
		user, err := u.repo.FindByID(ctx, emailOrID)
		if err != nil {
			return nil, fmt.Errorf("error getting user: %v", err)
		}

		return user, nil
	}
}
