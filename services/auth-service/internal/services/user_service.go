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
	GetAllUsersWithFilters(ctx context.Context, limit int, offset int, search string, isActive *bool, roleID *int) ([]domain.User, error)
	GetUsersCount(ctx context.Context) (int, error)
	GetUsersCountWithFilters(ctx context.Context, search string, isActive *bool, roleID *int) (int, error)
	UpdateUser(ctx context.Context, id int, u *domain.User) (*domain.User, error)
	ChangePassword(ctx context.Context, id int, currentPassword, newPassword string) error
	DeleteUser(ctx context.Context, id int) error
	// Profile methods
	GetProfile(ctx context.Context, userID int) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID int, updateData *domain.User) (*domain.User, error)
}

type userService struct {
	userRepo    repository.IUserRepository
	authService IAuthService
}

func NewUserService(userRepo repository.IUserRepository, authService IAuthService) IUserService {
	return &userService{userRepo: userRepo, authService: authService}
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

func (s *userService) GetAllUsersWithFilters(ctx context.Context, limit int, offset int, search string, isActive *bool, roleID *int) ([]domain.User, error) {
	return s.userRepo.GetAllUsersWithFilters(ctx, limit, offset, search, isActive, roleID)
}

func (s *userService) GetUsersCount(ctx context.Context) (int, error) {
	return s.userRepo.GetUsersCount(ctx)
}

func (s *userService) GetUsersCountWithFilters(ctx context.Context, search string, isActive *bool, roleID *int) (int, error) {
	return s.userRepo.GetUsersCountWithFilters(ctx, search, isActive, roleID)
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
	if updateData.MiddleName.Valid {
		updatedUser.MiddleName = updateData.MiddleName
	}
	if updateData.LastName.Valid {
		updatedUser.LastName = updateData.LastName
	}
	if updateData.Avatar.Valid {
		updatedUser.Avatar = updateData.Avatar
	}
	if updateData.Gender.Valid {
		updatedUser.Gender = updateData.Gender
	}
	// Check if DateOfBirth is valid
	if updateData.DateOfBirth.Valid {
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
	// Check if user exists
	_, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete user from database
	return s.userRepo.DeleteUser(ctx, id)
}

// GetProfile retrieves the current user's profile information
func (s *userService) GetProfile(ctx context.Context, userID int) (*domain.User, error) {
	return s.userRepo.GetUserByID(ctx, userID)
}

// UpdateProfile updates the current user's profile information
func (s *userService) UpdateProfile(ctx context.Context, userID int, updateData *domain.User) (*domain.User, error) {
	// Get the existing user to ensure it exists and merge with update data
	existingUser, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Create a copy of the existing user for updates
	updatedUser := *existingUser

	// Only update fields that are provided in the request (non-zero values)
	if updateData.FirstName != "" {
		updatedUser.FirstName = updateData.FirstName
	}
	if updateData.MiddleName.Valid {
		updatedUser.MiddleName = updateData.MiddleName
	}
	if updateData.LastName.Valid {
		updatedUser.LastName = updateData.LastName
	}
	if updateData.PhoneNumber.Valid {
		updatedUser.PhoneNumber = updateData.PhoneNumber
	}
	if updateData.Avatar.Valid {
		updatedUser.Avatar = updateData.Avatar
	}
	// Check if DateOfBirth is valid
	if updateData.DateOfBirth.Valid {
		updatedUser.DateOfBirth = updateData.DateOfBirth
	}

	// Update the user in the database
	err = s.userRepo.UpdateUser(ctx, userID, &updatedUser)
	if err != nil {
		return nil, err
	}

	// Return the updated user with all fields
	return &updatedUser, nil
}
