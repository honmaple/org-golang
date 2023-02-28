package parser

import (
	"regexp"
	"strings"
)

const (
	HrName        = "Hr"
	FootnoteName  = "Footnote"
	BlanklineName = "Blankline"
	ParagraghName = "Paragragh"
)

var (
	hrRegexp        = regexp.MustCompile(`^\s*\-{5,}\s*`)
	footnoteRegexp  = regexp.MustCompile(`^\[fn:([\w-]*?)\]\s+(.*)$`)
	blanklineRegexp = regexp.MustCompile(`^(\s*)(?:\r?\n|$)`)
	plainTextRegexp = regexp.MustCompile(`^(\s*)(.*)`)
	attributeRegexp = regexp.MustCompile(`(?:^|\s+)(:[-\w]+)\s+(.*)$`)
)

type Footnote struct {
	Label      string
	Inline     bool
	Definition []Node
}

func (Footnote) Name() string {
	return FootnoteName
}

type Blankline struct {
	Count int
}

func (Blankline) Name() string {
	return BlanklineName
}

type Paragragh struct {
	Children []Node
}

func (Paragragh) Name() string {
	return ParagraghName
}

type Hr struct{}

func (Hr) Name() string {
	return HrName
}

func (s *parser) ParseHr(d *Document, lines []string) (*Hr, int) {
	match := hrRegexp.FindStringSubmatch(lines[0])
	if match == nil || len(match) == 0 {
		return nil, 0
	}
	return &Hr{}, 1
}

func (s *parser) ParseBlankLine(d *Document, lines []string) (*Blankline, int) {
	idx, end := 0, len(lines)
	for idx < end {
		if !isBlankline(lines[idx]) {
			break
		}
		idx++
	}
	if idx > 0 {
		return &Blankline{idx}, idx
	}
	return nil, 0
}

// footnote defintion no prfix space
func (s *parser) ParseFootnote(d *Document, lines []string) (*Footnote, int) {
	match := footnoteRegexp.FindStringSubmatch(lines[0])
	if match == nil || len(match) == 0 {
		return nil, 0
	}
	idx, end := 1, len(lines)
	for idx < end {
		if footnoteRegexp.MatchString(lines[idx]) || headingRegexp.MatchString(lines[idx]) {
			break
		}
		idx++
	}
	fn := &Footnote{
		Label:      match[1],
		Definition: s.ParseAll(d, append([]string{match[2]}, lines[1:idx]...), false),
	}
	return fn, idx
}

func (s *parser) ParseParagragh(d *Document, lines []string) (*Paragragh, Node, int) {
	idx, end := 1, len(lines)
	for idx < end {
		if next, n := s.Parse(d, lines[idx:]); next != nil {
			return &Paragragh{s.ParseAllInline(d, strings.Join(lines[:idx], "\n"), false)}, next, idx + n
		}
		idx++
	}
	return &Paragragh{s.ParseAllInline(d, strings.Join(lines[:idx], "\n"), false)}, nil, idx
}
