package linestogo

import (
	"encoding/binary"
	"io"
)

type Segment struct {
	X         float32
	Y         float32
	Speed     float32
	Direction float32
	Width     float32
	Pressure  float32
}

func (s *Segment) readFrom(r io.Reader) error {
	var buf struct {
		X         float32
		Y         float32
		Speed     float32
		Direction float32
		Width     float32
		Pressure  float32
	}
	err := binary.Read(r, binary.LittleEndian, &buf)
	if err != nil {
		return err
	}
	/*
		s.P.X = int(buf.X)
		s.P.Y = int(buf.Y)
		if !s.P.In(image.Rect(0, 0, 1404, 1872)) {
			return fmt.Errorf("point %v is out of bound", s.P)
		}
	*/
	s.X = buf.X
	s.Y = buf.Y
	s.Speed = buf.Speed
	s.Direction = buf.Direction
	s.Width = buf.Width
	s.Pressure = buf.Pressure
	return nil
}
