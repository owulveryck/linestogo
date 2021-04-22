# linestogo

Tiny lib to manipulate the .line format (.rm in the reMarkable2) in Go

Example:

[embedmd]:# (example/cmd/main.go go)
```go
package main

import (
	"log"
	"os"

	"github.com/kr/pretty"
	linestogo "github.com/owulveryck/linesToGo"
)

func main() {
	p := &linestogo.Page{}
	err := linestogo.Read(os.Stdin, p)
	if err != nil {
		log.Fatal(err)
	}
	pretty.Print(p)
}
```
