package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type MLRecommendation struct {
	Name        string   `json:"name"`
	Price       float64  `json:"price"`
	Description string   `json:"description"`
	Score       float64  `json:"score"`
	RoomNumbers []string `json:"room_numbers"`
}

type MLClient struct {
	BaseURL string
	HTTP    *http.Client
}

func NewMLClient(baseURL string) *MLClient {
	return &MLClient{BaseURL: baseURL, HTTP: &http.Client{}}
}

func (c *MLClient) Recommend(message string) (*MLRecommendation, error) {
	body, _ := json.Marshal(map[string]string{"message": message})

	resp, err := c.HTTP.Post(
		c.BaseURL+"/recommend",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("ml service returned status %d", resp.StatusCode)
	}

	var results []MLRecommendation
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("empty recommendation result")
	}
	return &results[0], nil
}
