package bonusly

func (c *Client) MyBalance() (int, error) {
	me, err := c.me()
	if err != nil {
		return 0, err
	}

	return me.Result.GivingBalance, nil
}
