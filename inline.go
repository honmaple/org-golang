package main

import (
	"fmt"
	"regexp"
	"strings"
)

// Inline ..
type Inline struct {
	Name  string
	Label string
	Regex *regexp.Regexp
}

// InlineText ..
type InlineText struct {
	Text      string
	Escape    bool
	NeedParse bool
}

var inlineRegex = `(%[1]s(.*?)%[1]s)`

var comment = &Inline{
	Name:  "comment",
	Label: "%[1]s",
	Regex: regexp.MustCompile(`^(\s*)#(.*)$`),
}
var newline = &Inline{
	Name:  "newline",
	Label: "\n",
	Regex: regexp.MustCompile(`\\$`),
}
var italic = &Inline{
	Name:  "italic",
	Label: "<i>%[1]s</i>",
	Regex: regexp.MustCompile(fmt.Sprintf(inlineRegex, `\*\*`)),
}
var bold = &Inline{
	Name:  "bold",
	Label: "<b>%[1]s</b>",
	Regex: regexp.MustCompile(fmt.Sprintf(inlineRegex, `\*`)),
}
var underlined = &Inline{
	Name:  "underlined",
	Label: `<span style="text-decoration:underline">%[1]s</span>`,
	Regex: regexp.MustCompile(fmt.Sprintf(inlineRegex, `_`)),
}
var code = &Inline{
	Name:  "code",
	Label: "<code>%[1]s</code>",
	Regex: regexp.MustCompile(fmt.Sprintf(inlineRegex, `\=`)),
}
var delete = &Inline{
	Name:  "delete",
	Label: "<del>%[1]s</del>",
	Regex: regexp.MustCompile(fmt.Sprintf(inlineRegex, `\+`)),
}
var verbatim = &Inline{
	Name:  "verbatim",
	Label: "<code>%[1]s</code>",
	Regex: regexp.MustCompile(fmt.Sprintf(inlineRegex, `~`)),
}

var fn = &Inline{
	Name:  "fn",
	Label: "<code>%[1]s</code>",
	Regex: regexp.MustCompile(fmt.Sprintf(inlineRegex, `\=`)),
}
var link = &Inline{
	Name:  "link",
	Label: "<code>%[1]s</code>",
	Regex: regexp.MustCompile(fmt.Sprintf(inlineRegex, `\=`)),
}
var image = &Inline{
	Name:  "image",
	Label: "<code>%[1]s</code>",
	Regex: regexp.MustCompile(fmt.Sprintf(inlineRegex, `\=`)),
}

var inlines = []*Inline{
	comment,
	newline,
	italic,
	bold,
	underlined,
	code,
	delete,
	verbatim,
}

// HTML ..
func (s *Inline) HTML(text string) string {
	return s.Regex.ReplaceAllString(text, fmt.Sprintf(s.Label, "$2"))
}

// Match ..
func (s *Inline) Match(text string) bool {
	return s.Regex.FindString(text) != ""
}

// HTML ..
func (s *InlineText) HTML() string {
	if !s.NeedParse {
		return s.Text
	}
	return InlineHTML(s.Text, s.Escape)
}

// InlineHTML ..
func InlineHTML(text string, escape bool) string {
	if escape {
		text = strings.ReplaceAll(text, "<", "&lt;")
		text = strings.ReplaceAll(text, ">", "&gt;")
		// if quote {
		//	text = strings.ReplaceAll(text, `"`, "&quot;")
		//	text = strings.ReplaceAll(text, "'", "&#39;")
		// }
		// return text
	}
	for _, inline := range inlines {
		if inline.Match(text) {
			return InlineHTML(inline.HTML(text), false)
		}
	}
	return text
}
