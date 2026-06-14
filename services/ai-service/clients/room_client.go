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

type RoomDTO struct {
	ID           uint        `json:"id_room"`
	RoomNumber   string      `json:"room_number"`
	RoomTypeID   uint        `json:"id_room_type"`
	Availability bool        `json:"availability"`
	RoomType     RoomTypeDTO `json:"room_type,omitempty"`
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

func (c *RoomClient) GetRooms() ([]RoomDTO, error) {
	resp, err := c.HTTP.Get(c.BaseURL + "/api/rooms")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("gagal mengambil rooms (status %d)", resp.StatusCode)
	}

	var rooms []RoomDTO
	if err := json.NewDecoder(resp.Body).Decode(&rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}
