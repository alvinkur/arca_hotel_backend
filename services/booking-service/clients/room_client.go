package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type RoomClient struct {
	BaseURL string
	HTTP    *http.Client
}

type RoomDTO struct {
	ID           uint   `json:"id_room"`
	RoomNumber   string `json:"room_number"`
	RoomTypeID   uint   `json:"id_room_type"`
	Availability bool   `json:"availability"`
}

func NewRoomClient(baseURL string) *RoomClient {
	return &RoomClient{BaseURL: baseURL, HTTP: &http.Client{}}
}

func (c *RoomClient) GetRoom(id uint) (*RoomDTO, error) {
	resp, err := c.HTTP.Get(fmt.Sprintf("%s/api/rooms/%d", c.BaseURL, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("room tidak ditemukan")
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("gagal mengambil room (status %d)", resp.StatusCode)
	}

	var room RoomDTO
	if err := json.NewDecoder(resp.Body).Decode(&room); err != nil {
		return nil, err
	}
	return &room, nil
}

func (c *RoomClient) SetAvailability(id uint, available bool) error {
	payload := map[string]bool{"availability": available}
	data, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/api/rooms/%d", c.BaseURL, id),
		bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("gagal update availability room (status %d)", resp.StatusCode)
	}
	return nil
}
