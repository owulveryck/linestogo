package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"math"
	"os"

	linestogo "github.com/owulveryck/linesToGo"
)

func main() {
	var colors = []string{"black", "grey", "white"}
	p := &linestogo.Page{}
	err := linestogo.Read(os.Stdin, p)
	if err != nil {
		log.Fatal(err)
	}
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "    ")
	svg := xml.StartElement{
		Name: xml.Name{Local: "svg"},
		Attr: []xml.Attr{
			{
				Name:  xml.Name{Local: "width"},
				Value: "1872",
			},
			{
				Name:  xml.Name{Local: "height"},
				Value: "1404",
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
	enc.EncodeToken(svg)
	for i, layer := range p.Layers {
		grp := xml.StartElement{
			Name: xml.Name{Local: "g"},
			Attr: []xml.Attr{
				{
					Name:  xml.Name{Local: "id"},
					Value: fmt.Sprintf("layer_%v", i),
				},
			},
		}
		enc.EncodeToken(grp)
		for j, stroke := range layer.Strokes {
			grp := xml.StartElement{
				Name: xml.Name{Local: "g"},
				Attr: []xml.Attr{
					{
						Name:  xml.Name{Local: "id"},
						Value: fmt.Sprintf("stroke_%v_%v", i, j),
					},
				},
			}
			enc.EncodeToken(grp)
			col := colors[stroke.StrokeColor]
			var pl *Polyline
			var previousWidth float64
			var i int
			class := fmt.Sprintf("pen%v", stroke.Pen)
			for _, segment := range stroke.Segments {
				if previousWidth != math.Round(float64(segment.Width)*100)/100 {
					i = 0
					p := make([][2]float32, 0)
					if pl != nil {
						enc.Encode(pl)
						p = [][2]float32{pl.Points[len(pl.Points)-1]}
					}
					pl = &Polyline{
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
			enc.EncodeToken(grp.End())
		}
		enc.EncodeToken(grp.End())
	}
	enc.EncodeToken(svg.End())
	enc.Flush()
}

type Polyline struct {
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
