package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/domain"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/dto"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/repository"
)

type IAddressService interface {
	// Address operations
	CreateAddress(ctx context.Context, userID uint, req *dto.CreateAddressRequest) (*domain.UserAddress, error)
	GetAddressByID(ctx context.Context, userID uint, addressID uint) (*domain.UserAddress, error)
	GetAddressesByUserID(ctx context.Context, userID uint) ([]domain.UserAddress, error)
	GetDefaultAddressByUserID(ctx context.Context, userID uint) (*domain.UserAddress, error)
	UpdateAddress(ctx context.Context, userID uint, addressID uint, req *dto.UpdateAddressRequest) (*domain.UserAddress, error)
	DeleteAddress(ctx context.Context, userID uint, addressID uint) error
	SetDefaultAddress(ctx context.Context, userID uint, addressID uint) error

	// Phone number operations
	CreatePhoneNumber(ctx context.Context, userID uint, req *dto.CreatePhoneNumberRequest) (*domain.UserPhoneNumber, error)
	GetPhoneNumberByID(ctx context.Context, userID uint, phoneID uint) (*domain.UserPhoneNumber, error)
	GetPhoneNumbersByUserID(ctx context.Context, userID uint) ([]domain.UserPhoneNumber, error)
	GetPrimaryPhoneByUserID(ctx context.Context, userID uint) (*domain.UserPhoneNumber, error)
	UpdatePhoneNumber(ctx context.Context, userID uint, phoneID uint, req *dto.UpdatePhoneNumberRequest) (*domain.UserPhoneNumber, error)
	DeletePhoneNumber(ctx context.Context, userID uint, phoneID uint) error
	SetPrimaryPhone(ctx context.Context, userID uint, phoneID uint) error
}

type addressService struct {
	addressRepo repository.IAddressRepository
}

func NewAddressService(addressRepo repository.IAddressRepository) IAddressService {
	return &addressService{
		addressRepo: addressRepo,
	}
}

// Address operations

func (s *addressService) CreateAddress(ctx context.Context, userID uint, req *dto.CreateAddressRequest) (*domain.UserAddress, error) {
	// Clean and validate input
	address := &domain.UserAddress{
		UserID:       userID,
		AddressType:  strings.TrimSpace(req.AddressType),
		FirstName:    strings.TrimSpace(req.FirstName),
		LastName:     strings.TrimSpace(req.LastName),
		Company:      strings.TrimSpace(req.Company),
		AddressLine1: strings.TrimSpace(req.AddressLine1),
		AddressLine2: strings.TrimSpace(req.AddressLine2),
		City:         strings.TrimSpace(req.City),
		State:        strings.TrimSpace(req.State),
		Country:      strings.TrimSpace(req.Country),
		PostalCode:   strings.TrimSpace(req.PostalCode),
		Phone:        strings.TrimSpace(req.Phone),
		Email:        strings.TrimSpace(req.Email),
		IsDefault:    req.IsDefault,
		IsVerified:   false, // New addresses start as unverified
	}

	// If this is set as default, unset other default addresses first
	if address.IsDefault {
		// Get current default address
		currentDefault, err := s.addressRepo.GetDefaultAddressByUserID(ctx, userID)
		if err == nil && currentDefault != nil {
			// Unset current default
			currentDefault.IsDefault = false
			if err := s.addressRepo.UpdateAddress(ctx, currentDefault); err != nil {
				return nil, fmt.Errorf("failed to unset current default address: %w", err)
			}
		}
	}

	if err := s.addressRepo.CreateAddress(ctx, address); err != nil {
		return nil, fmt.Errorf("failed to create address: %w", err)
	}

	return address, nil
}

func (s *addressService) GetAddressByID(ctx context.Context, userID uint, addressID uint) (*domain.UserAddress, error) {
	address, err := s.addressRepo.GetAddressByID(ctx, addressID)
	if err != nil {
		return nil, fmt.Errorf("failed to get address: %w", err)
	}

	// Verify the address belongs to the user
	if address.UserID != userID {
		return nil, fmt.Errorf("address not found")
	}

	return address, nil
}

func (s *addressService) GetAddressesByUserID(ctx context.Context, userID uint) ([]domain.UserAddress, error) {
	addresses, err := s.addressRepo.GetAddressesByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get addresses: %w", err)
	}

	return addresses, nil
}

func (s *addressService) GetDefaultAddressByUserID(ctx context.Context, userID uint) (*domain.UserAddress, error) {
	address, err := s.addressRepo.GetDefaultAddressByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get default address: %w", err)
	}

	return address, nil
}

func (s *addressService) UpdateAddress(ctx context.Context, userID uint, addressID uint, req *dto.UpdateAddressRequest) (*domain.UserAddress, error) {
	// Get existing address
	address, err := s.GetAddressByID(ctx, userID, addressID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.AddressType != "" {
		address.AddressType = strings.TrimSpace(req.AddressType)
	}
	if req.FirstName != "" {
		address.FirstName = strings.TrimSpace(req.FirstName)
	}
	if req.LastName != "" {
		address.LastName = strings.TrimSpace(req.LastName)
	}
	if req.Company != "" {
		address.Company = strings.TrimSpace(req.Company)
	}
	if req.AddressLine1 != "" {
		address.AddressLine1 = strings.TrimSpace(req.AddressLine1)
	}
	if req.AddressLine2 != "" {
		address.AddressLine2 = strings.TrimSpace(req.AddressLine2)
	}
	if req.City != "" {
		address.City = strings.TrimSpace(req.City)
	}
	if req.State != "" {
		address.State = strings.TrimSpace(req.State)
	}
	if req.Country != "" {
		address.Country = strings.TrimSpace(req.Country)
	}
	if req.PostalCode != "" {
		address.PostalCode = strings.TrimSpace(req.PostalCode)
	}
	if req.Phone != "" {
		address.Phone = strings.TrimSpace(req.Phone)
	}
	if req.Email != "" {
		address.Email = strings.TrimSpace(req.Email)
	}
	if req.IsDefault != nil {
		address.IsDefault = *req.IsDefault
	}

	// If setting as default, unset other default addresses
	if address.IsDefault {
		currentDefault, err := s.addressRepo.GetDefaultAddressByUserID(ctx, userID)
		if err == nil && currentDefault != nil && currentDefault.ID != addressID {
			currentDefault.IsDefault = false
			if err := s.addressRepo.UpdateAddress(ctx, currentDefault); err != nil {
				return nil, fmt.Errorf("failed to unset current default address: %w", err)
			}
		}
	}

	if err := s.addressRepo.UpdateAddress(ctx, address); err != nil {
		return nil, fmt.Errorf("failed to update address: %w", err)
	}

	return address, nil
}

func (s *addressService) DeleteAddress(ctx context.Context, userID uint, addressID uint) error {
	// Verify the address belongs to the user
	_, err := s.GetAddressByID(ctx, userID, addressID)
	if err != nil {
		return err
	}

	if err := s.addressRepo.DeleteAddress(ctx, addressID); err != nil {
		return fmt.Errorf("failed to delete address: %w", err)
	}

	return nil
}

func (s *addressService) SetDefaultAddress(ctx context.Context, userID uint, addressID uint) error {
	// Verify the address belongs to the user
	_, err := s.GetAddressByID(ctx, userID, addressID)
	if err != nil {
		return err
	}

	if err := s.addressRepo.SetDefaultAddress(ctx, userID, addressID); err != nil {
		return fmt.Errorf("failed to set default address: %w", err)
	}

	return nil
}

// Phone number operations

func (s *addressService) CreatePhoneNumber(ctx context.Context, userID uint, req *dto.CreatePhoneNumberRequest) (*domain.UserPhoneNumber, error) {
	// Clean and validate input
	phone := &domain.UserPhoneNumber{
		UserID:      userID,
		PhoneType:   strings.TrimSpace(req.PhoneType),
		PhoneNumber: strings.TrimSpace(req.PhoneNumber),
		CountryCode: strings.TrimSpace(req.CountryCode),
		IsPrimary:   req.IsPrimary,
		IsVerified:  false, // New phone numbers start as unverified
	}

	// If this is set as primary, unset other primary phones first
	if phone.IsPrimary {
		currentPrimary, err := s.addressRepo.GetPrimaryPhoneByUserID(ctx, userID)
		if err == nil && currentPrimary != nil {
			// Unset current primary
			currentPrimary.IsPrimary = false
			if err := s.addressRepo.UpdatePhoneNumber(ctx, currentPrimary); err != nil {
				return nil, fmt.Errorf("failed to unset current primary phone: %w", err)
			}
		}
	}

	if err := s.addressRepo.CreatePhoneNumber(ctx, phone); err != nil {
		return nil, fmt.Errorf("failed to create phone number: %w", err)
	}

	return phone, nil
}

func (s *addressService) GetPhoneNumberByID(ctx context.Context, userID uint, phoneID uint) (*domain.UserPhoneNumber, error) {
	phone, err := s.addressRepo.GetPhoneNumberByID(ctx, phoneID)
	if err != nil {
		return nil, fmt.Errorf("failed to get phone number: %w", err)
	}

	// Verify the phone belongs to the user
	if phone.UserID != userID {
		return nil, fmt.Errorf("phone number not found")
	}

	return phone, nil
}

func (s *addressService) GetPhoneNumbersByUserID(ctx context.Context, userID uint) ([]domain.UserPhoneNumber, error) {
	phones, err := s.addressRepo.GetPhoneNumbersByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get phone numbers: %w", err)
	}

	return phones, nil
}

func (s *addressService) GetPrimaryPhoneByUserID(ctx context.Context, userID uint) (*domain.UserPhoneNumber, error) {
	phone, err := s.addressRepo.GetPrimaryPhoneByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get primary phone: %w", err)
	}

	return phone, nil
}

func (s *addressService) UpdatePhoneNumber(ctx context.Context, userID uint, phoneID uint, req *dto.UpdatePhoneNumberRequest) (*domain.UserPhoneNumber, error) {
	// Get existing phone
	phone, err := s.GetPhoneNumberByID(ctx, userID, phoneID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.PhoneType != "" {
		phone.PhoneType = strings.TrimSpace(req.PhoneType)
	}
	if req.PhoneNumber != "" {
		phone.PhoneNumber = strings.TrimSpace(req.PhoneNumber)
	}
	if req.CountryCode != "" {
		phone.CountryCode = strings.TrimSpace(req.CountryCode)
	}
	if req.IsPrimary != nil {
		phone.IsPrimary = *req.IsPrimary
	}

	// If setting as primary, unset other primary phones
	if phone.IsPrimary {
		currentPrimary, err := s.addressRepo.GetPrimaryPhoneByUserID(ctx, userID)
		if err == nil && currentPrimary != nil && currentPrimary.ID != phoneID {
			currentPrimary.IsPrimary = false
			if err := s.addressRepo.UpdatePhoneNumber(ctx, currentPrimary); err != nil {
				return nil, fmt.Errorf("failed to unset current primary phone: %w", err)
			}
		}
	}

	if err := s.addressRepo.UpdatePhoneNumber(ctx, phone); err != nil {
		return nil, fmt.Errorf("failed to update phone number: %w", err)
	}

	return phone, nil
}

func (s *addressService) DeletePhoneNumber(ctx context.Context, userID uint, phoneID uint) error {
	// Verify the phone belongs to the user
	_, err := s.GetPhoneNumberByID(ctx, userID, phoneID)
	if err != nil {
		return err
	}

	if err := s.addressRepo.DeletePhoneNumber(ctx, phoneID); err != nil {
		return fmt.Errorf("failed to delete phone number: %w", err)
	}

	return nil
}

func (s *addressService) SetPrimaryPhone(ctx context.Context, userID uint, phoneID uint) error {
	// Verify the phone belongs to the user
	_, err := s.GetPhoneNumberByID(ctx, userID, phoneID)
	if err != nil {
		return err
	}

	if err := s.addressRepo.SetPrimaryPhone(ctx, userID, phoneID); err != nil {
		return fmt.Errorf("failed to set primary phone: %w", err)
	}

	return nil
}
