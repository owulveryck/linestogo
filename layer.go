package linestogo

import (
	"encoding/binary"
	"encoding/xml"
	"fmt"
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

func (l *Layer) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	err := e.EncodeToken(start)
	if err != nil {
		return err
	}
	for i, stroke := range l.Strokes {
		err := e.EncodeElement(stroke, group(fmt.Sprintf("%v_stroke_%v", start.Attr[0].Value, i)))
		if err != nil {
			return err
		}
	}
	return e.EncodeToken(start.End())
}
