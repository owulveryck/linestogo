package linestogo

import (
	"encoding/binary"
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
		_           int32
		Width       float32
		_           int32
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
