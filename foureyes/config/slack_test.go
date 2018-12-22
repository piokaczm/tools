package config

// func TestSlackParse(t *testing.T) {
// 	t.Run("happy path", func(t *testing.T) {
// 		path := "./fixtures/config.yml"
// 		conf := New()

// 		err := conf.ReadConfig(path)
// 		assert.Nil(t, err)

// 		srs := conf.SlackChannels
// 		fmt.Println(srs)
// 		assert.Len(t, srs, 2)
// 		assert.Equal(t, srs[0].Name, "random")
// 		assert.Equal(t, srs[0].IntervalTime, 5*time.Second)
// 		assert.Equal(t, srs[1].Name, "general")
// 		assert.Equal(t, srs[1].IntervalTime, 10*time.Second)
// 	})
// }
