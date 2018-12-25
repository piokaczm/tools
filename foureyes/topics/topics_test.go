package topics

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"testing"
)

func BenchmarkProcess(b *testing.B) {
	cleaner := NewCleaner()
	cleaner.BuildPipeline(
		OnlyWithNouns,
		NotShorterThan(3),
		WithStemming,
		WithLemmatizing,
	)
	top := New(cleaner)
	data := fixture("fixtures/article_paragraph.txt")

	for i := 0; i < b.N; i++ {
		_, err := top.Process(data)
		if err != nil {
			b.Errorf("unexpected error occured: %s", err)
		}
	}
}
func TestProcess(t *testing.T) {
	tcs := []struct {
		name          string
		document      []string
		shouldContain []string
	}{
		{
			"with wiki article",
			fixture("fixtures/wiki.txt"),
			[]string{"katowice", "silesia", "poland", "city"},
		},
		{
			"with article",
			fixture("fixtures/article.txt"),
			[]string{"brexit", "vote"},
			// TODO: add warning about bad topics within provided ones?
			// --- FAIL: TestProcess/with_article (29.93s)
			// topics_test.go:57: expected [[brexit deal vote]] to contain
			// topics_test.go:57: expected [[brexit deal vote]] to contain
			// topics_test.go:57: expected [[brexit deal vote]] to contain britain
		},
	}

	cleaner := NewCleaner()
	cleaner.BuildPipeline(
		OnlyWithNouns,
		NotShorterThan(3),
		WithStemming,
		WithLemmatizing,
	)
	top := New(cleaner)

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			expectedTopics, err := cleaner.Clean(tc.shouldContain)
			if err != nil {
				t.Fatalf("expected cleaning base topics to return no errors, got %s", err)
			}

			actual, err := top.Process(tc.document)
			if err != nil {
				t.Fatalf("expected no errors, got %s", err)
			}

			fmt.Printf("extracted topics: %v\n", actual)

			for _, expected := range expectedTopics {
				if containsTopic(expected, actual) {
					continue
				}

				t.Errorf("expected %v to contain %s", actual, expected)
			}
		})
	}
}

func containsTopic(topic string, extractedTopics [][]string) bool {
	for _, extracted := range extractedTopics {
		for _, e := range extracted {
			if topic == e {
				return true
			}
		}
	}

	return false
}

func fixture(path string) []string {
	var doc []string

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		doc = append(doc, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return doc
}
