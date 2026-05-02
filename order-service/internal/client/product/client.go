package product

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/QosmuratSamat/order-service/internal/lib/errs"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type Product struct {
	ID       string  `json:"id"`
	Price    float64 `json:"price"`
	SellerID string  `json:"seller_id"`
	Stock    int     `json:"stock"`
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) GetProductByID(ctx context.Context, id string) (*Product, error) {
	u := fmt.Sprintf("%s/products/%s", c.baseURL, url.PathEscape(id))

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
		return nil, errs.ProductNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("product-service returned status %d", resp.StatusCode)
	}

	var product Product
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &product, nil
}
