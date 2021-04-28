package linestogo

import (
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"io"
	"math"
)

var colors = []string{"black", "grey", "white"}

func (s *Stroke) readFrom(r io.Reader) error {
	var buf struct {
		Pen         int32
		Color       int32
		_           [4]byte
		Width       float32
		_           [4]byte
		NumSegments int32
	}
	err := binary.Read(r, binary.LittleEndian, &buf)
	if err != nil {
		return err
	}
	s.Pen = uint32(buf.Pen)
	s.StrokeColor = uint32(buf.Color)
	s.PenWidth = float32(buf.Width)
	s.Segments = make([]*Segment, buf.NumSegments)
	for i := 0; i < int(buf.NumSegments); i++ {
		s.Segments[i] = &Segment{}
		err := s.Segments[i].readFrom(r)
		if err != nil {
			return err
		}
	}
	return nil
}

// MarshalXML creates a SVG representation of the stroke
func (s *Stroke) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	err := e.EncodeToken(start)
	if err != nil {
		return err
	}
	col := colors[s.StrokeColor]
	var pl *polyline
	var previousWidth float64
	var i int
	class := fmt.Sprintf("pen%v", s.Pen)
	for _, segment := range s.Segments {
		if previousWidth != math.Round(float64(segment.Width)*100)/100 {
			i = 0
			p := make([][2]float32, 0)
			if pl != nil {
				pl.Points = append(pl.Points, [2]float32{segment.X, segment.Y})
				e.Encode(pl)
				p = [][2]float32{pl.Points[len(pl.Points)-1]}
			}
			pl = &polyline{
				Stroke:      col,
				Fill:        "none",
				Points:      p,
				Class:       class,
				StrokeWidth: strokeWidth(math.Round(float64(segment.Width)*100) / 100),
			}
			previousWidth = math.Round(float64(segment.Width)*100) / 100
		}
		pl.Points = append(pl.Points, [2]float32{segment.X, segment.Y})
		i++
	}
	return e.EncodeToken(start.End())
}

type polyline struct {
	XMLName     xml.Name    `xml:"polyline"`
	Points      points      `xml:"points,attr"`
	Stroke      string      `xml:"stroke,attr"`
	Fill        string      `xml:"fill,attr"`
	Class       string      `xml:"class,attr"`
	StrokeWidth strokeWidth `xml:"stroke-width,attr"`
}

type strokeWidth float64

func (s strokeWidth) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{
		Name:  name,
		Value: fmt.Sprintf("%vpx", s),
	}, nil
}

type points [][2]float32

func (p points) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	var val string
	for i := 0; i < len(p); i++ {
		val = fmt.Sprintf("%v, %.2f %.2f", val, p[i][0], p[i][1])
	}
	return xml.Attr{
		Name:  name,
		Value: val[2:],
	}, nil
}
