package topics

import (
	"github.com/pkg/errors"

	"github.com/james-bowman/nlp"
)

type Extractor struct {
	topicsNum int
	wordsNum  int
	cleaner   *Cleaner
}

func New(topicsNum, wordsNum int, c *Cleaner) *Extractor {
	return &Extractor{
		topicsNum: topicsNum,
		wordsNum:  wordsNum,
		cleaner:   c,
	}
}

func (e *Extractor) Process(document []string) ([][]string, error) {
	document, err := e.cleaner.Clean(document)
	if err != nil {
		return nil, err
	}

	vectoriser := nlp.NewCountVectoriser()
	lda := nlp.NewLatentDirichletAllocation(e.topicsNum)
	pipeline := nlp.NewPipeline(vectoriser, lda)

	_, err = pipeline.FitTransform(document...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to model topics")
	}

	for i := 0; i < 100; i++ {
		_ = lda.Components()
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
