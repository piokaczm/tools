package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/piokaczm/tools/foureyes/comparator"
	"github.com/piokaczm/tools/foureyes/glasses"
	"github.com/piokaczm/tools/foureyes/services/slack"
	"github.com/piokaczm/tools/foureyes/topics"

	_ "github.com/joho/godotenv/autoload"
)

type notifier struct{}

func (n notifier) Notify(msg string) error {
	fmt.Println(msg)
	return nil
}

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	s, err := slack.New(os.Getenv("SLACK_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	n, err := s.NewNotifier("piokaczm")
	if err != nil {
		log.Fatal(err)
	}

	t := topics.New(
		2,
		15,
		topics.NewCleaner(
			topics.OnlyWithNouns,
		),
	)

	pooler, err := s.NewChannelPooler("random", t)
	if err != nil {
		log.Fatal(err)
	}

	c := comparator.New([]string{"Silesia", "Katowice", "Warsaw", "York"})
	w, err := glasses.New(n)
	if err != nil {
		log.Fatal(err)
	}

	err = w.Register(pooler, c, 40*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	w.Start()
	<-stop
	w.Stop()
}
