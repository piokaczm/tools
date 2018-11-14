package bonusly

import (
	"io/ioutil"
	"net/http"
)

type UsersResponse struct {
	Users []User `json:"result"`
}

type User struct {
	Username string `json:"username"`
}

func (c *Client) users() ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+usersEndpoint, nil)
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

	return body, nil
}
