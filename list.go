package main

import (
	"regexp"
	"strings"
)

// UnorderList ..
type UnorderList struct {
	Block
}

// OrderList ..
type OrderList struct {
	Block
}

// ListItem ..
type ListItem struct {
	Block
}

var listregex = regexp.MustCompile(`(\s*)(.+)$`)

var orderlist = &OrderList{
	Block: Block{
		Name:      "orderlist",
		Regex:     regexp.MustCompile(`(\s*)\d+(\.|\))\s+(.+)$`),
		Label:     "<ul>\n%[1]s\n</ul>",
		NeedParse: true,
	},
}

var unorderlist = &UnorderList{
	Block: Block{
		Name:      "unorderlist",
		Regex:     regexp.MustCompile(`(\s*)(-|\+)\s+(.+)$`),
		Label:     "<ul>\n%[1]s\n</ul>",
		NeedParse: true,
	},
}

var listitem = &ListItem{
	Block: Block{
		Name:      "listitem",
		Label:     "<li>%[1]s</li>",
		NeedParse: true,
	},
}

// Open ..
func (s *ListItem) Open(firstline string) BlockType {
	b := &ListItem{
		Block: *s.open(firstline),
	}
	b.AddChild(inlineblock.Open(firstline))
	return b
}

// Open ..
func (s *UnorderList) Open(firstline string) BlockType {
	b := &UnorderList{
		Block: *s.open(firstline),
	}
	match := s.Regex.FindStringSubmatch(firstline)
	title := match[3]
	b.AddChild(listitem.Open(title))
	return b
}

// Append ..
func (s *UnorderList) Append(text string) {
	match := s.Regex.FindStringSubmatch(s.FirstLine)
	depth := len(match[1])
	if s.Match(text) {
		match1 := s.Regex.FindStringSubmatch(text)
		depth1 := len(match1[1])
		title1 := match1[3]
		if depth == depth1 {
			s.AddChild(listitem.Open(title1))
			return
		}
	}
	s.Children[len(s.Children)-1].Append(text)
}

// MatchEnd ..
func (s *UnorderList) MatchEnd(text string) bool {
	if text == "" {
		return false
	}
	match := s.Regex.FindStringSubmatch(s.FirstLine)
	depth := len(match[1])
	if !s.Match(text) && depth >= len(text)-len(strings.TrimSpace(text)) {
		return true
	}
	return false
}

// Open ..
func (s *OrderList) Open(firstline string) BlockType {
	b := &OrderList{
		Block: *s.open(firstline),
	}
	match := b.Regex.FindStringSubmatch(firstline)
	title := match[3]
	b.AddChild(listitem.Open(title))
	return b
}

// Append ..
func (s *OrderList) Append(text string) {
	match := s.Regex.FindStringSubmatch(s.FirstLine)
	depth := len(match[1])
	if s.Match(text) {
		match1 := s.Regex.FindStringSubmatch(text)
		depth1 := len(match1[1])
		title1 := match1[3]
		if depth == depth1 {
			s.AddChild(listitem.Open(title1))
			return
		}
	}
	s.Children[len(s.Children)-1].Append(text)
}

// MatchEnd ..
func (s *OrderList) MatchEnd(text string) bool {
	if text == "" {
		return false
	}
	match := s.Regex.FindStringSubmatch(s.FirstLine)
	depth := len(match[1])
	if !s.Match(text) && depth >= len(text)-len(strings.TrimSpace(text)) {
		return true
	}
	return false
}
