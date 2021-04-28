package main

import (
	"compress/gzip"
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	linestogo "github.com/owulveryck/linesToGo"
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
	log.Fatal(http.ListenAndServe(":8080", nil))
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
</style>
<script src="" height="100%"></script>
<body>

<img src="/svg" >

</body>
</html>
`
