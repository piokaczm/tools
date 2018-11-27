package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/piokaczm/tools/spreadlove/bonusly"
)

var (
	apiKey   = os.Getenv("BONUSLY_API_KEY")
	lbCount  = flag.Int("lucky-folks", 0, "to how many random folks do you want to send some shiny coins?")
	tcLimit  = flag.Int("coins-limit", 0, "if you don't want to spend ALL your coins at once you can specify a limit using this one.")
	lbNames  = flag.String("lucky-names", "", "(strings separated by a comma without spaces) if you want to give TCs to a specific group of people instead of random folks, just provide this flag.")
	msg      = flag.String("message", "", "if you want to add you custom message for a bonus, provide it here.")
	findUser = flag.String("find-user", "", "if you're not sure about the username of one you want to send a bonus to, this option will provide you with all matching names")
)

func main() {
	c := bonusly.New(apiKey)

	flag.Parse()
	if *findUser != "" {
		out, err := c.FindName(*findUser)
		if err != nil {
			fmt.Println("ERROR: ", err)
			os.Exit(1)
		}

		fmt.Println("These are usernames matching passed string:")
		for i, name := range out {
			fmt.Printf("%d: %s\n", i, name)
		}
		os.Exit(0)
	}

	if *lbNames == "" && *lbCount < 1 {
		fmt.Println("ERROR: You have to at least provide 'lucky-names' or 'lucky-folks' flags!")
		help()
		os.Exit(1)
	}

	if *lbNames != "" && *lbCount > 0 {
		fmt.Println("ERROR: You have to provide only one of 'lucky-names' and 'lucky-folks' flags!")
		help()
		os.Exit(1)
	}

	out, err := c.SpreadFuriousLove(*lbCount, *tcLimit, *lbNames, *msg)
	if err != nil {
		fmt.Println("ERROR: ", err)
		os.Exit(1)
	}

	fmt.Println("Bonus sent! -> ", out)
	os.Exit(0)
}

func help() {
	fmt.Println(`
This app is here to help you spend all your coins and have some fun with it, you can use some flags to customize the outcome:
	
	'lucky-folks=<int>' -> <int> to how many random folks do you want to send some shiny coins?
	'lucky-names=<strings separated by a comma without spaces>' -> if you want to give TCs to a specific group of people instead of random folks, just provide this flag.
	'coins-limit=<int>' -> if you don't want to spend ALL your coins at once you can specify a limit using this one.
	'message=<string>' -> if you want to add you custom message for a bonus, provide it here.
	'find-user=<string>' -> if you're not sure about the username of one you want to send a bonus to, this option will provide you with all matching names`)
}
