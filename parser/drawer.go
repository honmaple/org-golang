package parser

import (
	"regexp"
)

const (
	DrawerName = "Drawer"
)

var (
	beginDrawerRegexp = regexp.MustCompile(`^(\s*):(\S+):\s*$`)
	endDrawerRegexp   = regexp.MustCompile(`^(?i)(\s*):END:\s*$`)
	propertyRegexp    = regexp.MustCompile(`^(\s*):(\S+):(\s+(.*)$|$)`)
)

type Drawer struct {
	Type       string
	Level      int
	Properties map[string]string
	Children   []Node
}

func (Drawer) Name() string {
	return DrawerName
}

func (s *Drawer) Get(key string) string {
	return s.Properties[key]
}

func (s *parser) ParseDrawer(d *Document, lines []string) (*Drawer, int) {
	match := beginDrawerRegexp.FindStringSubmatch(lines[0])
	if match == nil {
		return nil, 0
	}

	idx, end := 1, len(lines)
	for idx < end {
		if m := endDrawerRegexp.FindStringSubmatch(lines[idx]); m != nil {
			return &Drawer{
				Type:       match[2],
				Level:      len(match[1]),
				Properties: make(map[string]string),
				Children:   s.ParseAll(d, lines[1:idx], false),
			}, idx + 1
		}
		idx++
	}
	return nil, 0
}
