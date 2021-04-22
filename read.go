package linestogo

import (
	"io"
)

// Read from r and fill Page
func Read(r io.Reader, p *Page) error {
	return p.readFrom(r)
}
