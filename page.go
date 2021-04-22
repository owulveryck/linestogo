package linestogo

import (
	"encoding/binary"
	"errors"
	"io"
	"strings"
)

func (p *Page) String() string {
	return p.Header
}

type Page struct {
	Header string
	Layers []*Layer
}

// readFrom r and fill p
func (p *Page) readFrom(r io.Reader) error {
	var buf struct {
		Header    [43]byte
		NumLayers int32
	}
	err := binary.Read(r, binary.LittleEndian, &buf)
	if err != nil {
		return err
	}
	p.Header = string(buf.Header[:])
	if !strings.Contains(p.Header, "version=5") {
		return errors.New("only version 5 is supported")
	}
	p.Layers = make([]*Layer, buf.NumLayers)
	for i := 0; i < int(buf.NumLayers); i++ {
		p.Layers[i] = &Layer{}
		err := p.Layers[i].readFrom(r)
		if err != nil {
			return err
		}
	}
	return nil
}
