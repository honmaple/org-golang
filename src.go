package main

import (
	"fmt"
	"regexp"
	"strings"
)

// Src ..
type Src struct {
	Block
}

// Example ..
type Example struct {
	Block
}

var src = &Src{
	Block: Block{
		Name:      "src",
		Regex:     regexp.MustCompile(`\s*#\+(BEGIN_SRC|begin_src)\s+(.+)$`),
		Label:     "<pre class=\"%[1]s\">\n%[2]s\n</pre>",
		NeedParse: false,
	},
}
var example = &Example{
	Block: Block{
		Name:      "example",
		Regex:     regexp.MustCompile(`\s*#\+(BEGIN_EXAMPLE|begin_example)$`),
		Label:     "<pre class=\"%[1]s\">\n%[2]s\n</pre>",
		NeedParse: false,
	},
}

// HTML ..
func (s *Src) HTML() string {
	language := "language"
	if match := s.Regex.FindStringSubmatch(s.FirstLine); len(match) > 2 {
		language = match[2]
	}
	strs := make([]string, 0)
	for _, str := range s.Children {
		strs = append(strs, str.HTML())
	}
	return fmt.Sprintf(s.Label, language, strings.Join(strs, "\n"))
}

// HTML ..
func (s *Example) HTML() string {
	language := "language"
	strs := make([]string, 0)
	for _, str := range s.Children {
		strs = append(strs, str.HTML())
	}
	return fmt.Sprintf(s.Label, language, strings.Join(strs, "\n"))
}

// Open ..
func (s *Src) Open(firstline string) BlockType {
	return &Src{Block: *s.open(firstline)}
}

// Open ..
func (s *Example) Open(firstline string) BlockType {
	return &Example{Block: *s.open(firstline)}
}
