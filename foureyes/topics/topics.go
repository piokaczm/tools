package topics

import (
	"bufio"
	"os"

	"github.com/pkg/errors"

	"github.com/james-bowman/nlp"
)

type Extractor struct {
	topicsNum int
	wordsNum  int
	StopWords []string
}

type ExtractorOption func(*Extractor) error

func New(topicsNum, wordsNum int, opts ...ExtractorOption) *Extractor {
	e := &Extractor{
		topicsNum: topicsNum,
		wordsNum:  wordsNum,
		StopWords: make([]string, 0),
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

func WithStopWordsFromSlice(stopWords []string) ExtractorOption {
	return func(e *Extractor) error {
		e.StopWords = append(e.StopWords, stopWords...)
		return nil
	}
}

func WithStopWordsFromFile(filePath string) ExtractorOption {
	return func(e *Extractor) error {
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		var lines []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		e.StopWords = append(e.StopWords, lines...)
		return nil
	}
}

func (e *Extractor) Process(document []string) ([][]string, error) {
	vectoriser := nlp.NewCountVectoriser(e.StopWords...)
	lda := nlp.NewLatentDirichletAllocation(e.topicsNum)
	pipeline := nlp.NewPipeline(vectoriser, lda)

	_, err := pipeline.FitTransform(document...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to model topics")
	}

	topicsOverWords := lda.Components()
	tr, tc := topicsOverWords.Dims()

	vocab := make([]string, len(vectoriser.Vocabulary))
	for k, v := range vectoriser.Vocabulary {
		vocab[v] = k
	}

	aggr := newAggregator(tr)
	aggr.buildScoresMap(tc, topicsOverWords, vocab)
	return aggr.getTopWords(e.wordsNum), nil
}
