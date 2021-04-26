package main

import (
	"image"
	"image/draw"
	"reflect"
	"testing"

	"golang.org/x/image/vector"
)

func TestMe(t *testing.T) {
	z := vector.NewRasterizer(2, 2)
	z.DrawOp = draw.Src
	// . | .
	// -----
	// . | .

	z.MoveTo(0, 0)
	z.LineTo(1, 0)
	// x | x
	// -----
	// . | .

	dst := image.NewAlpha(z.Bounds())
	z.Draw(dst, dst.Bounds(), image.Opaque, image.Point{})
	expected := []uint8{255, 255, 0, 0}
	if !reflect.DeepEqual(dst.Pix, expected) {
		t.Fatal(dst.Pix)

	}
	/*

		p := &linestogo.Page{}
		err := linestogo.Read(os.Stdin, p)
		if err != nil {
			log.Fatal(err)
		}
		raster := vector.NewRasterizer(canvas.Dx(), canvas.Dy())
		dst := image.NewAlpha(raster.Bounds())
		for _, layer := range p.Layers {
			for _, stroke := range layer.Strokes {
				raster.Reset(canvas.Max.X, canvas.Max.Y)
				for i, segment := range stroke.Segments {
					if i == 0 {
						raster.MoveTo(float32(segment.P.X), float32(segment.P.Y))
					} else {
						raster.LineTo(float32(segment.P.X), float32(segment.P.Y))
					}
				}
				raster.Draw(dst, dst.Bounds(), image.Opaque, image.Point{})
			}
		}
		log.Println(dst.Pix)
		//png.Encode(os.Stdout, img)
	*/
}
