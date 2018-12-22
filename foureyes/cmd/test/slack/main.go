package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/piokaczm/tools/foureyes/config"

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

	conf := config.New()
	err := conf.ReadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("creating slack API")
	s, err := slack.New(conf.SlackConfig.ApiToken)
	if err != nil {
		log.Fatal(err)
	}

	n, err := s.NewNotifier(conf.SlackUsername)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("creating cleaner")
	cleaner := topics.NewCleaner()
	cleaner.BuildPipeline(
		topics.OnlyWithNouns,
		topics.NotShorterThan(3),
		topics.WithStemming,
		topics.WithLemmatizing,
	)

	// TODO: make initialization dynamic basing on input characteristics
	t := topics.New(
		5,
		15,
		cleaner,
	)

	log.Println("cleaning topics to watch for")
	// TODO: create map for proper human readable topics within notifications?
	topicsToWatchFor, err := cleaner.Clean(conf.Topics)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("parsed topics: ", topicsToWatchFor)

	log.Println("creating comparator")
	c := comparator.New(topicsToWatchFor)

	log.Println("creating glasses watcher")
	w, err := glasses.New(n)
	if err != nil {
		log.Fatal(err)
	}

	for _, ch := range conf.SlackConfig.Channels {
		log.Printf("creating pooler for %s\n", ch.Name)
		pooler, err := s.NewChannelPooler(ch.Name, t)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("registering pooler")
		err = w.Register(pooler, c, ch.IntervalTime)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("starting workers")
	w.Start()
	log.Println("program is running!")
	<-stop
	w.Stop()
}
