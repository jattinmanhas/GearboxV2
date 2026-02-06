package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/dto"
)

// AuthServiceClient handles communication with the auth service
type AuthServiceClient struct {
	baseURL    string
	httpClient *http.Client
}

// UserAddress represents a user address from the auth service
type UserAddress struct {
	ID           uint   `json:"id"`
	UserID       uint   `json:"user_id"`
	AddressType  string `json:"address_type"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Company      string `json:"company"`
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
	PostalCode   string `json:"postal_code"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	IsVerified   bool   `json:"is_verified"`
	IsDefault    bool   `json:"is_default"`
}

// NewAuthServiceClient creates a new auth service client
func NewAuthServiceClient(baseURL string) *AuthServiceClient {
	return &AuthServiceClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetUserAddress fetches a user address by ID
func (c *AuthServiceClient) GetUserAddress(ctx context.Context, userID uint, addressID uint, authToken string) (*UserAddress, error) {
	url := fmt.Sprintf("%s/api/v1/auth/addresses/%d", c.baseURL, addressID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("address not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth service returned status %d", resp.StatusCode)
	}

	var response struct {
		Success bool        `json:"success"`
		Message string      `json:"message"`
		Data    UserAddress `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("auth service error: %s", response.Message)
	}

	return &response.Data, nil
}

// GetUserDefaultAddress fetches the user's default address
func (c *AuthServiceClient) GetUserDefaultAddress(ctx context.Context, userID uint, authToken string) (*UserAddress, error) {
	url := fmt.Sprintf("%s/api/v1/auth/addresses/default", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("no default address found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth service returned status %d", resp.StatusCode)
	}

	var response struct {
		Success bool        `json:"success"`
		Message string      `json:"message"`
		Data    UserAddress `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("auth service error: %s", response.Message)
	}

	return &response.Data, nil
}

// ConvertUserAddressToOrderAddress converts a user address to order address format
func (c *AuthServiceClient) ConvertUserAddressToOrderAddress(userAddr *UserAddress) map[string]interface{} {
	return map[string]interface{}{
		"first_name":  userAddr.FirstName,
		"last_name":   userAddr.LastName,
		"company":     userAddr.Company,
		"address1":    userAddr.AddressLine1,
		"address2":    userAddr.AddressLine2,
		"city":        userAddr.City,
		"state":       userAddr.State,
		"country":     userAddr.Country,
		"postal_code": userAddr.PostalCode,
		"phone":       userAddr.Phone,
		"email":       userAddr.Email,
	}
}

// GetUserProfile fetches the user's profile from the auth service
func (c *AuthServiceClient) GetUserProfile(ctx context.Context, userID uint, authToken string) (*dto.UserProfileResponse, error) {
	url := fmt.Sprintf("%s/api/v1/auth/profile", c.baseURL) // Note: GetProfile usually doesn't need ID if it uses the token, but we might have an admin endpoint or just use the token's profile

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth service returned status %d", resp.StatusCode)
	}

	var response struct {
		Success bool                    `json:"success"`
		Message string                  `json:"message"`
		Data    dto.UserProfileResponse `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("auth service error: %s", response.Message)
	}

	return &response.Data, nil
}
