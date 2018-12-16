package slack

import (
	"errors"

	sapi "github.com/nlopes/slack"
)

var ErrNoToken = errors.New("slack: no token provided")

type Client struct {
	sapiClient *sapi.Client
}

type Notifier struct {
	*Client
	sendTo string
}

func New(token string) (*Client, error) {
	if token == "" {
		return nil, ErrNoToken
	}

	return &Client{sapi.New(token)}, nil
}

func (c *Client) NewNotifier(username string) (*Notifier, error) {
	userID, err := c.findUserID(username)
	if err != nil {
		return nil, err
	}

	convID, err := c.findDMChannel(userID)
	if err != nil {
		return nil, err
	}

	return &Notifier{c, convID}, nil
}

func (c *Client) findUserID(username string) (string, error) {
	usrs, err := c.sapiClient.GetUsers()
	if err != nil {
		return "", err
	}

	var userID string
	for _, u := range usrs {
		if u.Name == username {
			userID = u.ID
			break
		}
	}

	if userID == "" {
		return "", errors.New("user not found")
	}
	return userID, nil
}

func (c *Client) findDMChannel(userID string) (string, error) {
	convs, err := c.sapiClient.GetIMChannels()
	if err != nil {
		return "", err
	}

	var convID string
	for _, conv := range convs {
		if conv.User == userID {
			convID = conv.ID
		}
	}

	if convID == "" {
		return "", errors.New("conversation not found")
	}

	return convID, nil
}

func (n *Notifier) Notify(msg string) error {
	_, _, err := n.sapiClient.PostMessage(n.sendTo, sapi.MsgOptionText(msg, false))
	return err
}
