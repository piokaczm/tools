package glasses

import (
	"errors"
	"log"
	"time"
)

var (
	ErrNoSource     = errors.New("glasses: no source provided")
	ErrNoNotifier   = errors.New("glasses: no notifier provided")
	ErrNoComparator = errors.New("glasses: no comparator provided")
	ErrBadInterval  = errors.New("glasses: bad interval provided")
)

type worker struct {
	source     Source
	notifier   Notifier
	comparator Comparator
	interval   time.Duration
	stop       chan struct{}
}

func newWorker(s Source, n Notifier, c Comparator, i time.Duration) (*worker, error) {
	if s == nil {
		return nil, ErrNoSource
	}

	if n == nil {
		return nil, ErrNoNotifier
	}

	if c == nil {
		return nil, ErrNoComparator
	}

	if i.Nanoseconds() < 0 {
		return nil, ErrBadInterval
	}

	return &worker{s, n, c, i, make(chan struct{}, 1)}, nil
}

func (w *worker) start() {
	t := time.NewTicker(w.interval)

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
	msg, err := w.comparator.Match(w.source)
	if err != nil {
		return err
	}

	return w.notifier.Notify(msg)
}
