package parser

import (
	"regexp"
	"strings"
)

const (
	HrName        = "Hr"
	BlanklineName = "Blankline"
	ParagraghName = "Paragragh"
)

var (
	hrRegexp        = regexp.MustCompile(`^\s*\-{5,}\s*`)
	blanklineRegexp = regexp.MustCompile(`^(\s*)(?:\r?\n|$)`)
	plainTextRegexp = regexp.MustCompile(`^(\s*)(.*)`)
	attributeRegexp = regexp.MustCompile(`(?:^|\s+)(:[-\w]+)\s+(.*)$`)
)

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

func (s *parser) Hr(d *Document, lines []string) (*Hr, int) {
	match := hrRegexp.FindStringSubmatch(lines[0])
	if match == nil || len(match) == 0 {
		return nil, 0
	}
	return &Hr{}, 1
}

func (s *parser) BlankLine(d *Document, lines []string) (*Blankline, int) {
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

func (s *parser) Paragragh(d *Document, lines []string) (*Paragragh, Node, int) {
	idx, end := 1, len(lines)
	for idx < end {
		if next, n := s.Parse(d, lines[idx:]); next != nil {
			return &Paragragh{s.ParseAllInline(d, strings.Join(lines[:idx], "\n"), false)}, next, idx + n
		}
		idx++
	}
	return &Paragragh{s.ParseAllInline(d, strings.Join(lines[:idx], "\n"), false)}, nil, idx
}
