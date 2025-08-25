package repositories

import (
	"context"

	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(ctx context.Context, u *models.User) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *models.User) error {
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
