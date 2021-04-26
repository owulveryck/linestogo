package linestogo

import (
	"encoding/binary"
	"image"
	"image/draw"
	"io"
)

type Stroke struct {
	Pen         uint32
	StrokeColor uint32
	PenWidth    float32
	Segments    []*Segment
}

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

// Draw aligns r.Min in dst with sp in src and then replaces the
// rectangle r in dst with the result of drawing src on dst.
func (s *Stroke) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	panic("not implemented") // TODO: Implement
}
