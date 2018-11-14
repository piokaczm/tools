package bonusly

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type MeResponse struct {
	Result struct {
		GivingBalance int    `json:"giving_balance"`
		Username      string `json:"username"`
	} `json:"result"`
}

func (c *Client) me() (*MeResponse, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+meEndpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.execute(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	meResponse := MeResponse{}
	err = json.Unmarshal(body, &meResponse)
	if err != nil {
		return nil, err
	}

	return &meResponse, nil
}
