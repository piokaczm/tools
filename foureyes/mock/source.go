package mock

type Source struct {
	topics [][]string
	name   string
	err    error
}

func NewSource(name string, topics [][]string, err error) Source {
	return Source{topics, name, err}
}

func (s Source) Name() string {
	return s.name
}

func (s Source) Topics() ([][]string, error) {
	return s.topics, s.err
}
