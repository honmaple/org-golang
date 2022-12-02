package parser

import (
	"bufio"
	"io"
	"strings"
	"sync"
)

var pool = sync.Pool{
	New: func() interface{} {
		return new(parser)
	},
}

type (
	Node interface {
		Name() string
	}
	Parser interface {
		Parse(*Document, []string) (Node, int)
		ParseAll(*Document, []string, bool) []Node
		ParseInline(*Document, string, int) (Node, int)
		ParseAllInline(*Document, string, bool) []Node
	}
	Document struct {
		Children        []Node
		Sections        *Section
		Keywords        map[string]string
		Properties      map[string]string
		Hyperlinks      []string
		TimestampFormat string
	}
)

func isBlankline(line string) bool {
	return strings.TrimLeft(line, " ") == ""
}

func lineIndent(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}

func (d *Document) Get(k string) string {
	return d.Keywords[k]
}

func (d *Document) Set(k, v string) {
	d.Keywords[k] = v
}

func ParseFromLines(d *Document, lines []string) []Node {
	p := pool.Get().(*parser)
	defer pool.Put(p)
	return p.ParseAll(d, lines, false)
}

func ParseFromText(d *Document, text string) []Node {
	return ParseFromLines(d, strings.Split(text, "\n"))
}

func Parse(d *Document, r io.Reader) []Node {
	scanner := bufio.NewScanner(r)

	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return ParseFromLines(d, lines)
}

type parser struct{}

func (s *parser) Parse(d *Document, lines []string) (Node, int) {
	if node, idx := s.ParseBlankLine(d, lines); node != nil {
		return node, idx
	}
	if node, idx := s.ParseHeading(d, lines); node != nil {
		return node, idx
	}
	if node, idx := s.ParseTable(d, lines); node != nil {
		return node, idx
	}
	if node, idx := s.ParseList(d, lines); node != nil {
		return node, idx
	}
	if node, idx := s.ParseDrawer(d, lines); node != nil {
		return node, idx
	}
	if node, idx := s.ParseBlock(d, lines); node != nil {
		return node, idx
	}
	if node, idx := s.ParseBlockResult(d, lines); node != nil {
		return node, idx
	}
	if node, idx := s.ParseKeyword(d, lines); node != nil {
		return node, idx
	}
	if node, idx := s.ParseHr(d, lines); node != nil {
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
		node, next, i := s.ParseParagragh(d, lines[idx:])
		if node != nil {
			nodes = append(nodes, node)
		}
		if next != nil {
			nodes = append(nodes, next)
		}
		idx = idx + i
	}
	return nodes
}
