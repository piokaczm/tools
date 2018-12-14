package jira

import (
	"encoding/json"
	"errors"
	"net/http"

	perrors "github.com/pkg/errors"
)

var (
	endpoints = map[string]string{
		// "get-all-dashboards": "/rest/api/3/dashboard",
		"get-all-dashboards": "/rest/agile/1.0/board",
	}

	ErrNoAPIToken   = errors.New("jira: no API token provided")
	ErrNoBaseURL    = errors.New("jira: no base url for API provided")
	ErrNoHTTPClient = errors.New("jira: no http client provided")
	ErrUnauthorized = errors.New("jira: unauthorized request")
)

type Client struct {
	*http.Client
	apiToken string
	baseURL  string
	email    string
}

func New(apiToken, baseURL, email string, client *http.Client) (*Client, error) {
	if apiToken == "" {
		return nil, ErrNoAPIToken
	}

	if baseURL == "" {
		return nil, ErrNoBaseURL
	}

	if client == nil {
		return nil, ErrNoHTTPClient
	}

	return &Client{
		apiToken: apiToken,
		baseURL:  baseURL,
		Client:   client,
		email:    email,
	}, nil
}

func (c *Client) authorizedDo(req *http.Request) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.email, c.apiToken)

	resp, err := c.Do(req)
	if err != nil {
		return nil, perrors.Wrap(err, "jira: couldn't do a request")
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}

	return resp, nil
}

type ListResponse struct {
	StartAt    int         `json:"startAt"`
	MaxResults int         `json:"maxResults"`
	Total      int         `json:"total"`
	Prev       string      `json:"prev"`
	Next       string      `json:"next"`
	Dashboards []Dashboard `json:"dashboards"`
}

type Dashboard struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	View string `json:"view"`
}

func (c *Client) ListDashboards() (interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+endpoints["get-all-dashboards"], nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.authorizedDo(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var content interface{}
	if err := json.NewDecoder(resp.Body).Decode(&content); err != nil {
		return nil, perrors.Wrap(err, "jira: couldn't decode list dashoboards response")
	}

	return &content, nil
}
