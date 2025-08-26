package services

import (
	"context"
	"fmt"

	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/domain"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	RegisterNewUser(ctx context.Context, u *domain.User) error
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
	GetAllUsers(ctx context.Context, limit int, offset int) ([]domain.User, error)
}

type userService struct {
	userRepo repository.IUserRepository
}

func NewUserService(userRepo repository.IUserRepository) IUserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) RegisterNewUser(ctx context.Context, u *domain.User) error {
	if(u.Username == "") {
		return fmt.Errorf("username is required")
	}

	if(u.Password == "" || len(u.Password) < 6) {
		return fmt.Errorf("password is required and must be at least 6 characters long")
	}

	if(u.Email == "") {
		return fmt.Errorf("email is required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hash)

	return s.userRepo.RegisterNewUser(ctx, u)
}

func (s *userService) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

func (s *userService) GetAllUsers(ctx context.Context, limit int, offset int) ([]domain.User, error) {
	return s.userRepo.GetAllUsers(ctx, limit, offset)
}