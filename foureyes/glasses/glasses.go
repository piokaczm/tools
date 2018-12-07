package glasses

import "time"

type Watcher struct {
	sources        []Source       // APIs we ping for new data
	topicExtractor TopicExtractor // entity for extracting topics from new data
	notifier       Notifier       // entity for sending alerts
}

type Source interface {
	FetchNewData() ([]string, error)
}

type TopicExtractor interface {
	Process([][]string) ([][]string, error)
}

type Notifier interface {
	Notify(string) error
}

func (w *Watcher) Schedule(interval time.Duration) error {
	return nil
}
