package repository

import (
	"context"
	
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/domain"
	"github.com/jmoiron/sqlx"
)

type IUserRepository interface {
	RegisterNewUser(ctx context.Context, u *domain.User) error
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
	GetAllUsers(ctx context.Context, limit int, offset int) ([]domain.User, error)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) IUserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) RegisterNewUser(ctx context.Context, u *domain.User) error {
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

	return &user, nil
}

func (r *userRepository) GetAllUsers(ctx context.Context, limit int, offset int) ([]domain.User, error) {
	query := `
		SELECT * FROM users LIMIT $1 OFFSET $2;
	`

	var users []domain.User
	if err := r.db.SelectContext(ctx, &users, query, limit, offset); err != nil {
		return nil, err
	}

	return users, nil
}
