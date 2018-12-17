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

	log.Println("creating slack API")
	s, err := slack.New(os.Getenv("SLACK_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	// n, err := s.NewNotifier("piokaczm")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	log.Println("creating cleaner")
	cleaner := topics.NewCleaner()
	cleaner.BuildPipeline(
		topics.OnlyWithNouns,
		topics.NotShorterThan(3),
		topics.WithStemming,
		topics.WithLemmatizing,
	)
	t := topics.New(
		5,
		15,
		cleaner,
	)

	log.Println("creating #random pooler")
	pooler, err := s.NewChannelPooler("random", t)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("cleaning topics to watch for")
	// TODO: create map for proper human readable topics within notifications?
	topicsToWatchFor, err := cleaner.Clean([]string{"Silesia", "Katowice", "Warsaw", "York"})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("parsed topics: ", topicsToWatchFor)

	log.Println("creating comparator")
	c := comparator.New(topicsToWatchFor)

	log.Println("creating glasses watcher")
	w, err := glasses.New(notifier{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("registering pooler")
	err = w.Register(pooler, c, 3*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("starting workers")
	w.Start()
	log.Println("program is running!")
	<-stop
	w.Stop()
}