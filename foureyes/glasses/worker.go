package glasses

import (
	"errors"
	"log"
	"time"
)

var (
	ErrNoSource    = errors.New("glasses: no source provided for worker")
	ErrNoNotifier  = errors.New("glasses: no notifier provided for worker")
	ErrNoExtractor = errors.New("glasses: no topic extractor provided for worker")
)

type worker struct {
	source     Source
	notifier   Notifier
	extractor  TopicExtractor
	comparator Comparator
	stop       chan struct{}
}

func newWorker(s Source, n Notifier, e TopicExtractor, c Comparator) (*worker, error) {
	if s == nil {
		return nil, ErrNoSource
	}

	if n == nil {
		return nil, ErrNoNotifier
	}

	if e == nil {
		return nil, ErrNoExtractor
	}

	return &worker{s, n, e, c, make(chan struct{}, 1)}, nil
}

func (w *worker) start() {
	t := time.NewTicker(w.source.Interval())

	for {
		select {
		case <-w.stop:
			defer func() { w.stop <- struct{}{} }()
			return
		case <-t.C:
			// TODO: add several retries before stopping it!
			err := w.process()
			if err != nil {
				log.Println(err) // todo: better logging!
			}
		}
	}
}

func (w *worker) exit() {
	w.stop <- struct{}{}
	<-w.stop
}

func (w *worker) process() error {
	newData, err := w.source.FetchNewData()
	if err != nil {
		return err
	}

	topics, err := w.extractor.Process(newData)
	if err != nil {
		return err
	}

	msg, err := w.comparator.Match(topics)
	if err != nil {
		return err
	}

	if msg == "" {
		return nil
	}

	return w.notifier.Notify(msg)
}
