package viacep

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Alef-Daniel/challenge-deploy-cloud-run/internal/entity"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient() *Client {
	return &Client{
		BaseURL:    "https://viacep.com.br",
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type response struct {
	Localidade string `json:"localidade"`
}

func (c *Client) FindCity(ctx context.Context, cep string) (string, error) {
	url := fmt.Sprintf("%s/ws/%s/json/", c.BaseURL, cep)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("viacep request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusNotFound {
		return "", entity.ErrZipcodeNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("viacep: unexpected status %d", resp.StatusCode)
	}

	var body response
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("viacep decode: %w", err)
	}
	if body.Localidade == "" {
		return "", entity.ErrZipcodeNotFound
	}

	return body.Localidade, nil
}
