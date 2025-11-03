package repository

import (
	"context"
	"fmt"

	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/domain"
	"github.com/jmoiron/sqlx"
)

type IUserRepository interface {
	RegisterNewUser(ctx context.Context, u *domain.User) error
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetAllUsers(ctx context.Context, limit int, offset int) ([]domain.User, error)
	GetAllUsersWithFilters(ctx context.Context, limit int, offset int, search string, isActive *bool, roleID *int) ([]domain.User, error)
	GetUsersCount(ctx context.Context) (int, error)
	GetUsersCountWithFilters(ctx context.Context, search string, isActive *bool, roleID *int) (int, error)
	UpdateUser(ctx context.Context, id int, u *domain.User) error
	DeleteUser(ctx context.Context, id int) error
	// User Analytics
	GetUserAnalytics(ctx context.Context) (*domain.UserAnalytics, error)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) IUserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) RegisterNewUser(ctx context.Context, u *domain.User) error {
	// Don't set role_id here - it will be set by the migration or default constraint
	// u.RoleID = domain.RoleIDUser // Remove this line

	query := `
		INSERT INTO users (
			username, password, email, first_name, middle_name, last_name, avatar, gender, date_of_birth
		) VALUES (
			:username, :password, :email, :first_name, :middle_name, :last_name, :avatar, :gender, :date_of_birth
		) RETURNING id;
	`

	// use NamedQueryRowx to bind struct fields by db tags
	rows, err := r.db.NamedQueryContext(ctx, query, u)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		if err := rows.Scan(&u.ID); err != nil {
			return err
		}
	}
	return rows.Err()
}

func (r *userRepository) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	query := `
		SELECT * FROM users WHERE id = $1;
	`

	var user domain.User
	if err := r.db.GetContext(ctx, &user, query, id); err != nil {
		return nil, err
	}

	roleName, ok := domain.RoleNames[int(user.RoleID)]
	if !ok {
		user.RoleID, user.Role = domain.GetDefaultRole()
	} else {
		user.Role = roleName
	}

	return &user, nil
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT * FROM users WHERE username = $1;
	`

	var user domain.User
	if err := r.db.GetContext(ctx, &user, query, username); err != nil {
		return nil, err
	}

	roleName, ok := domain.RoleNames[int(user.RoleID)]
	if !ok {
		user.RoleID, user.Role = domain.GetDefaultRole()
	} else {
		user.Role = roleName
	}

	return &user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT * FROM users WHERE email = $1 AND is_deleted = false;
	`

	var user domain.User
	if err := r.db.GetContext(ctx, &user, query, email); err != nil {
		return nil, err
	}

	roleName, ok := domain.RoleNames[int(user.RoleID)]
	if !ok {
		user.RoleID, user.Role = domain.GetDefaultRole()
	} else {
		user.Role = roleName
	}

	return &user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, id int, u *domain.User) error {
	query := `
		UPDATE users SET 
			first_name = :first_name,
			middle_name = :middle_name,
			last_name = :last_name,
			phone_number = :phone_number,
			avatar = :avatar,
			gender = :gender,
			date_of_birth = :date_of_birth,
			password = :password,
			updated_at = NOW()
		WHERE id = :id;
	`

	// Add the ID to the struct for the WHERE clause
	u.ID = uint(id)

	// Execute the update query
	result, err := r.db.NamedExecContext(ctx, query, u)
	if err != nil {
		return err
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found with id %d", id)
	}

	return nil
}

func (r *userRepository) GetAllUsers(ctx context.Context, limit int, offset int) ([]domain.User, error) {
	query := `
		SELECT * FROM users LIMIT $1 OFFSET $2;
	`

	var users []domain.User
	if err := r.db.SelectContext(ctx, &users, query, limit, offset); err != nil {
		return nil, err
	}

	for i := range users {
		roleName, ok := domain.RoleNames[int(users[i].RoleID)]
		if !ok {
			users[i].RoleID, users[i].Role = domain.GetDefaultRole()
		} else {
			users[i].Role = roleName
		}
	}

	return users, nil
}

func (r *userRepository) GetAllUsersWithFilters(ctx context.Context, limit int, offset int, search string, isActive *bool, roleID *int) ([]domain.User, error) {
	query := `
		SELECT * FROM users 
		WHERE is_deleted = false
	`
	args := []interface{}{}
	argIndex := 1

	// Add search filter
	if search != "" {
		query += ` AND (username ILIKE $` + fmt.Sprintf("%d", argIndex) + ` OR email ILIKE $` + fmt.Sprintf("%d", argIndex+1) + ` OR first_name ILIKE $` + fmt.Sprintf("%d", argIndex+2) + ` OR last_name ILIKE $` + fmt.Sprintf("%d", argIndex+3) + `)`
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern, searchPattern)
		argIndex += 4
	}

	// Add is_active filter
	if isActive != nil {
		query += ` AND is_active = $` + fmt.Sprintf("%d", argIndex)
		args = append(args, *isActive)
		argIndex++
	}

	// Add role_id filter
	if roleID != nil {
		query += ` AND role_id = $` + fmt.Sprintf("%d", argIndex)
		args = append(args, *roleID)
		argIndex++
	}

	// Add ordering and pagination
	query += ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", argIndex) + ` OFFSET $` + fmt.Sprintf("%d", argIndex+1)
	args = append(args, limit, offset)

	var users []domain.User
	if err := r.db.SelectContext(ctx, &users, query, args...); err != nil {
		return nil, err
	}

	// Set role names for each user
	for i := range users {
		roleName, ok := domain.RoleNames[int(users[i].RoleID)]
		if !ok {
			users[i].RoleID, users[i].Role = domain.GetDefaultRole()
		} else {
			users[i].Role = roleName
		}
	}

	return users, nil
}

func (r *userRepository) GetUsersCount(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*) FROM users WHERE is_deleted = false;
	`

	var count int
	if err := r.db.GetContext(ctx, &count, query); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *userRepository) GetUsersCountWithFilters(ctx context.Context, search string, isActive *bool, roleID *int) (int, error) {
	query := `
		SELECT COUNT(*) FROM users 
		WHERE is_deleted = false
	`
	args := []interface{}{}
	argIndex := 1

	// Add search filter
	if search != "" {
		query += ` AND (username ILIKE $` + fmt.Sprintf("%d", argIndex) + ` OR email ILIKE $` + fmt.Sprintf("%d", argIndex+1) + ` OR first_name ILIKE $` + fmt.Sprintf("%d", argIndex+2) + ` OR last_name ILIKE $` + fmt.Sprintf("%d", argIndex+3) + `)`
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern, searchPattern)
		argIndex += 4
	}

	// Add is_active filter
	if isActive != nil {
		query += ` AND is_active = $` + fmt.Sprintf("%d", argIndex)
		args = append(args, *isActive)
		argIndex++
	}

	// Add role_id filter
	if roleID != nil {
		query += ` AND role_id = $` + fmt.Sprintf("%d", argIndex)
		args = append(args, *roleID)
		argIndex++
	}

	var count int
	if err := r.db.GetContext(ctx, &count, query, args...); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	query := `
		DELETE FROM users WHERE id = $1;
	`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// GetUserAnalytics retrieves user analytics
func (r *userRepository) GetUserAnalytics(ctx context.Context) (*domain.UserAnalytics, error) {
	analytics := &domain.UserAnalytics{}

	// Get total users count
	err := r.db.GetContext(ctx, &analytics.TotalUsers, `SELECT COUNT(*) FROM users WHERE is_deleted = false`)
	if err != nil {
		return nil, fmt.Errorf("failed to get total users count: %w", err)
	}

	// Get active users count
	err = r.db.GetContext(ctx, &analytics.ActiveUsers, `SELECT COUNT(*) FROM users WHERE is_active = true AND is_deleted = false`)
	if err != nil {
		return nil, fmt.Errorf("failed to get active users count: %w", err)
	}

	// Get new users today
	err = r.db.GetContext(ctx, &analytics.NewUsersToday, `
		SELECT COUNT(*) FROM users 
		WHERE DATE(created_at) = CURRENT_DATE AND is_deleted = false`)
	if err != nil {
		return nil, fmt.Errorf("failed to get new users today: %w", err)
	}

	// Get new users this week
	err = r.db.GetContext(ctx, &analytics.NewUsersThisWeek, `
		SELECT COUNT(*) FROM users 
		WHERE created_at >= DATE_TRUNC('week', CURRENT_DATE) AND is_deleted = false`)
	if err != nil {
		return nil, fmt.Errorf("failed to get new users this week: %w", err)
	}

	// Get new users this month
	err = r.db.GetContext(ctx, &analytics.NewUsersThisMonth, `
		SELECT COUNT(*) FROM users 
		WHERE created_at >= DATE_TRUNC('month', CURRENT_DATE) AND is_deleted = false`)
	if err != nil {
		return nil, fmt.Errorf("failed to get new users this month: %w", err)
	}

	// Get users by role
	var usersByRole []domain.UserRoleCount
	err = r.db.SelectContext(ctx, &usersByRole, `
		SELECT r.name as role, COUNT(u.id) as count
		FROM roles r
		LEFT JOIN users u ON r.id = u.role_id AND u.is_deleted = false
		GROUP BY r.id, r.name
		ORDER BY count DESC`)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by role: %w", err)
	}
	analytics.UsersByRole = usersByRole

	// Get user registration trend (last 30 days)
	var registrationTrend []domain.UserRegistrationData
	err = r.db.SelectContext(ctx, &registrationTrend, `
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as count
		FROM users 
		WHERE created_at >= CURRENT_DATE - INTERVAL '30 days' 
		AND is_deleted = false
		GROUP BY DATE(created_at)
		ORDER BY date ASC`)
	if err != nil {
		return nil, fmt.Errorf("failed to get user registration trend: %w", err)
	}
	analytics.UserRegistrationTrend = registrationTrend

	return analytics, nil
}
