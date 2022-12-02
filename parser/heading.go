package parser

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

const HeadingName = "Heading"

var (
	headingRegexp      = regexp.MustCompile(`^(\*+)\s+(.*?)(?:\r?\n|$)`)
	headingTitleRegexp = regexp.MustCompile(`^(?:\[#([A-C])\])?\s*(.+?)(?:\s+:(.+?):)?$`)
)

type Section struct {
	*Heading
	idx      string
	last     *Section
	parent   *Section
	Children []*Section
}

func (s *Section) add(node *Heading) string {
	var parent *Section

	if s.last == nil {
		parent = s
	} else if node.Stars > s.last.Stars {
		parent = s.last
	} else if node.Stars == s.last.Stars {
		parent = s.last.parent
	} else {
		parent = s.last.parent
		for parent.Heading != nil && node.Stars <= parent.Stars {
			parent = parent.parent
		}
	}
	sec := &Section{Heading: node, parent: parent}
	parent.Children = append(parent.Children, sec)
	if parent.Heading == nil {
		sec.idx = fmt.Sprintf("%d", len(parent.Children))
	} else {
		sec.idx = fmt.Sprintf("%s.%d", parent.idx, len(parent.Children))
	}
	s.last = sec
	return sec.idx
}

// STARS KEYWORD PRIORITY TITLE TAGS
type Heading struct {
	Index      string
	Stars      int
	Keyword    string
	Priority   string
	Title      []Node
	Tags       []string
	Properties *Drawer
	Children   []Node
}

func (Heading) Name() string {
	return HeadingName
}

func (s *Heading) Id() string {
	if s.Properties != nil {
		if id := s.Properties.Get("CUSTOM_ID"); id != "" {
			return id
		}
	}
	return fmt.Sprintf("heading-%s", s.Index)
}

func (s *parser) ParseHeading(d *Document, lines []string) (*Heading, int) {
	match := headingRegexp.FindStringSubmatch(lines[0])
	if len(match) == 0 {
		return nil, 0
	}
	title := match[2]
	keyword := ""
	if v := strings.SplitN(title, " ", 2); len(v) >= 2 {
		todo := strings.FieldsFunc(d.Get("TODO"), func(r rune) bool { return unicode.IsSpace(r) || r == '|' })
		for i := range todo {
			if v[0] == todo[i] {
				keyword = v[0]
				title = v[1]
				break
			}
		}
	}
	b := &Heading{
		Stars:   len(match[1]),
		Keyword: keyword,
	}
	b.Index = d.Sections.add(b)

	tmatch := headingTitleRegexp.FindStringSubmatch(title)
	b.Priority = tmatch[1]
	b.Title = s.ParseAllInline(d, tmatch[2], false)
	b.Tags = strings.FieldsFunc(tmatch[3], func(r rune) bool { return r == ':' })

	idx, end := 1, len(lines)
	for idx < end {
		if m := headingRegexp.FindStringSubmatch(lines[idx]); m != nil && len(m[1]) <= b.Stars {
			break
		}
		idx++
	}
	children := s.ParseAll(d, lines[1:idx], false)
	if len(children) > 0 && children[0].Name() == DrawerName {
		b.Properties = children[0].(*Drawer)
		children = children[1:]
	}
	b.Children = children

	return b, idx
}
