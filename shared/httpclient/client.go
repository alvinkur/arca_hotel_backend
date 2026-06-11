package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ServiceClient struct {
	BaseURL string
	HTTP    *http.Client
}

func NewServiceClient(baseURL string) *ServiceClient {
	return &ServiceClient{
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTP:    &http.Client{},
	}
}

func (c *ServiceClient) Get(path string, result interface{}) error {
	resp, err := c.HTTP.Get(c.BaseURL + path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

func (c *ServiceClient) GetRaw(path string) (*http.Response, error) {
	resp, err := c.HTTP.Get(c.BaseURL + path)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}
	return resp, nil
}

func (c *ServiceClient) Post(path string, body, result interface{}) error {
	bodyReader, contentType := encodeBody(body)
	resp, err := c.HTTP.Post(c.BaseURL+path, contentType, bodyReader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}
	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

func (c *ServiceClient) Put(path string, body, result interface{}) error {
	bodyReader, contentType := encodeBody(body)
	req, _ := http.NewRequest("PUT", c.BaseURL+path, bodyReader)
	req.Header.Set("Content-Type", contentType)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}
	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

func (c *ServiceClient) Delete(path string) error {
	req, _ := http.NewRequest("DELETE", c.BaseURL+path, nil)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

func encodeBody(v interface{}) (io.Reader, string) {
	if v == nil {
		return nil, "application/json"
	}
	data, _ := json.Marshal(v)
	return strings.NewReader(string(data)), "application/json"
}
