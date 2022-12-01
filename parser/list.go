package parser

import (
	"regexp"
	"strings"
)

const (
	ListName             = "List"
	ListItemName         = "ListItem"
	OrderlistName        = "OrderList"
	UnorderlistName      = "UnorderList"
	DescriptiveName      = "Descriptive"
	descriptiveDescName  = "DescriptiveDesc"
	descriptiveTitleName = "DescriptiveTitle"
	DescriptiveItemName  = "DescriptiveItem"
)

var (
	listRegexp        = regexp.MustCompile(`^(\s*)(([0-9]+|[a-zA-Z])[.)]|[+*-])(\s+(.*)|$)`)
	descriptiveRegexp = regexp.MustCompile(`^(\s*)([+*-])\s+(.*)::(\s|$)`)
	listStatusRegexp  = regexp.MustCompile(`\[( |X|-)\]\s`)
	levelRegexp       = regexp.MustCompile(`(\s*)(.+)$`)
)

type List struct {
	Type     string
	Level    int
	Children []Node
}

type ListItem struct {
	Level    int
	Bullet   string
	Status   string
	Title    string
	Children []Node
}

type DescriptiveItem struct {
	Level    int
	Bullet   string
	Status   string
	Title    string
	Descs    []Node
	Children []Node
}

func (List) Name() string {
	return ListName
}

func (ListItem) Name() string {
	return ListItemName
}

func (s ListItem) Kind() string {
	if strings.ContainsAny(s.Bullet, "-*+") {
		return UnorderlistName
	}
	return OrderlistName
}

func (DescriptiveItem) Name() string {
	return DescriptiveItemName
}

func (s *parser) ParseListItem(d *Document, lines []string) (*ListItem, int) {
	match := listRegexp.FindStringSubmatch(lines[0])
	if match == nil {
		return nil, 0
	}
	status, title := "", match[4]
	if m := listStatusRegexp.FindStringSubmatch(title); m != nil {
		status, title = m[1], title[len("[ ] "):]
	}
	b := &ListItem{
		Level:  len(match[1]),
		Title:  title,
		Status: status,
		Bullet: match[2],
	}
	spa := 0
	idx, end := 1, len(lines)
	for idx < end {
		if isBlankline(lines[idx]) {
			spa++
			// 连续两次空行
			if spa == 2 {
				break
			}
			idx++
			continue
		}
		spa = 0
		if level := lineIndent(lines[idx]); level <= b.Level {
			break
		}
		idx++
	}
	b.Children = s.ParseAll(d, append([]string{title}, lines[1:idx]...), false)
	return b, idx
}

func (s *parser) ParseList(d *Document, lines []string) (*List, int) {
	item, idx := s.ParseListItem(d, lines)
	if item == nil {
		return nil, 0
	}
	l := &List{
		Type:     item.Kind(),
		Children: []Node{item},
	}

	end := len(lines)
	for idx < end {
		if level := lineIndent(lines[idx]); level < item.Level {
			break
		}
		item, ln := s.ParseListItem(d, lines[idx:])
		if item != nil && item.Level == item.Level && item.Kind() == l.Type {
			l.Children = append(l.Children, item)
			idx = idx + ln
			continue
		}
		break
	}
	return l, idx
}
