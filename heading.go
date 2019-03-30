package org

import (
	// "crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

// Toc ..
type Toc struct {
	Block
}

// Heading ..
type Heading struct {
	Block

	ID    string
	Title string
	Level int
}

var toc = &Toc{
	Block: Block{
		Name:      "toc",
		NeedParse: true,
		Label:     "<div id=\"table-of-contents\">\n<h2>Table of Contents</h2>\n<div id=\"text-table-of-contents\">\n%[1]s\n</div></div>\n\n",
	},
}

var heading = &Heading{
	Block: Block{
		Name:      "heading",
		Regex:     regexp.MustCompile(`^(\*+)\s+(.+)$`),
		Label:     "<h%[1]d id=\"%[3]s\">%[2]s</h%[1]d>",
		NeedParse: true,
	},
}

// headingID ..
func headingID(text string) string {
	m := sha1.New()
	m.Write([]byte(text))
	return hex.EncodeToString(m.Sum(nil))[:8]
}

// MatchEnd ..
func (s *Heading) MatchEnd(text string) bool {
	if s.Parent == nil || !s.Regex.MatchString(text) {
		return false
	}
	if match := s.Regex.FindStringSubmatch(text); len(match) > 2 && len(match[1]) <= s.Level {
		return true
	}
	return false
}

// HTML ..
func (s *Heading) HTML() string {
	strs := make([]string, 0)

	strs = append(strs, fmt.Sprintf(s.Label, s.Level + options.Offset, s.Title, s.ID))

	for _, child := range s.Children {
		strs = append(strs, child.HTML())
	}
	return strings.Join(strs, "\n")
}

// Open ..
func (s *Toc) Open(firstline string) BlockType {
	return &Toc{Block: *s.open(firstline)}
}

// HTML ..
func (s *Toc) HTML() string {
	if len(s.Children) == 0 {
		return ""
	}
	m := org.Open("")
	for _, child := range s.Children {
		head := child.(*Heading)
		m.Append(fmt.Sprintf("%[1]s- <a href=\"#%[3]s\">%[2]s</a>", strings.Repeat(" ", head.Level), head.Title, head.ID))
	}
	return fmt.Sprintf(s.Label, m.HTML())
}

// Open ..
func (s *Heading) Open(firstline string) BlockType {
	match := s.Regex.FindStringSubmatch(firstline)
	b := &Heading{
		Block: *s.open(firstline),
		Level: len(match[1]),
		Title: match[2],
		ID:    headingID(firstline),
	}
	toc.AddChild(b)
	return b
}
