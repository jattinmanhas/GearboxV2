package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/domain"
	"github.com/jmoiron/sqlx"
)

type IRefreshTokenRepository interface {
	CreateRefreshToken(ctx context.Context, token *domain.RefreshToken) error
	GetRefreshTokenByToken(ctx context.Context, refreshToken string) (*domain.RefreshToken, error)
	GetRefreshTokensByUserID(ctx context.Context, userID uint) ([]domain.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, refreshToken string) error
	RevokeAllUserTokens(ctx context.Context, userID uint) error
	CleanupExpiredTokens(ctx context.Context) error
	DeleteRefreshToken(ctx context.Context, refreshToken string) error
}

type refreshTokenRepository struct {
	db *sqlx.DB
}

func NewRefreshTokenRepository(db *sqlx.DB) IRefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

// CreateRefreshToken stores a new refresh token in the database
func (r *refreshTokenRepository) CreateRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (
			user_id, refresh_token, user_agent, ip_address, expires_at, created_at, is_revoked
		) VALUES (
			:user_id, :refresh_token, :user_agent, :ip_address, :expires_at, :created_at, :is_revoked
		) RETURNING id;
	`

	rows, err := r.db.NamedQueryContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("failed to create refresh token: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&token.ID); err != nil {
			return fmt.Errorf("failed to scan refresh token ID: %w", err)
		}
	}

	return rows.Err()
}

// GetRefreshTokenByToken retrieves a refresh token by its token string
func (r *refreshTokenRepository) GetRefreshTokenByToken(ctx context.Context, refreshToken string) (*domain.RefreshToken, error) {
	query := `
		SELECT * FROM refresh_tokens 
		WHERE refresh_token = $1 AND is_revoked = false AND expires_at > NOW();
	`

	var token domain.RefreshToken
	if err := r.db.GetContext(ctx, &token, query, refreshToken); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("refresh token not found or expired")
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	return &token, nil
}

// GetRefreshTokensByUserID retrieves all refresh tokens for a specific user
func (r *refreshTokenRepository) GetRefreshTokensByUserID(ctx context.Context, userID uint) ([]domain.RefreshToken, error) {
	query := `
		SELECT * FROM refresh_tokens 
		WHERE user_id = $1 AND is_revoked = false AND expires_at > NOW()
		ORDER BY created_at DESC;
	`

	var tokens []domain.RefreshToken
	if err := r.db.SelectContext(ctx, &tokens, query, userID); err != nil {
		return nil, fmt.Errorf("failed to get refresh tokens: %w", err)
	}

	return tokens, nil
}

// RevokeRefreshToken marks a specific refresh token as revoked
func (r *refreshTokenRepository) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	query := `
		UPDATE refresh_tokens 
		SET is_revoked = true, last_used_at = NOW()
		WHERE refresh_token = $1;
	`

	result, err := r.db.ExecContext(ctx, query, refreshToken)
	if err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("refresh token not found")
	}

	return nil
}

// RevokeAllUserTokens revokes all refresh tokens for a specific user
func (r *refreshTokenRepository) RevokeAllUserTokens(ctx context.Context, userID uint) error {
	query := `
		UPDATE refresh_tokens 
		SET is_revoked = true
		WHERE user_id = $1 AND is_revoked = false;
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke user tokens: %w", err)
	}

	return nil
}

// CleanupExpiredTokens removes expired tokens from the database
func (r *refreshTokenRepository) CleanupExpiredTokens(ctx context.Context) error {
	query := `
		DELETE FROM refresh_tokens 
		WHERE expires_at < NOW() OR is_revoked = true;
	`

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired tokens: %w", err)
	}

	return nil
}

// DeleteRefreshToken permanently removes a refresh token
func (r *refreshTokenRepository) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	query := `
		DELETE FROM refresh_tokens 
		WHERE refresh_token = $1;
	`

	result, err := r.db.ExecContext(ctx, query, refreshToken)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("refresh token not found")
	}

	return nil
}
