package clients

import (
	"fmt"
	"net/http"
)

type AuthClient struct {
	BaseURL string
	HTTP    *http.Client
}

func NewAuthClient(baseURL string) *AuthClient {
	return &AuthClient{BaseURL: baseURL, HTTP: &http.Client{}}
}

func (c *AuthClient) ValidateCustomer(id uint) error {
	resp, err := c.HTTP.Get(fmt.Sprintf("%s/api/customers/%d", c.BaseURL, id))
	if err != nil {
		return fmt.Errorf("gagal validasi customer: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("customer tidak ditemukan")
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("gagal validasi customer (status %d)", resp.StatusCode)
	}
	return nil
}
