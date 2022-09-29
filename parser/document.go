package parser

import (
	"bufio"
	"io"
	"strings"
	"sync"
)

const (
	todoKeywords    = "TODO | DONE | CANCELED"
	timestampFormat = "2006-01-02 Mon 15:04"
)

var parserPool = sync.Pool{
	New: func() interface{} {
		return new(parser)
	},
}

type Node interface {
	Name() string
}

type Parser interface {
	Parse(*Document, []string) (Node, int)
	ParseAll(*Document, []string, bool) []Node
	ParseInline(*Document, string, int) (Node, int)
	ParseAllInline(*Document, string, bool) []Node
}

func isBlankline(line string) bool {
	return strings.TrimLeft(line, " ") == ""
}

func lineIndent(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}

type Document struct {
	Children        []Node
	Sections        *Section
	Keywords        map[string]string
	Properties      map[string]string
	TimestampFormat string
}

func (d *Document) Get(k string) string {
	return d.Keywords[k]
}

func (d *Document) Set(k, v string) {
	d.Keywords[k] = v
}

func ParseFromLines(lines []string) *Document {
	p := parserPool.Get().(*parser)
	defer parserPool.Put(p)

	d := &Document{
		Sections: &Section{},
		Keywords: map[string]string{
			"TODO": todoKeywords,
		},
		TimestampFormat: timestampFormat,
	}
	d.Children = p.ParseAll(d, lines, false)
	return d
}

func ParseFromReader(reader io.Reader) *Document {
	scanner := bufio.NewScanner(reader)

	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return ParseFromLines(lines)
}

func Parse(content string) *Document {
	return ParseFromLines(strings.Split(content, "\n"))
}

type parser struct{}

func (s *parser) Parse(d *Document, lines []string) (Node, int) {
	if node, idx := s.BlankLine(d, lines); node != nil {
		return node, idx
	}
	if node, idx := s.Headline(d, lines); node != nil {
		return node, idx
	}
	if node, idx := s.Table(d, lines); node != nil {
		return node, idx
	}
	if node, idx := s.List(d, lines); node != nil {
		return node, idx
	}
	if node, idx := s.Drawer(d, lines); node != nil {
		return node, idx
	}
	if node, idx := s.Block(d, lines); node != nil {
		return node, idx
	}
	if node, idx := s.BlockResult(d, lines); node != nil {
		return node, idx
	}
	if node, idx := s.Keyword(d, lines); node != nil {
		return node, idx
	}
	if node, idx := s.Hr(d, lines); node != nil {
		return node, idx
	}
	return nil, 0
}

func (s *parser) ParseAll(d *Document, lines []string, raw bool) []Node {
	if raw && len(lines) > 0 {
		return s.ParseAllInline(d, strings.Join(lines, "\n"), raw)
	}
	idx, end, nodes := 0, len(lines), make([]Node, 0)
	for idx < end {
		if node, i := s.Parse(d, lines[idx:]); node != nil {
			idx = idx + i
			nodes = append(nodes, node)
			continue
		}
		node, next, i := s.Paragragh(d, lines[idx:])
		if node != nil {
			nodes = append(nodes, node)
		}
		if next != nil {
			nodes = append(nodes, next)
		}
		idx = idx + i + 1
	}
	return nodes
}
