package topics

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/james-bowman/nlp"
)

type Extractor struct {
	cleaner *Cleaner
}

func New(c *Cleaner) *Extractor {
	return &Extractor{
		cleaner: c,
	}
}

func (e *Extractor) Process(document []string) ([][]string, error) {
	if len(document) == 0 {
		return [][]string{}, nil
	}

	document, err := e.cleaner.Clean(document)
	if err != nil {
		return nil, err
	}

	wordsNum, topicsNum := e.determineTopicsAndWordsNum(document)

	vectoriser := nlp.NewCountVectoriser()
	lda := nlp.NewLatentDirichletAllocation(topicsNum)
	pipeline := nlp.NewPipeline(vectoriser, lda)

	_, err = pipeline.FitTransform(document...)
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
	return aggr.getTopWords(wordsNum), nil
}

func (e *Extractor) determineTopicsAndWordsNum(document []string) (wordsNum int, topicsNum int) {
	var words int

	for _, d := range document {
		words += len(strings.Split(d, " ")) // naive
	}

	wordsNum = int(float64(words) * 0.008)
	topicsNum = int(float64(words) * 0.002)

	fmt.Printf("words: %d | wNum: %d | tNum: %d\n", words, wordsNum, topicsNum)

	if topicsNum == 0 {
		topicsNum = 1
	}

	if wordsNum < 10 {
		wordsNum = 10
	}

	return
}
