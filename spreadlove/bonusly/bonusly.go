package bonusly

import (
	"fmt"
	"net/http"
	"time"
)

const (
	baseURL             = "https://bonus.ly/api/v1"
	meEndpoint          = "/users/me"
	usersEndpoint       = "/users?user_mode=normal"
	createBonusEndpoint = "/bonuses"
)

type Client struct {
	*http.Client
	apiKey string
}

func New(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) execute(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received unexpected status %q from bonusly API", resp.Status)
	}

	return resp, nil
}
