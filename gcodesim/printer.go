package gcodesim

import (
	"bytes"
	"fmt"
	"math"

	svg "github.com/ajstarks/svgo"
)

type Printer struct {
	Relative         bool
	RelativeExtruder bool

	HeadPosition    Position
	Offset          Position
	CurrentFeedRate float64

	FanSpeed            uint8
	ExtruderTemperature float64
	BedTemperature      float64

	Layers map[float64]*Layer // indexed by Z pos

	ExtrudedFilament float64
}

func NewPrinter() *Printer {
	return &Printer{
		Layers: make(map[float64]*Layer),
	}
}

type Position struct {
	X float64
	Y float64
	Z float64
	E float64
}

type Line struct {
	x0, y0, x1, y1 float64
	e              float64
}

type Layer struct {
	Lines []Line

	ExtrudeMovements int
}

func NewLayer() *Layer {
	var l Layer
	return &l
}

func (l *Layer) SVG() []byte {
	buf := new(bytes.Buffer)
	canvas := svg.New(buf)

	// Determine min/max x/y
	var xMin float64 = math.MaxFloat64
	var yMin float64 = math.MaxFloat64
	var xMax, yMax float64
	for _, l := range l.Lines {
		if l.x0 < xMin {
			xMin = l.x0
		}
		if l.y0 < yMin {
			yMin = l.y0
		}
		if l.x0 > xMax {
			xMax = l.x0
		}
		if l.y0 > yMax {
			yMax = l.y0
		}

		if l.x1 < xMin {
			xMin = l.x1
		}
		if l.y1 < yMin {
			yMin = l.y1
		}
		if l.x1 > xMax {
			xMax = l.x1
		}
		if l.y1 > yMax {
			yMax = l.y1
		}
	}

	canvas.Start(int(xMax-xMin)+1, int(yMax-yMin)+1)

	for _, l := range l.Lines {
		canvas.Path(fmt.Sprintf("M %v,%v L %v,%v", l.x0-xMin, l.y0-yMin, l.x1-xMin, l.y1-yMin), `stroke="black"`)
	}

	canvas.End()
	return buf.Bytes()
}
