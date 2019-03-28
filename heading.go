package main

import (
	"fmt"
	"regexp"
	"strings"
)

// Heading ..
type Heading struct {
	Block
}

var heading = &Heading{
	Block: Block{
		Name:      "heading",
		Regex:     regexp.MustCompile(`^(\*+)\s+(.+)$`),
		Label:     "<h%[1]d>%[2]s</h%[1]d>",
		NeedParse: true,
	},
}

// MatchEnd ..
func (s *Heading) MatchEnd(text string) bool {
	if s.Parent == nil {
		return false
	}
	if match := s.Regex.FindStringSubmatch(text); len(match) > 2 {
		if match1 := s.Regex.FindStringSubmatch(s.Parent.GetFirstLine()); len(match1) > 2 && len(match[1]) <= len(match1[1]) {
			return true
		}
		return false
	}
	return false
}

// HTML ..
func (s *Heading) HTML() string {
	strs := make([]string, 0)
	if match := s.Regex.FindStringSubmatch(s.FirstLine); len(match) > 2 {
		strs = append(strs, fmt.Sprintf(s.Label, len(match[1]), match[2]))
	}
	for _, child := range s.Children {
		strs = append(strs, child.HTML())
	}
	return strings.Join(strs, "\n")
}

// Open ..
func (s *Heading) Open(firstline string) BlockType {
	return &Heading{Block: *s.open(firstline)}
}
