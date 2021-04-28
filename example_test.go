package linestogo_test

import (
	"context"
	"encoding/xml"
	"log"
	"os"

	linestogo "github.com/owulveryck/linesToGo"
)

func Example_svg() {
	p := &linestogo.Page{}
	err := linestogo.Read(os.Stdin, p)
	if err != nil {
		log.Fatal(err)
	}
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "    ")
	err = enc.Encode(p)
	if err != nil {
		log.Fatal(err)
	}
	enc.Flush()
}

func Example_startPolling() {
	ctx := context.Background()
	pageC, cancel, err := linestogo.StartPolling(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer cancel()
	for p := range pageC {
		log.Println("current page is:", p)
	}
}
