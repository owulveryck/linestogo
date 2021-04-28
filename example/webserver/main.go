package main

import (
	"compress/gzip"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	linestogo "github.com/owulveryck/linesToGo"
	"gonum.org/v1/gonum/stat"
)

func main() {
	ctx := context.Background()
	pageC, cancel, err := linestogo.StartPolling(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer cancel()
	http.HandleFunc("/", frontPage)
	http.HandleFunc("/svg", serveSVG(pageC))
	http.HandleFunc("/svgcorrected", serveSVGCorrected(pageC))
	http.HandleFunc("/raw", serveRaw(pageC))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveRaw(pageC <-chan string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Try to get a fresh picture
		select {
		case page := <-pageC:
			f, err := os.Open(page)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer f.Close()
			w.Header().Set("content-type", "application/octet-stream")
			w.Header().Set("Content-Encoding", "gzip")
			zr := gzip.NewWriter(w)
			defer zr.Close()
			io.Copy(zr, f)
			return
		case <-time.After(10 * time.Millisecond):
		}
		http.Error(w, "no content", http.StatusOK)
	}

}
func serveSVG(pageC <-chan string) http.HandlerFunc {
	p := &linestogo.Page{}
	return func(w http.ResponseWriter, r *http.Request) {
		// Try to get a fresh picture
		select {
		case page := <-pageC:
			f, err := os.Open(page)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer f.Close()
			err = linestogo.Read(f, p)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case <-time.After(10 * time.Millisecond):
		}
		if p != nil {
			w.Header().Set("content-type", "image/svg+xml")
			w.Header().Set("Content-Encoding", "gzip")
			zr := gzip.NewWriter(w)
			defer zr.Close()
			enc := xml.NewEncoder(zr)
			enc.Indent("", "    ")
			err := enc.Encode(p)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			enc.Flush()
			return
		}
		http.Error(w, "no content", http.StatusOK)
	}
}
func serveSVGCorrected(pageC <-chan string) http.HandlerFunc {
	p := &linestogo.Page{}
	return func(w http.ResponseWriter, r *http.Request) {
		// Try to get a fresh picture
		select {
		case page := <-pageC:
			f, err := os.Open(page)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer f.Close()
			err = linestogo.Read(f, p)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			linearize(p)
		case <-time.After(10 * time.Millisecond):
		}
		if p != nil {
			w.Header().Set("content-type", "image/svg+xml")
			w.Header().Set("Content-Encoding", "gzip")
			zr := gzip.NewWriter(w)
			defer zr.Close()
			enc := xml.NewEncoder(zr)
			enc.Indent("", "    ")
			err := enc.Encode(p)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			enc.Flush()
			return
		}
		http.Error(w, "no content", http.StatusOK)
	}
}

func frontPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, skeleton)
}

const skeleton = `
<!DOCTYPE html>
<html lang="en">
<meta charset="UTF-8">
<title>reMarkable</title>
<meta name="viewport" content="width=device-width,initial-scale=1">
<style>
.fullscreen-map {
    position: absolute;
    top: 0;
    left: 0;
    height: 110vh;
    min-width: 1020px;
    object-fit: cover;
    object-position: 0;
    font-family: 'object-fit: cover;';
    z-index: -1;
	transform: rotate(90deg);
}

</style>
<script src="" height="100%"></script>
<body>
<div class="map">
  <object class="fullscreen-map" type="image/svg+xml" data="/svg">
	<img src="/svg" >
  </object>
</div>


</body>
</html>
`

func linearize(p *linestogo.Page) {
	for _, l := range p.Layers {
		for _, s := range l.Strokes {
			x := make([]float64, len(s.Segments))
			y := make([]float64, len(s.Segments))
			for i := 0; i < len(s.Segments); i++ {
				x[i] = float64(s.Segments[i].X)
				y[i] = float64(s.Segments[i].Y)
			}
			alpha, beta := stat.LinearRegression(x, y, nil, false)
			// replace the segments
			s.Segments[len(s.Segments)-1] = &linestogo.Segment{
				X: s.Segments[len(s.Segments)-1].X,
				Y: float32(beta)*s.Segments[len(s.Segments)-1].X + float32(alpha),
			}
			s.Segments[1] = s.Segments[len(s.Segments)-1]
			s.Segments = s.Segments[:2]
		}

	}
}
