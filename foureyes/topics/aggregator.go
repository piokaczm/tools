package topics

import (
	"sort"

	"gonum.org/v1/gonum/mat"
)

type aggregator struct {
	mapsByScore []map[float64][]string
	scores      [][]float64
	length      int
}

func newAggregator(length int) *aggregator {
	return &aggregator{
		make([]map[float64][]string, length),
		make([][]float64, length),
		length,
	}
}

func (a *aggregator) buildScoresMap(tc int, topicsOverWords mat.Matrix, vocab []string) {
	for topic := 0; topic < a.length; topic++ {
		for word := 0; word < tc; word++ {
			score := topicsOverWords.At(topic, word)
			a.addScore(vocab[word], score, topic)
		}
	}
}

func (t *aggregator) addScore(word string, score float64, topicIdx int) {
	if t.mapsByScore[topicIdx] == nil {
		t.mapsByScore[topicIdx] = make(map[float64][]string)
	}

	if _, ok := t.mapsByScore[topicIdx][score]; !ok {
		t.mapsByScore[topicIdx][score] = make([]string, 0)
	}

	t.scores[topicIdx] = append(t.scores[topicIdx], score)
	t.mapsByScore[topicIdx][score] = append(t.mapsByScore[topicIdx][score], word)
}

func (a *aggregator) getTopWords(n int) [][]string {
	topWords := make([][]string, a.length)

	for i, m := range a.mapsByScore {
		sort.Float64s(a.scores[i])

		// TODO: normalize for wordsNum > len
		for j := 0; j < n; j++ {
			score := a.scores[i][len(a.scores[i])-j-1]
			if topWords[i] == nil {
				topWords[i] = make([]string, 0)
			}
			// fmt.Printf("words: %v, score: %9.9f\n", m[score], score)

			// TODO: update j if m[score] is more than one word
			topWords[i] = append(topWords[i], m[score]...) // there might be multiple words with the same score
		}
	}

	return topWords
}
