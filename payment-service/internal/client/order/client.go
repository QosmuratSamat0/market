package order

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) GetOrderAmount(ctx context.Context, orderID string) (int64, error) {
	url := fmt.Sprintf("%s/internal/orders/%s", c.baseURL, orderID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("order-service returned status %d", resp.StatusCode)
	}

	var result struct {
		Amount int64   `json:"amount"`
		Total  float64 `json:"total"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	if result.Amount == 0 && result.Total > 0 {
		return int64(math.Round(result.Total)), nil
	}

	return result.Amount, nil
}
