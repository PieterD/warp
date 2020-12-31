package glutil

import (
	"fmt"
	"io"
)

type Generator struct {
	next func() ([]byte, error)

	buf    []byte
	bufPos int
	eof    bool
}

func NewGenerator(next func() ([]byte, error)) *Generator {
	return &Generator{
		next: next,
	}
}

func (g *Generator) Read(b []byte) (int, error) {
	if g.eof {
		return 0, io.EOF
	}
	if len(g.buf) > g.bufPos {
		src := g.buf[g.bufPos:]
		n := copy(b, src)
		g.bufPos += n
		if n == len(src) {
			g.bufPos = 0
			g.buf = g.buf[:0]
		}
		return n, nil
	}
	nextData, err := g.next()
	if err == io.EOF {
		g.eof = true
		return 0, io.EOF
	}
	if err != nil {
		return 0, fmt.Errorf("fetching next: %w", err)
	}
	n := copy(b, nextData)
	if n < len(nextData) {
		g.buf = append(g.buf, nextData[n:]...)
	}
	return n, nil
}

var _ io.Reader = &Generator{}
