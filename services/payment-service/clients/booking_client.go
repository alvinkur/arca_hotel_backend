package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type BookingClient struct {
	BaseURL string
	HTTP    *http.Client
}

func NewBookingClient(baseURL string) *BookingClient {
	return &BookingClient{BaseURL: baseURL, HTTP: &http.Client{}}
}

func (c *BookingClient) ValidateBooking(id uint) error {
	resp, err := c.HTTP.Get(fmt.Sprintf("%s/api/bookings/%d", c.BaseURL, id))
	if err != nil {
		return fmt.Errorf("gagal validasi booking: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("booking tidak ditemukan")
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("gagal validasi booking (status %d)", resp.StatusCode)
	}
	return nil
}

func (c *BookingClient) UpdatePaymentStatus(id uint, status string) error {
	payload := map[string]string{"status_payment": status}
	data, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/api/bookings/%d", c.BaseURL, id),
		bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("gagal update status booking (status %d)", resp.StatusCode)
	}
	return nil
}
