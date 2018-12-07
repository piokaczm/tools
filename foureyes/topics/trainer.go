package topics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"

	prose "gopkg.in/jdkato/prose.v2"
)

// As described here: https://medium.com/errata-ai/prodigy-prose-radically-efficient-machine-teaching-in-go-93389bf2d772
func TeachAboutEntities(path, name string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	train, test := Split(ReadAnnotations(data))

	// Here, we're training a new model named <name> with the training portion
	// of our annotated data.
	//
	// Depending on your hardware, this should take around 1 - 3 minutes.
	model := prose.ModelFromData(name, prose.UsingEntities(train))

	// Now, let's test our model:
	correct := 0.0
	for _, entry := range test {
		// Create a document without segmentation, which isn't required for NER.
		doc, err := prose.NewDocument(
			entry.Text,
			prose.WithSegmentation(false),
			prose.UsingModel(model))

		if err != nil {
			return err
		}
		ents := doc.Entities()

		if entry.Answer != "accept" && len(ents) == 0 {
			// If we rejected this entity during annotation, prose shouldn't
			// have labeled it.
			correct++
		} else {
			// Otherwise, we need to verify that we found the correct entities.
			expected := []string{}
			for _, span := range entry.Spans {
				expected = append(expected, entry.Text[span.Start:span.End])
			}
			if reflect.DeepEqual(expected, ents) {
				correct++
			}
		}
	}
	fmt.Println(test)
	fmt.Printf("correct %f, len: %d\n", correct, len(test))
	fmt.Printf("Correct (%%): %f\n", correct/float64(len(test)))

	savePath := fmt.Sprintf("entities/%s", name)
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		model.Write(savePath) // Save the model to disk.
	}
	return nil
}

func ReadAnnotations(jsonLines []byte) []Annotations {
	dec := json.NewDecoder(bytes.NewReader(jsonLines))
	entries := []Annotations{}
	for {
		ent := Annotations{}
		err := dec.Decode(&ent)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		entries = append(entries, ent)
	}
	return entries
}

type Annotations struct {
	Text   string                `json:"text"`
	Spans  []prose.LabeledEntity `json:"spans"`
	Answer string                `json:"answer"`
}

type Span struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// Split divides our human-annotated data set into two groups: one for training
// our model and one for testing it.
//
// We're using an 80-20 split here, although you may want to use a different
// split.
func Split(data []Annotations) ([]prose.EntityContext, []Annotations) {
	cutoff := int(float64(len(data)) * 0.8)

	train, test := []prose.EntityContext{}, []Annotations{}
	for i, entry := range data {
		if i < cutoff {
			train = append(train, prose.EntityContext{
				Text:   entry.Text,
				Spans:  entry.Spans,
				Accept: entry.Answer == "accept"})
		} else {
			test = append(test, entry)
		}
	}

	return train, test
}
