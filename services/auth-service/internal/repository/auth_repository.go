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
	GetAllUsers(ctx context.Context, limit int, offset int) ([]domain.User, error)
	UpdateUser(ctx context.Context, id int, u *domain.User) error
	DeleteUser(ctx context.Context, id int) error
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

func (r *userRepository) UpdateUser(ctx context.Context, id int, u *domain.User) error {
	query := `
		UPDATE users SET 
			first_name = :first_name,
			middle_name = :middle_name,
			last_name = :last_name,
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

func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	query := `
		DELETE FROM users WHERE id = $1;
	`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
