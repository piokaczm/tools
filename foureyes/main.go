package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/james-bowman/nlp"
	topics "github.com/patrikeh/go-topics"
)

// topics "github.com/patrikeh/go-topics"
// TODO: test with https://github.com/james-bowman/nlp/blob/master/lda.go

var (
	jiraAPIToken = os.Getenv("ATLASSIAN_API_TOKEN")
)

func main() {
	mainLDA()
	mainGoTopics()
}

func mainLDA() {
	corpus, err := readFile("./fixtures/comments.txt")
	if err != nil {
		log.Fatal(err)
	}

	stopWords, err := readFile("./fixtures/stopwords-en.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Create a pipeline with a count vectoriser and LDA transformer for 2 topics
	vectoriser := nlp.NewCountVectoriser(stopWords...)
	lda := nlp.NewLatentDirichletAllocation(2)
	pipeline := nlp.NewPipeline(vectoriser, lda)

	_, err = pipeline.FitTransform(corpus...)
	if err != nil {
		fmt.Printf("Failed to model topics for documents because %v", err)
		return
	}

	// Examine Document over topic probability distribution
	// dr, dc := docsOverTopics.Dims()
	// for doc := 0; doc < dc; doc++ {
	// fmt.Printf("\nTopic distribution for document: '%s' -\n", corpus[doc])
	// for topic := 0; topic < dr; topic++ {
	// if topic > 0 {
	// 	fmt.Printf(",")
	// }
	// fmt.Printf(" Topic #%d=%f\n", topic, docsOverTopics.At(topic, doc))
	// }
	// }

	// Examine Topic over word probability distribution
	topicsOverWords := lda.Components()
	tr, tc := topicsOverWords.Dims()
	mapByScore := make([]map[float64][]string, tr)
	scores := make([][]float64, tr)

	vocab := make([]string, len(vectoriser.Vocabulary))
	for k, v := range vectoriser.Vocabulary {
		vocab[v] = k
	}
	for topic := 0; topic < tr; topic++ {
		mapByScore[topic] = make(map[float64][]string)
		// fmt.Printf("\nWord distribution for Topic #%d -\n", topic)
		for word := 0; word < tc; word++ {
			score := topicsOverWords.At(topic, word)
			scores[topic] = append(scores[topic], score)

			if _, ok := mapByScore[topic][score]; !ok {
				mapByScore[topic][score] = make([]string, 0)
			}

			mapByScore[topic][score] = append(mapByScore[topic][score], vocab[word])

			// if word > 0 {
			// 	fmt.Printf(",")
			// }
			// fmt.Printf(" '%s'=%f\n", vocab[word], topicsOverWords.At(topic, word))
		}
	}

	for i, m := range mapByScore {
		fmt.Printf("\ntopic: %d\n", i)
		sort.Float64s(scores[i])
		topWords := 8
		for j := 0; j < topWords; j++ {
			score := scores[i][len(scores[i])-j-1]
			fmt.Printf("%s -> %f\n", m[score], score)
		}
	}
}

// zbierz username ze slacka i jiry -> zbuduj funkcje do skipowania
func mainGoTopics() {
	processor := topics.NewProcessor(
		topics.Transformations{
			topics.ToLower,
			topics.Sanitize,
			topics.MinLen,
			skipNames,
			skipSlackNames,
			skipTimestamps,
			topics.GetStopwordFilter("./fixtures/stopwords-en.txt"),
		},
	)
	corpus, err := processor.ImportSingleFileCorpus(topics.NewCorpus(), "./fixtures/slack.txt")
	if err != nil {
		log.Fatal(err)
	}

	lda := topics.NewLDA(&topics.Configuration{Verbose: false})
	err = lda.Init(corpus, 2, 0, 0)
	if err != nil {
		log.Fatal(err)
	}

	_, err = lda.Train(10000)
	if err != nil {
		log.Fatal(err)
	}
	lda.PrintTopWords(8)
	// words := lda.GetTopWords(5)
	// fmt.Println(words)
}

func skipNames(word string) (new string, keep bool) {
	names := []string{}
	if include(names, word) {
		return "", false
	}

	return word, true
}

func skipSlackNames(word string) (new string, keep bool) {
	if strings.Contains(word, "@") {
		return "", false
	}

	return word, true
}

func skipTimestamps(word string) (new string, keep bool) {
	r := regexp.MustCompile(`[\d\d:\d\d]`)
	if r.MatchString(word) {
		return "", false
	}

	return word, true
}

func include(a []string, e string) bool {
	for _, w := range a {
		if w == e {
			return true
		}
	}
	return false
}

func readFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
