package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/domain"
	"github.com/jmoiron/sqlx"
)

type IPasswordResetRepository interface {
	CreatePasswordResetToken(ctx context.Context, token *domain.PasswordResetToken) error
	GetPasswordResetTokenByToken(ctx context.Context, token string) (*domain.PasswordResetToken, error)
	GetPasswordResetTokenByUserID(ctx context.Context, userID uint) (*domain.PasswordResetToken, error)
	MarkTokenAsUsed(ctx context.Context, token string) error
	DeleteExpiredTokens(ctx context.Context) error
	DeleteTokenByUserID(ctx context.Context, userID uint) error
}

type passwordResetRepository struct {
	db *sqlx.DB
}

func NewPasswordResetRepository(db *sqlx.DB) IPasswordResetRepository {
	return &passwordResetRepository{db: db}
}

// CreatePasswordResetToken stores a new password reset token in the database
func (r *passwordResetRepository) CreatePasswordResetToken(ctx context.Context, token *domain.PasswordResetToken) error {
	query := `
		INSERT INTO password_reset_tokens (
			user_id, token, expires_at, used, created_at, updated_at
		) VALUES (
			:user_id, :token, :expires_at, :used, :created_at, :updated_at
		) RETURNING id;
	`

	rows, err := r.db.NamedQueryContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("failed to create password reset token: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&token.ID); err != nil {
			return fmt.Errorf("failed to scan password reset token ID: %w", err)
		}
	}

	return rows.Err()
}

// GetPasswordResetTokenByToken retrieves a password reset token by its token string
func (r *passwordResetRepository) GetPasswordResetTokenByToken(ctx context.Context, token string) (*domain.PasswordResetToken, error) {
	query := `
		SELECT * FROM password_reset_tokens 
		WHERE token = $1 AND used = false AND expires_at > NOW();
	`

	var resetToken domain.PasswordResetToken
	if err := r.db.GetContext(ctx, &resetToken, query, token); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("password reset token not found or expired")
		}
		return nil, fmt.Errorf("failed to get password reset token: %w", err)
	}

	return &resetToken, nil
}

// GetPasswordResetTokenByUserID retrieves the most recent unused password reset token for a user
func (r *passwordResetRepository) GetPasswordResetTokenByUserID(ctx context.Context, userID uint) (*domain.PasswordResetToken, error) {
	query := `
		SELECT * FROM password_reset_tokens 
		WHERE user_id = $1 AND used = false AND expires_at > NOW()
		ORDER BY created_at DESC
		LIMIT 1;
	`

	var resetToken domain.PasswordResetToken
	if err := r.db.GetContext(ctx, &resetToken, query, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no active password reset token found for user")
		}
		return nil, fmt.Errorf("failed to get password reset token: %w", err)
	}

	return &resetToken, nil
}

// MarkTokenAsUsed marks a password reset token as used
func (r *passwordResetRepository) MarkTokenAsUsed(ctx context.Context, token string) error {
	query := `
		UPDATE password_reset_tokens 
		SET used = true, updated_at = NOW()
		WHERE token = $1;
	`

	result, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("password reset token not found")
	}

	return nil
}

// DeleteExpiredTokens removes expired tokens from the database
func (r *passwordResetRepository) DeleteExpiredTokens(ctx context.Context) error {
	query := `
		DELETE FROM password_reset_tokens 
		WHERE expires_at < NOW() OR used = true;
	`

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete expired tokens: %w", err)
	}

	return nil
}

// DeleteTokenByUserID deletes all password reset tokens for a specific user
func (r *passwordResetRepository) DeleteTokenByUserID(ctx context.Context, userID uint) error {
	query := `
		DELETE FROM password_reset_tokens 
		WHERE user_id = $1;
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete tokens for user: %w", err)
	}

	return nil
}

