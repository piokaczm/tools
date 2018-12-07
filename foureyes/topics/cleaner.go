package topics

import (
	"strings"

	"github.com/pkg/errors"
	prose "gopkg.in/jdkato/prose.v2"
)

var (
	verbTags      = []string{"VB", "VBD", "VBG", "VBN", "VBP", "VBZ"}
	adverbTags    = []string{"RB", "RBR", "RBS", "RP"}
	pronounTags   = []string{"PRP", "PRP$"}
	adjectiveTags = []string{"JJ", "JJR", "JJS"}
)

type Cleaner struct {
	bannedTags map[string]struct{}
	only       string
	models     []*prose.Model
}

type CleanerOption func(c *Cleaner)

func NewCleaner(opts ...CleanerOption) *Cleaner {
	c := &Cleaner{
		bannedTags: make(map[string]struct{}),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithCustomModels(paths ...string) CleanerOption {
	return func(c *Cleaner) {
		models := make([]*prose.Model, len(paths))
		for i, path := range paths {
			models[i] = prose.ModelFromDisk(path)
		}
	}
}

func OnlyWithNouns(c *Cleaner) {
	c.only = "NN"
}

func WithoutPronouns(c *Cleaner) {
	c.banTag(pronounTags)
}

func WithoutAdjectives(c *Cleaner) {
	c.banTag(adjectiveTags)
}

func WithoutAdverbs(c *Cleaner) {
	c.banTag(adverbTags)
}

func WithoutVerbs(c *Cleaner) {
	c.banTag(verbTags)
}

func (c *Cleaner) banTag(tagsList []string) {
	for _, tag := range tagsList {
		c.bannedTags[tag] = struct{}{}
	}
}

func (c *Cleaner) Clean(docs []string) ([]string, error) {
	cleanDocs := make([]string, len(docs))

	for i, doc := range docs {
		cleanString, err := c.cleanupSingleString(doc)
		if err != nil {
			return nil, err
		}

		cleanDocs[i] = cleanString
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

	for _, tok := range doc.Tokens() {
		if c.only == "" {
			if _, exists := c.bannedTags[tok.Tag]; exists {
				continue
			}
		} else {
			// maybe just add a whitelist for a good start, to make it simpler?
			if !strings.HasPrefix(tok.Tag, c.only) && tok.Label != "APPLICATION" && tok.Label != "B-GPE" { // todo: make it dynamic
				continue
			}
		}
		cleanStrings = append(cleanStrings, tok.Text)
	}

	// not the exact representation, but should be close enough for later processing
	return strings.Join(cleanStrings, " "), nil
}
