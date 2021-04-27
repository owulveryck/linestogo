package linestogo_test

import (
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
