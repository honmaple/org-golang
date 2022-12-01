package parser

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

const HeadlineName = "Headline"

var (
	headlineRegexp      = regexp.MustCompile(`^(\*+)\s+(.*?)(?:\r?\n|$)`)
	headlineTitleRegexp = regexp.MustCompile(`^(?:\[#([A-C])\])?\s*(.+?)(?:\s+:(.+?):)?$`)
)

type Section struct {
	*Headline
	idx      string
	last     *Section
	parent   *Section
	Children []*Section
}

func (s *Section) add(node *Headline) string {
	var parent *Section

	if s.last == nil {
		parent = s
	} else if node.Stars > s.last.Stars {
		parent = s.last
	} else if node.Stars == s.last.Stars {
		parent = s.last.parent
	} else {
		parent = s.last.parent
		for parent.Headline != nil && node.Stars <= parent.Stars {
			parent = parent.parent
		}
	}
	sec := &Section{Headline: node, parent: parent}
	parent.Children = append(parent.Children, sec)
	if parent.Headline == nil {
		sec.idx = fmt.Sprintf("%d", len(parent.Children))
	} else {
		sec.idx = fmt.Sprintf("%s.%d", parent.idx, len(parent.Children))
	}
	s.last = sec
	return sec.idx
}

type Headline struct {
	Index      string
	Stars      int
	Keyword    string
	Priority   string
	Title      []Node
	Tags       []string
	Properties *Drawer
	Children   []Node
}

func (Headline) Name() string {
	return HeadlineName
}

func (s *Headline) Id() string {
	if s.Properties != nil {
		if id := s.Properties.Get("CUSTOM_ID"); id != "" {
			return id
		}
	}
	return fmt.Sprintf("headline-%s", s.Index)
}

func (s *parser) ParseHeadline(d *Document, lines []string) (*Headline, int) {
	match := headlineRegexp.FindStringSubmatch(lines[0])
	if match == nil {
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
	b := &Headline{
		Stars:   len(match[1]),
		Keyword: keyword,
	}
	b.Index = d.Sections.add(b)

	tmatch := headlineTitleRegexp.FindStringSubmatch(title)
	b.Priority = tmatch[1]
	b.Title = s.ParseAllInline(d, tmatch[2], false)
	b.Tags = strings.FieldsFunc(tmatch[3], func(r rune) bool { return r == ':' })

	idx, end := 1, len(lines)
	for idx < end {
		if m := headlineRegexp.FindStringSubmatch(lines[idx]); m != nil && len(m[1]) <= b.Stars {
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
