package slack

import (
	"errors"
	"fmt"

	sapi "github.com/nlopes/slack"
)

var ErrNoToken = errors.New("slack: no token provided")

type Client struct {
	sapiClient *sapi.Client
}

func New(token string) (*Client, error) {
	if token == "" {
		return nil, ErrNoToken
	}

	return &Client{sapi.New(token)}, nil
}

func (c *Client) GetChannelHistory(name string) ([]string, error) {
	id, err := c.getChannelID(name)
	if err != nil {
		return nil, err
	}

	params := sapi.HistoryParameters{Unreads: false, Count: 100}
	h, err := c.sapiClient.GetChannelHistory(id, params)
	if err != nil {
		return nil, err
	}

	var msgs []string
	for _, msg := range h.Messages {
		msgs = append(msgs, msg.Text)
	}

	return msgs, nil
}

func (c *Client) getChannelID(name string) (string, error) {
	chans, err := c.sapiClient.GetChannels(false)
	if err != nil {
		return "", nil
	}

	for _, c := range chans {
		if c.Name == name {
			return c.ID, nil
		}
	}

	return "", fmt.Errorf("channel %q not found within the workspace", name)
}
