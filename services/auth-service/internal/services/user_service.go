package services

import (
	"context"

	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/domain"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	RegisterNewUser(ctx context.Context, u *domain.User) error
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
	GetAllUsers(ctx context.Context, limit int, offset int) ([]domain.User, error)
	UpdateUser(ctx context.Context, id int, u *domain.User) (*domain.User, error)
	ChangePassword(ctx context.Context, id int, currentPassword, newPassword string) error
	DeleteUser(ctx context.Context, id int) error
}

type userService struct {
	userRepo repository.IUserRepository
}

func NewUserService(userRepo repository.IUserRepository) IUserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) RegisterNewUser(ctx context.Context, u *domain.User) error {
	// Hash the password
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

func (s *userService) UpdateUser(ctx context.Context, id int, updateData *domain.User) (*domain.User, error) {
	// Get the existing user to ensure it exists and merge with update data
	existingUser, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Create a copy of the existing user for updates
	updatedUser := *existingUser

	// Only update fields that are provided in the request (non-zero values)
	if updateData.FirstName != "" {
		updatedUser.FirstName = updateData.FirstName
	}
	if updateData.MiddleName != "" {
		updatedUser.MiddleName = updateData.MiddleName
	}
	if updateData.LastName != "" {
		updatedUser.LastName = updateData.LastName
	}
	if updateData.Avatar != "" {
		updatedUser.Avatar = updateData.Avatar
	}
	if updateData.Gender != "" {
		updatedUser.Gender = updateData.Gender
	}
	// Check if DateOfBirth is not zero (time.Time zero value is 0001-01-01)
	if !updateData.DateOfBirth.IsZero() {
		updatedUser.DateOfBirth = updateData.DateOfBirth
	}

	// Update the user in the database
	err = s.userRepo.UpdateUser(ctx, id, &updatedUser)
	if err != nil {
		return nil, err
	}

	// Return the updated user with all fields
	return &updatedUser, nil
}

func (s *userService) ChangePassword(ctx context.Context, id int, currentPassword, newPassword string) error {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hash)

	return s.userRepo.UpdateUser(ctx, id, user)
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	return s.userRepo.DeleteUser(ctx, id)
}
