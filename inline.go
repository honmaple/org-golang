package org

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
	NeedParse bool
	Customize map[string]*InlineBlock
}

// _inline_regex = r'((?:^|\s|[\u4e00-\u9fa5])(?![/\\])){0}([^\s]*?|[^\s]+.*?[^\s]+)(?<![/\\]|\s){0}(\B|[\u4e00-\u9fa5])'

var inlineRegex = `(?:^|\s|[^a-zA-Z0-9<\\\*])(%[1]s([^\s].*?[^\s])%[1]s)\B`

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
	Regex: regexp.MustCompile(fmt.Sprintf(inlineRegex, `(?:\*\*|\/)`)),
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
	Regex: regexp.MustCompile(fmt.Sprintf(inlineRegex, "[`~]")),
}

var fn = &Inline{
	Name:  "fn",
	Label: "<sup><a id=\"fnr:%[1]s\" class=\"footref\" href=\"#fn.%[1]s\">%[1]s</a></sup>",
	Regex: regexp.MustCompile(`\[fn:([\w-]+?)(:(.*?))?\]`),
}
var link = &Inline{
	Name:  "link",
	Label: "<a href=\"%[2]s\">%[1]s</a>",
	Regex: regexp.MustCompile(`\[\[(.+?)\](?:\[(.+?)\])?\]`),
}
var image = &Inline{
	Name:  "image",
	Label: "<img src=\"%[1]s\"/>",
	Regex: regexp.MustCompile(`\[\[(.+?)\]\]`),
}

var inlines = []*Inline{
	comment,
	newline,
	link,
	image,
	italic,
	bold,
	underlined,
	code,
	delete,
	verbatim,
}

// HTML ..
func (s *Inline) HTML(text string) string {
	if s.Name == "link" {
		return s.Regex.ReplaceAllString(text, fmt.Sprintf(s.Label, "$1", "$2"))
	}
	return s.Regex.ReplaceAllString(text, fmt.Sprintf(s.Label, "$2"))
}

// Match ..
func (s *Inline) Match(text string) bool {
	return s.Regex.MatchString(text)
}

// HTML ..
func (s *InlineText) HTML() string {
	if !s.NeedParse {
		return s.Text
	}
	return InlineHTML(s.Text, options.Escape)
}

// InlineHTML ..
func InlineHTML(text string, escape bool) string {
	if escape {
		text = strings.ReplaceAll(text, "<", "&lt;")
		text = strings.ReplaceAll(text, ">", "&gt;")
	}
	for _, inline := range inlines {
		if inline.Match(text) {
			return InlineHTML(inline.HTML(text), false)
		}
	}
	return text
}
