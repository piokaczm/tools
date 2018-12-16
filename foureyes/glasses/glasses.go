package glasses

import (
	"time"

	"github.com/piokaczm/tools/foureyes/comparator"
)

type Comparator interface {
	Match(comparator.Source) (string, error)
}

type Source interface {
	comparator.Source
}

type Notifier interface {
	Notify(string) error
}

type Watcher struct {
	workers  []*worker
	notifier Notifier // entity for sending alerts
}

func New(n Notifier) (*Watcher, error) {
	if n == nil {
		return nil, ErrNoNotifier
	}

	return &Watcher{make([]*worker, 0), n}, nil
}

func (w *Watcher) Register(s Source, c Comparator, i time.Duration) error {
	worker, err := newWorker(s, w.notifier, c, i)
	if err != nil {
		return err
	}

	w.workers = append(w.workers, worker)
	return nil
}

func (w *Watcher) Start() {
	for _, w := range w.workers {
		go w.start()
	}
}

func (w *Watcher) Stop() {
	for _, w := range w.workers {
		w.exit()
	}
}
