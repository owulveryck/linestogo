package linestogo

import (
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
)

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

func (p *Page) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	svg := xml.StartElement{
		Name: xml.Name{Local: "svg"},
		Attr: []xml.Attr{
			{
				Name:  xml.Name{Local: "width"},
				Value: "1404",
			},
			{
				Name:  xml.Name{Local: "height"},
				Value: "1872",
			},
			{
				Name:  xml.Name{Local: "viewBox"},
				Value: "0 0 1404 1872",
			},
			{
				Name:  xml.Name{Local: "xmlns"},
				Value: "http://www.w3.org/2000/svg",
			},
		},
	}
	err := e.EncodeToken(svg)
	if err != nil {
		return err
	}
	for i, layer := range p.Layers {
		e.EncodeElement(layer, group(fmt.Sprintf("layer_%v", i)))
	}
	return e.EncodeToken(svg.End())
}

type test struct {
	A string
}

func group(name string) xml.StartElement {
	return xml.StartElement{
		Name: xml.Name{Local: "g"},
		Attr: []xml.Attr{
			{
				Name:  xml.Name{Local: "id"},
				Value: name,
			},
		},
	}
}
