package bonusly

import (
	"encoding/json"
	"errors"
	"math/rand"
	"strings"
	"time"
)

var errNotEnoughFolks = errors.New("you're trying to spread love to more people than your organization has!")

func (c *Client) SpreadFuriousLove(lbCount, tcToSpread int, lbNames, msg string) (string, error) {
	var namesArray []string

	me, err := c.me()
	if err != nil {
		return "", err
	}

	if lbNames != "" {
		namesArray = strings.Split(lbNames, ",")
	} else {
		rawUsers, err := c.users()
		if err != nil {
			return "", err
		}

		namesArray, err = parseUsers(rawUsers, me.Result.Username)
		if err != nil {
			return "", err
		}

		if len(namesArray) < lbCount {
			return "", errNotEnoughFolks
		}

		namesArray = randomUsers(lbCount, namesArray)
	}

	var coins int

	if tcToSpread == 0 {
		coins = me.Result.GivingBalance
	} else {
		coins = tcToSpread
	}

	return c.bonus(namesArray, coins, msg)
}

func parseUsers(rawData []byte, me string) ([]string, error) {
	parsed := UsersResponse{}
	err := json.Unmarshal(rawData, &parsed)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(parsed.Users))
	for i, u := range parsed.Users {
		if u.Username == me {
			continue
		}

		names[i] = u.Username
	}

	return names, nil
}

func randomUsers(count int, allNames []string) []string {
	rand.Seed(time.Now().Unix())
	names := make([]string, count)

	for i := 0; i < count; i++ {
		idx := rand.Intn(len(allNames))
		names[i] = allNames[idx]
		allNames = append(allNames[:idx], allNames[idx+1:]...)
	}

	return names
}
