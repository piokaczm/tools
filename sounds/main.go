// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

// +build ignore

package main

import (
	"log"
	"math"
	"time"

	ui "github.com/gizak/termui"
	"github.com/gizak/termui/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	var idx int
	data := make([][]float64, 1)
	data[0] = make([]float64, 150, 150)
	dPipeOut := make(chan [][]float64, 1)
	dPipeIn := make(chan [][]float64, 1)

	p0 := widgets.NewPlot()
	p0.Data = data
	p0.SetRect(0, 0, 180, 30)
	p0.AxesColor = ui.ColorWhite
	p0.LineColors[0] = ui.ColorGreen
	ui.Render(p0)
	dPipeIn <- p0.Data
	uiEvents := ui.PollEvents()

	go func() {
		for {
			select {
			case data := <-dPipeIn:
				d2 := make([][]float64, 1)
				d2[0] = data[0][0 : len(data[0])-1]
				// fmt.Println(d2[0][0])
				// fmt.Println(d2[0][0:10])
				// fmt.Println(idx)
				d2[0] = append([]float64{1 + math.Sin(float64(idx)/5)}, d2[0]...)
				idx++
				time.Sleep(25 * time.Millisecond)
				dPipeOut <- d2
			}
		}
	}()

	go func() {
		for {
			select {
			case data := <-dPipeOut:
				p0.Data = data
				ui.Render(p0)
				dPipeIn <- p0.Data
			}
		}
	}()

	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
