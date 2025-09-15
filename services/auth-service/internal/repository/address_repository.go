package repository

import (
	"context"
	"fmt"

	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/domain"
	"github.com/jmoiron/sqlx"
)

type IAddressRepository interface {
	// Address operations
	CreateAddress(ctx context.Context, address *domain.UserAddress) error
	GetAddressByID(ctx context.Context, id uint) (*domain.UserAddress, error)
	GetAddressesByUserID(ctx context.Context, userID uint) ([]domain.UserAddress, error)
	GetDefaultAddressByUserID(ctx context.Context, userID uint) (*domain.UserAddress, error)
	UpdateAddress(ctx context.Context, address *domain.UserAddress) error
	DeleteAddress(ctx context.Context, id uint) error
	SetDefaultAddress(ctx context.Context, userID uint, addressID uint) error

	// Phone number operations
	CreatePhoneNumber(ctx context.Context, phone *domain.UserPhoneNumber) error
	GetPhoneNumberByID(ctx context.Context, id uint) (*domain.UserPhoneNumber, error)
	GetPhoneNumbersByUserID(ctx context.Context, userID uint) ([]domain.UserPhoneNumber, error)
	GetPrimaryPhoneByUserID(ctx context.Context, userID uint) (*domain.UserPhoneNumber, error)
	UpdatePhoneNumber(ctx context.Context, phone *domain.UserPhoneNumber) error
	DeletePhoneNumber(ctx context.Context, id uint) error
	SetPrimaryPhone(ctx context.Context, userID uint, phoneID uint) error
}

type addressRepository struct {
	db *sqlx.DB
}

func NewAddressRepository(db *sqlx.DB) IAddressRepository {
	return &addressRepository{db: db}
}

// Address operations

func (r *addressRepository) CreateAddress(ctx context.Context, address *domain.UserAddress) error {
	query := `
		INSERT INTO user_addresses (
			user_id, address_type, first_name, last_name, company,
			address_line_1, address_line_2, city, state, country,
			postal_code, phone, email, is_verified, is_default
		) VALUES (
			:user_id, :address_type, :first_name, :last_name, :company,
			:address_line_1, :address_line_2, :city, :state, :country,
			:postal_code, :phone, :email, :is_verified, :is_default
		) RETURNING id, created_at, updated_at`

	rows, err := r.db.NamedQueryContext(ctx, query, address)
	if err != nil {
		return fmt.Errorf("failed to create address: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&address.ID, &address.CreatedAt, &address.UpdatedAt); err != nil {
			return fmt.Errorf("failed to scan address result: %w", err)
		}
	}

	return rows.Err()
}

func (r *addressRepository) GetAddressByID(ctx context.Context, id uint) (*domain.UserAddress, error) {
	query := `
		SELECT * FROM user_addresses 
		WHERE id = $1 AND is_deleted = FALSE`

	var address domain.UserAddress
	if err := r.db.GetContext(ctx, &address, query, id); err != nil {
		return nil, fmt.Errorf("failed to get address: %w", err)
	}

	return &address, nil
}

func (r *addressRepository) GetAddressesByUserID(ctx context.Context, userID uint) ([]domain.UserAddress, error) {
	query := `
		SELECT * FROM user_addresses 
		WHERE user_id = $1 AND is_deleted = FALSE 
		ORDER BY is_default DESC, created_at DESC`

	addresses := make([]domain.UserAddress, 0)
	if err := r.db.SelectContext(ctx, &addresses, query, userID); err != nil {
		return nil, fmt.Errorf("failed to get addresses: %w", err)
	}

	return addresses, nil
}

func (r *addressRepository) GetDefaultAddressByUserID(ctx context.Context, userID uint) (*domain.UserAddress, error) {
	query := `
		SELECT * FROM user_addresses 
		WHERE user_id = $1 AND is_default = TRUE AND is_deleted = FALSE`

	var address domain.UserAddress
	if err := r.db.GetContext(ctx, &address, query, userID); err != nil {
		return nil, fmt.Errorf("failed to get default address: %w", err)
	}

	return &address, nil
}

func (r *addressRepository) UpdateAddress(ctx context.Context, address *domain.UserAddress) error {
	query := `
		UPDATE user_addresses SET 
			address_type = :address_type, first_name = :first_name, last_name = :last_name,
			company = :company, address_line_1 = :address_line_1, address_line_2 = :address_line_2,
			city = :city, state = :state, country = :country, postal_code = :postal_code,
			phone = :phone, email = :email, is_verified = :is_verified, is_default = :is_default,
			updated_at = NOW()
		WHERE id = :id AND is_deleted = FALSE`

	result, err := r.db.NamedExecContext(ctx, query, address)
	if err != nil {
		return fmt.Errorf("failed to update address: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("address with ID %d not found", address.ID)
	}

	return nil
}

func (r *addressRepository) DeleteAddress(ctx context.Context, id uint) error {
	query := `
		UPDATE user_addresses SET 
			is_deleted = TRUE, updated_at = NOW()
		WHERE id = $1 AND is_deleted = FALSE`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete address: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("address with ID %d not found", id)
	}

	return nil
}

func (r *addressRepository) SetDefaultAddress(ctx context.Context, userID uint, addressID uint) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// First, unset all default addresses for the user
	_, err = tx.ExecContext(ctx, `
		UPDATE user_addresses SET 
			is_default = FALSE, updated_at = NOW()
		WHERE user_id = $1 AND is_deleted = FALSE`, userID)
	if err != nil {
		return fmt.Errorf("failed to unset default addresses: %w", err)
	}

	// Then set the specified address as default
	result, err := tx.ExecContext(ctx, `
		UPDATE user_addresses SET 
			is_default = TRUE, updated_at = NOW()
		WHERE id = $1 AND user_id = $2 AND is_deleted = FALSE`, addressID, userID)
	if err != nil {
		return fmt.Errorf("failed to set default address: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("address with ID %d not found for user %d", addressID, userID)
	}

	return tx.Commit()
}

// Phone number operations

func (r *addressRepository) CreatePhoneNumber(ctx context.Context, phone *domain.UserPhoneNumber) error {
	query := `
		INSERT INTO user_phone_numbers (
			user_id, phone_type, phone_number, country_code, is_verified, is_primary
		) VALUES (
			:user_id, :phone_type, :phone_number, :country_code, :is_verified, :is_primary
		) RETURNING id, created_at, updated_at`

	rows, err := r.db.NamedQueryContext(ctx, query, phone)
	if err != nil {
		return fmt.Errorf("failed to create phone number: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&phone.ID, &phone.CreatedAt, &phone.UpdatedAt); err != nil {
			return fmt.Errorf("failed to scan phone result: %w", err)
		}
	}

	return rows.Err()
}

func (r *addressRepository) GetPhoneNumberByID(ctx context.Context, id uint) (*domain.UserPhoneNumber, error) {
	query := `
		SELECT * FROM user_phone_numbers 
		WHERE id = $1 AND is_deleted = FALSE`

	var phone domain.UserPhoneNumber
	if err := r.db.GetContext(ctx, &phone, query, id); err != nil {
		return nil, fmt.Errorf("failed to get phone number: %w", err)
	}

	return &phone, nil
}

func (r *addressRepository) GetPhoneNumbersByUserID(ctx context.Context, userID uint) ([]domain.UserPhoneNumber, error) {
	query := `
		SELECT * FROM user_phone_numbers 
		WHERE user_id = $1 AND is_deleted = FALSE 
		ORDER BY is_primary DESC, created_at DESC`

	var phones []domain.UserPhoneNumber
	if err := r.db.SelectContext(ctx, &phones, query, userID); err != nil {
		return nil, fmt.Errorf("failed to get phone numbers: %w", err)
	}

	return phones, nil
}

func (r *addressRepository) GetPrimaryPhoneByUserID(ctx context.Context, userID uint) (*domain.UserPhoneNumber, error) {
	query := `
		SELECT * FROM user_phone_numbers 
		WHERE user_id = $1 AND is_primary = TRUE AND is_deleted = FALSE`

	var phone domain.UserPhoneNumber
	if err := r.db.GetContext(ctx, &phone, query, userID); err != nil {
		return nil, fmt.Errorf("failed to get primary phone: %w", err)
	}

	return &phone, nil
}

func (r *addressRepository) UpdatePhoneNumber(ctx context.Context, phone *domain.UserPhoneNumber) error {
	query := `
		UPDATE user_phone_numbers SET 
			phone_type = :phone_type, phone_number = :phone_number, country_code = :country_code,
			is_verified = :is_verified, is_primary = :is_primary, updated_at = NOW()
		WHERE id = :id AND is_deleted = FALSE`

	result, err := r.db.NamedExecContext(ctx, query, phone)
	if err != nil {
		return fmt.Errorf("failed to update phone number: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("phone number with ID %d not found", phone.ID)
	}

	return nil
}

func (r *addressRepository) DeletePhoneNumber(ctx context.Context, id uint) error {
	query := `
		UPDATE user_phone_numbers SET 
			is_deleted = TRUE, updated_at = NOW()
		WHERE id = $1 AND is_deleted = FALSE`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete phone number: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("phone number with ID %d not found", id)
	}

	return nil
}

func (r *addressRepository) SetPrimaryPhone(ctx context.Context, userID uint, phoneID uint) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// First, unset all primary phones for the user
	_, err = tx.ExecContext(ctx, `
		UPDATE user_phone_numbers SET 
			is_primary = FALSE, updated_at = NOW()
		WHERE user_id = $1 AND is_deleted = FALSE`, userID)
	if err != nil {
		return fmt.Errorf("failed to unset primary phones: %w", err)
	}

	// Then set the specified phone as primary
	result, err := tx.ExecContext(ctx, `
		UPDATE user_phone_numbers SET 
			is_primary = TRUE, updated_at = NOW()
		WHERE id = $1 AND user_id = $2 AND is_deleted = FALSE`, phoneID, userID)
	if err != nil {
		return fmt.Errorf("failed to set primary phone: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("phone number with ID %d not found for user %d", phoneID, userID)
	}

	return tx.Commit()
}
