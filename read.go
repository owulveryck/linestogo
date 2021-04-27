package linestogo

import (
	"io"
)

//go:generate protoc --go_opt=Mproto/defs.proto=. --go_out=module=github.com/owulveryck/linestogo:. --proto_path=proto defs.proto

// Read from r and fill Page
func Read(r io.Reader, p *Page) error {
	return p.readFrom(r)
}
