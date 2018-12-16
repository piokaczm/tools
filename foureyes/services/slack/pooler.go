package slack

import (
	"errors"
	"fmt"

	sapi "github.com/nlopes/slack"
)

type ChannelPooler struct {
	*Client
	extractor TopicExtractor
	channel   string
	lastCheck string
}

type TopicExtractor interface {
	Process([]string) ([][]string, error)
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

	return topics, nil
}

func (c *ChannelPooler) getChannelHistory() ([]string, error) {
	id, err := c.getChannelID(c.channel)
	if err != nil {
		return nil, err
	}

	params := sapi.HistoryParameters{Unreads: false, Count: 1000}
	if c.lastCheck != "" {
		params.Oldest = c.lastCheck
	}

	h, err := c.sapiClient.GetChannelHistory(id, params)
	if err != nil {
		return nil, err
	}

	var msgs []string
	for _, msg := range h.Messages {
		msgs = append(msgs, msg.Text)
	}

	if len(h.Messages) > 0 {
		c.lastCheck = h.Messages[0].Timestamp // todo: allow getting more than limit (see: has_more in slack api)
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
