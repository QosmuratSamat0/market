package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	domain "github.com/QosmuratSamat0/auth-service/internal/domain/auth"
	"github.com/QosmuratSamat0/auth-service/internal/usecase/auth"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// internal API response from user-service
type userResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"`
}

type createUserRequest struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"`
}

func (c *Client) GetUserByEmail(ctx context.Context, email string) (*auth.User, error) {
	u := fmt.Sprintf("%s/internal/users/by-email?email=%s", c.baseURL, url.QueryEscape(email))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user-service returned status %d", resp.StatusCode)
	}

	var ur userResponse
	if err := json.NewDecoder(resp.Body).Decode(&ur); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return toAuthUser(&ur), nil
}

func (c *Client) GetUserByID(ctx context.Context, id string) (*auth.User, error) {
	u := fmt.Sprintf("%s/internal/users/%s", c.baseURL, url.PathEscape(id))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("user not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user-service returned status %d", resp.StatusCode)
	}

	var ur userResponse
	if err := json.NewDecoder(resp.Body).Decode(&ur); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return toAuthUser(&ur), nil
}

func (c *Client) CreateUser(ctx context.Context, user *auth.User) error {
	body, err := json.Marshal(createUserRequest{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         string(user.Role),
	})
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	u := fmt.Sprintf("%s/internal/users/", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		return fmt.Errorf("email already exists")
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("user-service returned status %d", resp.StatusCode)
	}

	// Read back the created user to get the ID
	var ur userResponse
	if err := json.NewDecoder(resp.Body).Decode(&ur); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	user.ID = ur.ID

	return nil
}

func toAuthUser(ur *userResponse) *auth.User {
	return &auth.User{
		ID:           ur.ID,
		Name:         ur.Name,
		Email:        ur.Email,
		PasswordHash: ur.PasswordHash,
		Role:         domain.Role(ur.Role),
	}
}
