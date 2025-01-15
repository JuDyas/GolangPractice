package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/JuDyas/GolangPractice/pastebin_new/dto"

	"github.com/JuDyas/GolangPractice/pastebin_new/internal/repositories"
	"github.com/JuDyas/GolangPractice/pastebin_new/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, email, password string) error
	GetUser(ctx context.Context, emailOrID string) (*dto.User, error)
	Authenticate(ctx context.Context, email, password string) (*dto.User, error)
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

	var newUser = models.User{
		ID:        primitive.NewObjectID().Hex(),
		Email:     email,
		Password:  string(hashPassword),
		CreatedAt: time.Now(),
		Role:      models.RoleStandard,
	}

	err = u.repo.Create(ctx, &newUser)
	if err != nil {
		return fmt.Errorf("error creating new user: %v", err)
	}

	return nil
}

// TODO: ПЕРЕПИСАТЬ!!!!
func (u *userServiceImpl) GetUser(ctx context.Context, emailOrID string) (*dto.User, error) {
	if strings.Contains(emailOrID, "@") && strings.Contains(emailOrID, ".") {
		userRepo, err := u.repo.FindByEmail(ctx, emailOrID)
		if err != nil {
			return nil, fmt.Errorf("error getting user: %v", err)
		}

		var user = dto.User{
			ID:        userRepo.ID,
			Email:     userRepo.Email,
			Password:  userRepo.Password,
			CreatedAt: userRepo.CreatedAt,
			Role:      userRepo.Role,
		}

		return &user, nil
	} else {
		userRepo, err := u.repo.FindByID(ctx, emailOrID)
		if err != nil {
			return nil, fmt.Errorf("error getting user: %v", err)
		}

		var user = dto.User{
			ID:        userRepo.ID,
			Email:     userRepo.Email,
			Password:  userRepo.Password,
			CreatedAt: userRepo.CreatedAt,
			Role:      userRepo.Role,
		}

		return &user, nil
	}
}

func (u *userServiceImpl) Authenticate(ctx context.Context, email, password string) (*dto.User, error) {
	userRepo, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %v", err)
	}

	if userRepo == nil {
		return nil, dto.NotFoundErr
	}

	err = bcrypt.CompareHashAndPassword([]byte(userRepo.Password), []byte(password))
	if err != nil {
		return nil, dto.InvalidPasswordErr
	}

	var user = dto.User{
		ID:        userRepo.ID,
		Email:     userRepo.Email,
		Password:  userRepo.Password,
		CreatedAt: userRepo.CreatedAt,
		Role:      userRepo.Role,
	}

	return &user, nil
}
