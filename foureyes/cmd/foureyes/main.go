package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/piokaczm/tools/foureyes/comparator"
	"github.com/piokaczm/tools/foureyes/config"
	"github.com/piokaczm/tools/foureyes/glasses"
	"github.com/piokaczm/tools/foureyes/services/slack"
	"github.com/piokaczm/tools/foureyes/topics"
	"github.com/pkg/errors"
)

var (
	confFlag = flag.String("config", "config.yaml", "path to the config file")
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	conf, err := initConfig()
	if err != nil {
		log.Fatal(err)
	}

	s, err := slack.New(conf.SlackConfig.ApiToken)
	if err != nil {
		log.Fatal(err)
	}

	n, err := s.NewNotifier(conf.SlackUsername)
	if err != nil {
		log.Fatal(err)
	}

	cleaner := initCleaner()
	t := topics.New(cleaner)

	topicsToWatchFor, err := parseTopics(conf.Topics, cleaner)
	if err != nil {
		log.Fatal(err)
	}

	c := comparator.New(topicsToWatchFor)
	w, err := glasses.New(n)
	if err != nil {
		log.Fatal(err)
	}

	err = registerWorkers(conf, t, s, c, w)
	if err != nil {
		log.Fatal(err)
	}

	w.Start()
	<-stop
	w.Stop()
}

func initConfig() (*config.Config, error) {
	flag.Parse()

	conf := config.New()
	err := conf.ReadConfig(*confFlag)
	return conf, errors.Wrap(err, fmt.Sprintf("couldn't read config from %s", *confFlag))
}

// TODO: keep info about raw topics for readability
func parseTopics(rawTopics []string, cleaner *topics.Cleaner) ([]string, error) {
	topicsToWatchFor, err := cleaner.Clean(rawTopics)
	if err != nil {
		log.Fatal(err)
	}

	for i, t := range topicsToWatchFor {
		if t == "" {
			return nil, fmt.Errorf("topic %s was converted to empty string using current clean pipeline", rawTopics[i])
		}
	}

	return topicsToWatchFor, nil
}

func initCleaner() *topics.Cleaner {
	cleaner := topics.NewCleaner()
	// TODO: allow changing pipeline via config
	cleaner.BuildPipeline(
		topics.OnlyWithNouns,
		topics.NotShorterThan(3),
		topics.WithStemming,
		topics.WithLemmatizing,
	)
	return cleaner
}

func registerWorkers(conf *config.Config, t *topics.Extractor, s *slack.Client, c *comparator.Matcher, w *glasses.Watcher) error {
	for _, ch := range conf.SlackConfig.Channels {
		pooler, err := s.NewChannelPooler(ch.Name, t)
		if err != nil {
			return err
		}

		err = w.Register(pooler, c, ch.IntervalTime)
		if err != nil {
			return err
		}
	}
	return nil
}
