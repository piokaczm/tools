package slack

import (
	"errors"
	"fmt"
	"time"

	sapi "github.com/nlopes/slack"
)

var ErrNoToken = errors.New("slack: no token provided")

type Client struct {
	sapiClient *sapi.Client
}

type ChannelPooler struct {
	*Client
	extractor TopicExtractor
	channel   string
	lastCheck time.Time
}

type TopicExtractor interface {
	Process([]string) ([][]string, error)
}

func New(token string) (*Client, error) {
	if token == "" {
		return nil, ErrNoToken
	}

	return &Client{sapi.New(token)}, nil
}

func (c *Client) NewChannelPooler(channelName string, topicExtractor TopicExtractor) (*ChannelPooler, error) {
	if topicExtractor == nil {
		return nil, errors.New("no topic extractor")
	}

	if channelName == "" {
		return nil, errors.New("no channel name")
	}

	return &ChannelPooler{
		Client:    c,
		extractor: topicExtractor,
		channel:   channelName,
		lastCheck: time.Now(), // ???
	}, nil
}

func (c *ChannelPooler) Name() string {
	return fmt.Sprintf("slack channel %q", c.channel)
}

func (c *ChannelPooler) Topics() ([][]string, error) {
	msgs, err := c.getChannelHistory()
	if err != nil {
		return nil, err
	}

	topics, err := c.extractor.Process(msgs)
	if err != nil {
		return nil, err
	}

	c.lastCheck = time.Now()
	return topics, nil
}

func (c *ChannelPooler) getChannelHistory() ([]string, error) {
	id, err := c.getChannelID(c.channel)
	if err != nil {
		return nil, err
	}

	params := sapi.HistoryParameters{Unreads: false, Count: 100} // todo: fetch data since the last check
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

func (c *ChannelPooler) getChannelID(name string) (string, error) {
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
