package parser

import (
	"regexp"
)

const (
	KeywordName     = "Keyword"
	keywordAttrName = "KeywordAttr"
)

var (
	keywordRegexp = regexp.MustCompile(`^(\s*)#\+([^:]+):(\s+(.*)|\n|$)`)
)

type WithKeyword struct {
	Caption   map[string][]string
	HTMLAttrs map[string][]string
	Node   Node
}

type Keyword struct {
	Key   string
	Value string
}

type KeywordAttr struct {
	Key   string
	Value map[string]string
}

func (Keyword) Name() string {
	return KeywordName
}

func (s *Keyword) Get(string) (string, bool) {
	return "", true
}

func (s *parser) ParseKeyword(d *Document, lines []string) (*Keyword, int) {
	match := keywordRegexp.FindStringSubmatch(lines[0])
	if match == nil {
		return nil, 0
	}
	node := &Keyword{
		Key:   match[2],
		Value: match[4],
	}
	switch node.Key {
	case "CAPTION", "ATTR_HTML":
		// next, n := s.Parse(lines[1:])
		// if next != nil {

		// }
	default:
		d.Set(match[2], match[4])
	}
	return node, 1
}
