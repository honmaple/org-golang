package parser

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"
)

const (
	InlineTextName      = "InlineText"
	InlineLinkName      = "Link"
	InlinePercentName   = "Percent"
	InlineEmphasisName  = "Emphasis"
	InlineFootnoteName  = "Footnote"
	InlineLineBreakName = "LineBreak"
	InlineTimestampName = "Timestamp"
)

var (
	plainLinkRegexp   = regexp.MustCompile(`^(\w+)://`)
	angleLinkRegexp   = regexp.MustCompile(`^<(\w+):(.+)>`)
	regularLinkRegexp = regexp.MustCompile(`^\[\[(.+?)\](?:\[(.+?)\])?\]`)
	commentRegexp     = regexp.MustCompile(`^(\s*)#(.*)$`)
	percentRegexp     = regexp.MustCompile(`^\[(\d+/\d+|\d+%)\]`)
	footnoteRegexp    = regexp.MustCompile(`^\[fn:([\w-]*?)(:(.*?))?\]`)
	timestampRegexp   = regexp.MustCompile(`^<(\d{4}-\d{2}-\d{2})( [A-Za-z]+)?( \d{2}:\d{2})?( \+\d+[dwmy])?>`)
)

type InlineText struct {
	Content string
}

func (InlineText) Name() string {
	return InlineTextName
}

type InlineLink struct {
	URL      string
	Desc     string
	Protocol string
}

func (InlineLink) Name() string {
	return InlineLinkName
}

func (s *InlineLink) IsImage() bool {
	ext := filepath.Ext(s.URL)
	img := map[string]bool{
		"png":  true,
		"jpg":  true,
		"jpeg": true,
		"gif":  true,
		"svg":  true,
	}
	return s.Desc == "" && img[ext]
}

func (s *InlineLink) IsVideo() bool {
	ext := filepath.Ext(s.URL)
	vid := map[string]bool{
		"mp4":  true,
		"webm": true,
	}
	return s.Desc == "" && vid[ext]
}

type InlineEmphasis struct {
	Marker   string
	Children []Node
}

func (InlineEmphasis) Name() string {
	return InlineEmphasisName
}

type InlineFootnote struct {
	Label      string
	Definition string
}

func (InlineFootnote) Name() string {
	return InlineFootnoteName
}

type InlinePercent struct {
	Num string
}

func (InlinePercent) Name() string {
	return InlinePercentName
}

type InlineLineBreak struct {
	Count int
}

func (InlineLineBreak) Name() string {
	return InlineLineBreakName
}

type InlineTimestamp struct {
	Time     time.Time
	IsDate   bool
	Interval string
}

func (InlineTimestamp) Name() string {
	return InlineTimestampName
}

func isSpace(line string, index int) bool {
	if index >= len(line) {
		return false
	}
	return line[index] == ' '
}

func isInList(w string, ws []string) bool {
	for _, word := range ws {
		if word == w {
			return true
		}
	}
	return false
}

func isValidPreBorder(line string, index int) bool {
	if index < 0 {
		return true
	}
	r := rune(line[index])
	return unicode.IsSpace(r) || strings.ContainsRune(`-({'"`, r)
}

func isValidPostBorder(line string, index int) bool {
	if index >= len(line) {
		return true
	}
	r := rune(line[index])
	return unicode.IsSpace(r) || strings.ContainsRune(`-.,:!?;'")}[`, r)
}

func (s *parser) ParseInlineLineBreak(d *Document, line string, i int) (*InlineLineBreak, int) {
	idx, end := i, len(line)
	for idx < end {
		if line[idx] != '\n' {
			break
		}
		idx++
	}
	if count := idx - i; count > 0 {
		return &InlineLineBreak{count}, count
	}
	return nil, 0
}

func (s *parser) ParseInlineTimestamp(d *Document, line string, i int) (*InlineTimestamp, int) {
	if m := timestampRegexp.FindStringSubmatch(line[i:]); m != nil {
		date, datetime, interval, isDate := m[1], m[3], strings.TrimSpace(m[4]), false
		if datetime == "" {
			datetime, isDate = "00:00", true
		}
		t, err := time.Parse(d.TimestampFormat, fmt.Sprintf("%s Mon %s", date, datetime))
		if err != nil {
			return nil, 0
		}
		return &InlineTimestamp{t, isDate, interval}, len(m[0])
	}
	return nil, 0
}

func (s *parser) ParseInlineFootnote(d *Document, line string, i int) (*InlineFootnote, int) {
	match := footnoteRegexp.FindStringSubmatch(line[i:])
	if len(match) == 0 {
		return nil, 0
	}
	return &InlineFootnote{Label: match[1], Definition: match[3]}, len(match[0])
}

func (s *parser) ParseInlinePercent(d *Document, line string, i int) (*InlinePercent, int) {
	match := percentRegexp.FindStringSubmatch(line[i:])
	if len(match) == 0 {
		return nil, 0
	}
	return &InlinePercent{match[1]}, len(match[0])
}

func (s *parser) ParseInlineLink(d *Document, line string, i int) (*InlineLink, int) {
	match := plainLinkRegexp.FindStringSubmatch(line[i:])
	if len(match) > 0 && isInList(match[1], d.Hyperlinks) {
		start, idx := i+len(match[0]), i+len(match[0])
		for idx < len(line) {
			if unicode.IsSpace(rune(line[idx])) {
				break
			}
			idx++
		}
		if idx > start {
			return &InlineLink{Protocol: match[1], URL: line[start:idx]}, idx
		}
	}
	match = angleLinkRegexp.FindStringSubmatch(line[i:])
	if len(match) > 0 && isInList(match[1], d.Hyperlinks) {
		return &InlineLink{Protocol: match[1], URL: match[2]}, len(match[0])
	}

	match = regularLinkRegexp.FindStringSubmatch(line[i:])
	if len(match) == 0 {
		return nil, 0
	}
	return &InlineLink{URL: match[1], Desc: match[2]}, len(match[0])
}

func (s *parser) ParseInlineEmphasis(d *Document, line string, i int) (*InlineEmphasis, int) {
	marker := line[i]

	needparse := true
	switch marker {
	case '*', '/', '+', '_':
		needparse = true
	case '=', '~', '`':
		needparse = false
	default:
		return nil, 0
	}

	if !isValidPreBorder(line, i-1) {
		return nil, 0
	}
	idx, end := i+1, len(line)
	for idx < end {
		if line[idx] == marker && idx != i+1 && isValidPostBorder(line, idx+1) {
			b := &InlineEmphasis{Marker: string(marker), Children: s.ParseAllInline(d, line[i+1:idx], !needparse)}
			// if isSpace(line, idx+1) {
			//	return b, idx - i + 2
			// }
			return b, idx - i + 1
		}
		idx++
	}
	return nil, 0
}

func (s *parser) ParseInlineText(d *Document, line string, i int) (Node, Node, int) {
	idx, end := i+1, len(line)
	for idx < end {
		if next, n := s.ParseInline(d, line, idx); next != nil {
			return &InlineText{line[i:idx]}, next, idx - i + n
		}
		idx++
	}
	return &InlineText{line[i:idx]}, nil, idx - i
}

func (s *parser) ParseInline(d *Document, line string, i int) (Node, int) {
	if node, idx := s.ParseInlineLineBreak(d, line, i); node != nil {
		return node, idx
	}
	if node, idx := s.ParseInlineEmphasis(d, line, i); node != nil {
		return node, idx
	}
	if node, idx := s.ParseInlineLink(d, line, i); node != nil {
		return node, idx
	}
	if node, idx := s.ParseInlinePercent(d, line, i); node != nil {
		return node, idx
	}
	if node, idx := s.ParseInlineFootnote(d, line, i); node != nil {
		return node, idx
	}
	if node, idx := s.ParseInlineTimestamp(d, line, i); node != nil {
		return node, idx
	}
	return nil, i
}

func (s *parser) ParseAllInline(d *Document, line string, raw bool) []Node {
	if raw {
		return []Node{&InlineText{line}}
	}
	idx, end, nodes := 0, len(line), make([]Node, 0)
	for idx < end {
		if node, i := s.ParseInline(d, line, idx); node != nil {
			nodes = append(nodes, node)
			idx = idx + i
			continue
		}
		node, next, i := s.ParseInlineText(d, line, idx)
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
