package topics

import (
	"strings"
	"sync"

	"github.com/pkg/errors"
	prose "gopkg.in/jdkato/prose.v2"
)

type Cleaner struct {
	pipeline []Filter
	models   []*prose.Model // what to do with it?
}

type Filter func(input prose.Token) (output string, discard bool) // TODO: abstract prose.Token?

type CleanerOption func(c *Cleaner)

func NewCleaner(opts ...CleanerOption) *Cleaner {
	c := &Cleaner{}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Cleaner) BuildPipeline(filters ...Filter) {
	c.pipeline = append(c.pipeline, filters...)
}

func WithCustomModels(paths ...string) CleanerOption {
	return func(c *Cleaner) {
		models := make([]*prose.Model, len(paths))
		for i, path := range paths {
			models[i] = prose.ModelFromDisk(path)
		}
	}
}

func (c *Cleaner) Clean(docs []string) ([]string, error) {
	cleanDocs := make([]string, len(docs))
	wg := &sync.WaitGroup{}
	errs := make(chan error, len(docs))

	for i, doc := range docs {
		wg.Add(1)

		go func(doc string, i int, wg *sync.WaitGroup) {
			defer wg.Done()
			cleanString, err := c.cleanupSingleString(doc)
			if err != nil {
				errs <- err
				return
			}

			cleanDocs[i] = cleanString
		}(doc, i, wg)
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			return nil, err
		}
	}

	return cleanDocs, nil
}

func (c *Cleaner) cleanupSingleString(s string) (string, error) {
	opts := make([]prose.DocOpt, len(c.models)+1)
	for i, mod := range c.models {
		opts[i] = prose.UsingModel(mod)
	}
	opts[len(c.models)] = prose.WithSegmentation(false)

	doc, err := prose.NewDocument(
		s,
		opts...,
	)
	if err != nil {
		return "", errors.Wrap(err, "couldn't create new prose document")
	}

	var cleanStrings []string

TokenLoop:
	for _, tok := range doc.Tokens() {
		// 	// maybe just add a whitelist for a good start, to make it simpler?
		// 	if !strings.HasPrefix(tok.Tag, c.only) && tok.Label != "APPLICATION" && tok.Label != "B-GPE" { // todo: make it dynamic
		// 		continue
		// 	}
		var discard bool
		for _, filter := range c.pipeline {
			tok.Text, discard = filter(tok)
			if discard {
				continue TokenLoop
			}
		}

		cleanStrings = append(cleanStrings, tok.Text)
	}

	// not the exact representation, but should be close enough for later processing
	return strings.Join(cleanStrings, " "), nil
}
