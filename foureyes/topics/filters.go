package topics

import (
	"strings"

	"github.com/aaaton/golem"
	"github.com/kljensen/snowball"
	prose "gopkg.in/jdkato/prose.v2"
)

var (
	verbTags = map[string]struct{}{
		"VB":  struct{}{},
		"VBD": struct{}{},
		"VBG": struct{}{},
		"VBN": struct{}{},
		"VBP": struct{}{},
		"VBZ": struct{}{},
	}
	adverbTags = map[string]struct{}{
		"RB":  struct{}{},
		"RBR": struct{}{},
		"RBS": struct{}{},
		"RP":  struct{}{},
	}
	pronounTags = map[string]struct{}{
		"PRP":  struct{}{},
		"PRP$": struct{}{},
	}
	adjectiveTags = map[string]struct{}{
		"JJ":  struct{}{},
		"JJR": struct{}{},
		"JJS": struct{}{},
	}
)

func NotShorterThan(n int) Filter {
	return func(in prose.Token) (string, bool) {
		return in.Text, len(in.Text) <= n
	}
}

func WithStemming(in prose.Token) (string, bool) {
	stem, err := snowball.Stem(in.Text, "english", true)
	if err != nil {
		panic(err) // wtf
	}

	return stem, false
}

func WithLemmatizing(in prose.Token) (string, bool) {
	lem, err := golem.New("english")
	if err != nil {
		panic(err)
	}

	if !lem.InDict(in.Text) {
		return in.Text, false
	}

	return lem.Lemma(in.Text), false
}

func Downcase(in prose.Token) (string, bool) {
	return strings.ToLower(in.Text), false
}

func OnlyWithNouns(in prose.Token) (string, bool) {
	return in.Text, !strings.HasPrefix(in.Tag, "NN")
}

func WithoutPronouns(in prose.Token) (string, bool) {
	return banTag(pronounTags, in)
}

func WithoutAdjectives(in prose.Token) (string, bool) {
	return banTag(adjectiveTags, in)
}

func WithoutAdverbs(in prose.Token) (string, bool) {
	return banTag(adverbTags, in)
}

func WithoutVerbs(in prose.Token) (string, bool) {
	return banTag(verbTags, in)
}

func banTag(tags map[string]struct{}, in prose.Token) (string, bool) {
	_, ok := tags[in.Tag]
	return in.Text, ok
}
