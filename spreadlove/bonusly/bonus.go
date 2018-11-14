package bonusly

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	errInsufficientCoins = errors.New("Splitting remaining coins between provided folks leave them with 0 coins per person :(")
	errBonusAborted      = errors.New("Bonus was aborted!")
)

type BonusPayload struct {
	Reason string `json:"reason"`
}

func (c *Client) bonus(names []string, tcToSpread int, msg string) (string, error) {
	data, err := buildData(names, tcToSpread, msg)
	if err != nil {
		return "", err
	}

	fmt.Printf(
		"You are about to send following bonus:\n'%s'\nAre you sure you want to proceed? (y/n)",
		data,
	)
	var input string
	fmt.Scanln(&input)
	if input == "n" {
		return "", errBonusAborted
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+createBonusEndpoint, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.execute(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return string(data), nil
}

func buildData(names []string, coins int, msg string) ([]byte, error) {
	var message string
	if msg == "" {
		message = "spreading furious love #collaborate_to_thrive"
	} else {
		message = msg
	}

	splittedCoins := coins / len(names)
	if splittedCoins == 0 {
		return nil, errInsufficientCoins
	}
	formattedNames := ""

	for _, n := range names {
		formattedNames = formattedNames + "@" + n + " "
	}
	data := BonusPayload{
		Reason: fmt.Sprintf("+%d %s%s", splittedCoins, formattedNames, message),
	}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return body, nil
}
