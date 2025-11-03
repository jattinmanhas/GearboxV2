package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/domain"
	"github.com/jmoiron/sqlx"
)

type IOAuthRepository interface {
	CreateOAuthProvider(ctx context.Context, provider *domain.OAuthProviderLink) error
	GetOAuthProviderByProviderAndUserID(ctx context.Context, provider string, providerUserID string) (*domain.OAuthProviderLink, error)
	GetOAuthProvidersByUserID(ctx context.Context, userID uint) ([]*domain.OAuthProviderLink, error)
	UpdateOAuthProvider(ctx context.Context, provider *domain.OAuthProviderLink) error
	DeleteOAuthProvider(ctx context.Context, id uint) error
	UnlinkProvider(ctx context.Context, userID uint, provider string) error
	GetOAuthProviderByUserAndProvider(ctx context.Context, userID uint, provider string) (*domain.OAuthProviderLink, error)
}

type oauthRepository struct {
	db *sqlx.DB
}

func NewOAuthRepository(db *sqlx.DB) IOAuthRepository {
	return &oauthRepository{db: db}
}

// CreateOAuthProvider creates a new OAuth provider link for a user
func (r *oauthRepository) CreateOAuthProvider(ctx context.Context, provider *domain.OAuthProviderLink) error {
	query := `
		INSERT INTO oauth_providers (
			user_id, provider, provider_user_id, email, access_token, refresh_token, expires_at
		) VALUES (
			:user_id, :provider, :provider_user_id, :email, :access_token, :refresh_token, :expires_at
		) RETURNING id, created_at, updated_at;
	`

	rows, err := r.db.NamedQueryContext(ctx, query, provider)
	if err != nil {
		return fmt.Errorf("failed to create oauth provider: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&provider.ID, &provider.CreatedAt, &provider.UpdatedAt); err != nil {
			return fmt.Errorf("failed to scan oauth provider result: %w", err)
		}
	}

	return rows.Err()
}

// GetOAuthProviderByProviderAndUserID retrieves an OAuth provider by provider type and provider user ID
func (r *oauthRepository) GetOAuthProviderByProviderAndUserID(ctx context.Context, provider string, providerUserID string) (*domain.OAuthProviderLink, error) {
	query := `
		SELECT * FROM oauth_providers 
		WHERE provider = $1 AND provider_user_id = $2
		LIMIT 1;
	`

	var oauthProvider domain.OAuthProviderLink
	if err := r.db.GetContext(ctx, &oauthProvider, query, provider, providerUserID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get oauth provider: %w", err)
	}

	return &oauthProvider, nil
}

// GetOAuthProviderByUserAndProvider retrieves an OAuth provider by user ID and provider type
func (r *oauthRepository) GetOAuthProviderByUserAndProvider(ctx context.Context, userID uint, provider string) (*domain.OAuthProviderLink, error) {
	query := `
		SELECT * FROM oauth_providers 
		WHERE user_id = $1 AND provider = $2
		LIMIT 1;
	`

	var oauthProvider domain.OAuthProviderLink
	if err := r.db.GetContext(ctx, &oauthProvider, query, userID, provider); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get oauth provider: %w", err)
	}

	return &oauthProvider, nil
}

// GetOAuthProvidersByUserID retrieves all OAuth providers for a user
func (r *oauthRepository) GetOAuthProvidersByUserID(ctx context.Context, userID uint) ([]*domain.OAuthProviderLink, error) {
	query := `
		SELECT id, user_id, provider, provider_user_id, email, expires_at, created_at, updated_at
		FROM oauth_providers 
		WHERE user_id = $1
		ORDER BY created_at DESC;
	`

	var providers []*domain.OAuthProviderLink
	if err := r.db.SelectContext(ctx, &providers, query, userID); err != nil {
		return nil, fmt.Errorf("failed to get oauth providers: %w", err)
	}

	return providers, nil
}

// UpdateOAuthProvider updates an existing OAuth provider link
func (r *oauthRepository) UpdateOAuthProvider(ctx context.Context, provider *domain.OAuthProviderLink) error {
	query := `
		UPDATE oauth_providers SET 
			access_token = :access_token,
			refresh_token = :refresh_token,
			expires_at = :expires_at,
			updated_at = NOW()
		WHERE id = :id;
	`

	result, err := r.db.NamedExecContext(ctx, query, provider)
	if err != nil {
		return fmt.Errorf("failed to update oauth provider: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("oauth provider not found")
	}

	return nil
}

// DeleteOAuthProvider deletes an OAuth provider link by ID
func (r *oauthRepository) DeleteOAuthProvider(ctx context.Context, id uint) error {
	query := `
		DELETE FROM oauth_providers WHERE id = $1;
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete oauth provider: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("oauth provider not found")
	}

	return nil
}

// UnlinkProvider unlinks an OAuth provider from a user
func (r *oauthRepository) UnlinkProvider(ctx context.Context, userID uint, provider string) error {
	query := `
		DELETE FROM oauth_providers WHERE user_id = $1 AND provider = $2;
	`

	result, err := r.db.ExecContext(ctx, query, userID, provider)
	if err != nil {
		return fmt.Errorf("failed to unlink oauth provider: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("oauth provider link not found")
	}

	return nil
}
