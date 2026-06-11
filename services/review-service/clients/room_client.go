package clients

import (
	"fmt"
	"net/http"
)

type RoomClient struct {
	BaseURL string
	HTTP    *http.Client
}

func NewRoomClient(baseURL string) *RoomClient {
	return &RoomClient{BaseURL: baseURL, HTTP: &http.Client{}}
}

func (c *RoomClient) ValidateRoom(id uint) error {
	resp, err := c.HTTP.Get(fmt.Sprintf("%s/api/rooms/%d", c.BaseURL, id))
	if err != nil {
		return fmt.Errorf("gagal validasi room: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("room tidak ditemukan")
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("gagal validasi room (status %d)", resp.StatusCode)
	}
	return nil
}
