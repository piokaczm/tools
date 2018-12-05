package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/piokaczm/tools/foureyes/topics"
)

var (
	jiraAPIToken = os.Getenv("ATLASSIAN_API_TOKEN")
)

func main() {
	t := topics.New(
		1,
		10,
		topics.WithStopWordsFromFile("./fixtures/stopwords-en.txt"),
	)

	topics, err := t.Process(readFile("./fixtures/gop.txt"))
	if err != nil {
		log.Fatal(err)
	}

	for i, topic := range topics {
		fmt.Printf("[topic %d] - %q\n", i, topic)
	}
}

func readFile(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return lines
}
