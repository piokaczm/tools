package oscillator

import (
	"fmt"
	"time"

	ui "github.com/gizak/termui"
	"github.com/gizak/termui/widgets"
)

// Renderer is an interface describing an entity that's able to render a stream of data.
type Renderer interface {
	Render(point float64)
}

type TermUIRenderer struct {
	buffer                      []float64
	pl                          *widgets.Plot
	in, out                     chan float64
	stopListener, stopRendering chan struct{}
	freq                        time.Duration // in ms
}

func NewTermUIRenderer(bufSize, freq, width, height int) (*TermUIRenderer, error) {
	// !!! dataset length should be eq to width !!!
	// TODO: add input validations
	buffer := make([]float64, bufSize, bufSize)

	// create widget
	p := widgets.NewPlot()
	p.SetRect(0, 0, width, height)
	p.AxesColor = ui.ColorWhite
	p.LineColors[0] = ui.ColorGreen

	// assign values to renderer
	t := &TermUIRenderer{}
	t.buffer = buffer
	t.pl = p
	t.in = make(chan float64, 1)
	t.out = make(chan float64, 1)
	t.freq = time.Duration(freq) * time.Millisecond

	// start goroutines here?
	return t, nil
}

func (t *TermUIRenderer) Listen() chan struct{} {
	stop := make(chan struct{}, 1)
	go func() {
		for {
			select {
			case point := <-t.in:
				// add point to the buffer
				fmt.Println(point)
			case <-stop:
			}
		}
	}()

	return stop
}
