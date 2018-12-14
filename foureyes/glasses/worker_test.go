package glasses

import (
	"errors"
	"testing"
	"time"

	"github.com/piokaczm/tools/foureyes/glasses/mocks"
	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {
	newData := [][]string{[]string{"new data"}}
	topics := [][]string{[]string{"topics"}}
	msg := "msg"

	testCases := []struct {
		name      string
		mockCalls func(*mocks.Source, *mocks.Notifier, *mocks.TopicExtractor, *mocks.Comparator)
		assert    func(*testing.T, error, *mocks.Source, *mocks.Notifier, *mocks.TopicExtractor, *mocks.Comparator)
	}{
		{
			name: "happy path",
			mockCalls: func(s *mocks.Source, n *mocks.Notifier, e *mocks.TopicExtractor, c *mocks.Comparator) {
				s.On("FetchNewData").Return(newData, nil)
				e.On("Process", newData).Return(topics, nil)
				c.On("Match", topics).Return(msg, nil)
				n.On("Notify", msg).Return(nil)
			},
			assert: func(t *testing.T, err error, s *mocks.Source, n *mocks.Notifier, e *mocks.TopicExtractor, c *mocks.Comparator) {
				assert.NoError(t, err)
				s.AssertCalled(t, "FetchNewData")
				e.AssertCalled(t, "Process", newData)
				c.AssertCalled(t, "Match", topics)
				n.AssertCalled(t, "Notify", msg)
			},
		},
		{
			name: "with error",
			mockCalls: func(s *mocks.Source, n *mocks.Notifier, e *mocks.TopicExtractor, c *mocks.Comparator) {
				s.On("FetchNewData").Return(nil, errors.New("source error"))
			},
			assert: func(t *testing.T, err error, s *mocks.Source, n *mocks.Notifier, e *mocks.TopicExtractor, c *mocks.Comparator) {
				assert.EqualError(t, err, "source error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(*testing.T) {
			source := &mocks.Source{}
			notifier := &mocks.Notifier{}
			extractor := &mocks.TopicExtractor{}
			comparator := &mocks.Comparator{}

			w, err := newWorker(source, notifier, extractor, comparator)
			if err != nil {
				t.Errorf("expected no errors, got %s", err)
			}

			tc.mockCalls(source, notifier, extractor, comparator)
			err = w.process()
			tc.assert(t, err, source, notifier, extractor, comparator)
		})
	}
}

func TestLifetime(t *testing.T) {
	newData := [][]string{[]string{"new data"}}
	topics := [][]string{[]string{"topics"}}
	msg := "msg"

	s := &mocks.Source{}
	n := &mocks.Notifier{}
	e := &mocks.TopicExtractor{}
	c := &mocks.Comparator{}

	w, err := newWorker(s, n, e, c)
	if err != nil {
		t.Errorf("expected no errors, got %s", err)
	}

	s.On("Interval").Return(100 * time.Millisecond)
	s.On("FetchNewData").Return(newData, nil)
	e.On("Process", newData).Return(topics, nil)
	c.On("Match", topics).Return(msg, nil)
	n.On("Notify", msg).Return(nil)

	go w.start()
	time.Sleep(500 * time.Millisecond)
	w.exit()
}
