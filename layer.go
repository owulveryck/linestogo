package linestogo

import (
	"encoding/binary"
	"io"
)

type Layer struct {
	Strokes []*Stroke
}

// readFrom r and fill l
func (l *Layer) readFrom(r io.Reader) error {
	var numStrokes int32
	err := binary.Read(r, binary.LittleEndian, &numStrokes)
	if err != nil {
		return err
	}
	l.Strokes = make([]*Stroke, numStrokes)
	for i := 0; i < int(numStrokes); i++ {
		l.Strokes[i] = &Stroke{}
		err := l.Strokes[i].readFrom(r)
		if err != nil {
			return err
		}
	}
	return nil
}
