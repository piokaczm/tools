package glasses

import "time"

type Watcher struct {
	sources        []Source       // APIs we ping for new data
	topicExtractor TopicExtractor // entity for extracting topics from new data
	notifier       Notifier       // entity for sending alerts
}

type Comparator interface {
	Match([][]string) (string, error)
}

type Source interface {
	Name() string
	FetchNewData() ([][]string, error)
	Interval() time.Duration
}

type TopicExtractor interface {
	Process([][]string) ([][]string, error)
}

type Notifier interface {
	Notify(string) error
}
