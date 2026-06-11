package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RoomTypeDTO struct {
	ID          uint    `json:"id_room_type"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}

type RoomClient struct {
	BaseURL string
	HTTP    *http.Client
}

func NewRoomClient(baseURL string) *RoomClient {
	return &RoomClient{BaseURL: baseURL, HTTP: &http.Client{}}
}

func (c *RoomClient) GetRoomTypes() ([]RoomTypeDTO, error) {
	resp, err := c.HTTP.Get(c.BaseURL + "/api/room-types")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("gagal mengambil room types (status %d)", resp.StatusCode)
	}

	var types []RoomTypeDTO
	if err := json.NewDecoder(resp.Body).Decode(&types); err != nil {
		return nil, err
	}
	return types, nil
}
