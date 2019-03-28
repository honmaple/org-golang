package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

// Block ..
type Block struct {
	Name      string
	FirstLine string
	LastLine  string
	Parent    *Block
	Children  []*Block
	Current   *Block
	NeedParse bool
	Escape    bool
	Regex     *regexp.Regexp
	Label     string
	Attr      map[string]string
}

func checkInList(f []string, v string) bool {
	for i := range f {
		if f[i] == v {
			return true
		}
	}
	return false
}

func headingEnd(s *Block, text string) bool {
	if s.Parent == nil {
		return false
	}
	if match := s.Regex.FindStringSubmatch(text); len(match) > 2 {
		if match1 := s.Regex.FindStringSubmatch(s.Parent.FirstLine); len(match1) > 2 && len(match[1]) <= len(match1[1]) {
			return true
		}
		return false
	}
	return false
}

func listEnd(s *Block, text string) bool {
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

func tableEnd(s *Block, text string) bool {
	if s.Match(text) {
		return false
	}
	return true
	return tablesep.MatchString(text)
}

func headingHTML(s *Block) string {
	strs := make([]string, 0)
	if match := s.Regex.FindStringSubmatch(s.FirstLine); len(match) > 2 {
		strs = append(strs, fmt.Sprintf(s.Label, len(match[1]), match[2]))
	}
	for _, child := range s.Children {
		strs = append(strs, child.String())
	}
	return strings.Join(strs, "\n")
}

func srcHTML(s *Block) string {
	strs := make([]string, 0)
	language := "language"
	if match := s.Regex.FindStringSubmatch(s.FirstLine); len(match) > 2 {
		language = match[2]
	}
	for _, child := range s.Children {
		strs = append(strs, child.String())
	}
	return fmt.Sprintf(s.Label, language, strings.Join(strs, "\n"))
}

func inlineblockHTML(s *Block) string {
	inlinetext := &InlineText{
		Text:      s.FirstLine,
		NeedParse: s.NeedParse,
		Escape:    s.Escape,
	}
	if s.Label == "" {
		return inlinetext.HTML()
	}
	return fmt.Sprintf(s.Label, inlinetext.HTML())
}

// Open ..
func (s *Block) Open(firstline string) *Block {
	b := &Block{
		Name:      s.Name,
		Label:     s.Label,
		Regex:     s.Regex,
		NeedParse: s.NeedParse,
		FirstLine: firstline,
	}
	switch b.Name {
	case table.Name:
		b.AddChild(tablerow.Open(b.FirstLine))
	case tablerow.Name:
		match := b.Regex.FindStringSubmatch(b.FirstLine)
		for _, i := range strings.Split(match[1], "|") {
			b.AddChild(tablecell.Open(i))
		}
	case paragraph.Name:
		p := inlineblock.Open(firstline)
		p.NeedParse = p.NeedParse
		b.AddChild(p)
	case unorderlist.Name, orderlist.Name:
		match := b.Regex.FindStringSubmatch(b.FirstLine)
		title := match[3]
		b.AddChild(listitem.Open(title))
	case listitem.Name:
		b.AddChild(inlineblock.Open(b.FirstLine))
	}

	return b
}

// Close ..
func (s *Block) Close(lastline string) {
	s.LastLine = lastline
}

// String ..
func (s *Block) String() string {
	switch s.Name {
	case heading.Name:
		return headingHTML(s)
	case src.Name:
		return srcHTML(s)
	case example.Name:
		return srcHTML(s)
	case inlineblock.Name, tablecell.Name:
		return inlineblockHTML(s)
	default:
		strs := make([]string, 0)
		for _, child := range s.Children {
			strs = append(strs, child.String())
		}
		if s.Label == "" {
			return strings.Join(strs, "\n")
		}
		return fmt.Sprintf(s.Label, strings.Join(strs, "\n"))
	}
}

// Debug ..
func (s *Block) Debug() []string {
	strs := make([]string, 0)
	count := ""
	parent := s.Parent
	if parent != nil {
		count = count + "\n"
	}
	for {
		if parent == nil {
			break
		}
		parent = parent.Parent
		count = count + "\t"
	}
	for _, str := range s.Children {
		d := fmt.Sprintf("%s%s%s", count, str.Name, strings.Join(str.Debug(), ""))
		strs = append(strs, d)
	}
	return strs
}

// Match ..
func (s *Block) Match(text string) bool {
	if s.Regex != nil {
		return s.Regex.MatchString(text)
	}
	return false
}

// MatchEnd ..
func (s *Block) MatchEnd(text string) bool {
	switch s.Name {
	case heading.Name:
		return headingEnd(s, text)
	case unorderlist.Name, orderlist.Name:
		return listEnd(s, text)
	case table.Name:
		return tableEnd(s, text)
	}
	if regex, ok := regexs[s.Name]; ok {
		return regex.MatchString(text)
	}
	return false
}

// AddChild ..
func (s *Block) AddChild(block *Block) {
	block.Parent = s
	s.Children = append(s.Children, block)
}

// Append ..
func (s *Block) Append(text string) {
	if s.Current == nil {
		s.Current = s
	}

	p := s.Current

	for {
		if s.Current.Name != paragraph.Name {
			break
		}
		s.Current = s.Current.Parent
	}

	if s.Current.MatchEnd(text) {
		s.Current.Close(text)
		if !checkInList([]string{table.Name, unorderlist.Name, orderlist.Name}, s.Current.Name) {
			s.Current = s.Current.Parent
			return
		}
		s.Current = s.Current.Parent
	}

	if s.Current.Name == table.Name {
		s.Current.AddChild(tablerow.Open(text))
		return
	}
	if s.Current.Name == unorderlist.Name || s.Current.Name == orderlist.Name {
		match := s.Current.Regex.FindStringSubmatch(s.Current.FirstLine)
		depth := len(match[1])
		if s.Current.Match(text) {
			match1 := s.Current.Regex.FindStringSubmatch(text)
			depth1 := len(match1[1])
			title1 := match1[3]
			if depth == depth1 {
				s.Current.AddChild(listitem.Open(title1))
				return
			}
		}
		s.Current.Children[len(s.Current.Children)-1].Append(text)
		return
	}
	for _, block := range blocks {
		if block.Match(text) {
			b := block.Open(text)
			s.Current.AddChild(b)
			s.Current = b
			return
		}
	}
	if blankline.Match(text) {
		s.Current.AddChild(blankline.Open(text))
		return
	}
	if hr.Match(text) {
		s.Current.AddChild(hr.Open(text))
		return
	}

	if p.Name == paragraph.Name {
		b := inlineblock.Open(text)
		b.NeedParse = p.NeedParse
		p.AddChild(b)
		s.Current = p
		return
	}

	b := paragraph.Open(text)
	if !s.Current.NeedParse {
		b.Label = ""
		b.NeedParse = false
	}
	s.Current.AddChild(b)
	s.Current = b
	return
}

var regexs = map[string]*regexp.Regexp{
	"src":     regexp.MustCompile(`\s*#\+(END_SRC|end_src)\s*$`),
	"example": regexp.MustCompile(`\s*#\+(END_EXAMPLE|end_example)\s*$`),
	"center":  regexp.MustCompile(`\s*#\+(END_CENTER|end_center)\s*$`),
	"verse":   regexp.MustCompile(`\s*#\+(END_VERSE|end_verse)\s*$`),
	"export":  regexp.MustCompile(`\s*#\+(END_EXPORT|end_export)\s*$`),
	"quote":   regexp.MustCompile(`\s*#\+(END_QUOTE|end_quote)\s*$`),
}

var tablesep = regexp.MustCompile(`^(\s*)\|((?:\+|-)*?)\|?$`)

var org = &Block{
	Name:      "org",
	NeedParse: true,
}

var hr = &Block{
	Name:      "hr",
	Regex:     regexp.MustCompile(`^\s*\-{5,}\s*`),
	NeedParse: false,
	Label:     "<hr/>%[1]s",
}

var blankline = &Block{
	Name:      "blankline",
	Regex:     regexp.MustCompile(`^$`),
	NeedParse: false,
}

var heading = &Block{
	Name:      "heading",
	Regex:     regexp.MustCompile(`^(\*+)\s+(.+)$`),
	Label:     "<h%[1]d>%[2]s</h%[1]d>",
	NeedParse: true,
}

var src = &Block{
	Name:      "src",
	Regex:     regexp.MustCompile(`\s*#\+(BEGIN_SRC|begin_src)\s+(.+)$`),
	Label:     "<pre class=\"%[1]s\">\n%[2]s\n</pre>",
	NeedParse: false,
}
var example = &Block{
	Name:      "example",
	Regex:     regexp.MustCompile(`\s*#\+(BEGIN_EXAMPLE|begin_example)$`),
	Label:     "<pre class=\"%[1]s\">\n%[2]s\n</pre>",
	NeedParse: false,
}

var center = &Block{
	Name:      "center",
	Regex:     regexp.MustCompile(`\s*#\+(BEGIN_CENTER|begin_center)$`),
	Label:     "<p class=\"org-center\">\n%[1]s\n</p>",
	NeedParse: true,
}

var export = &Block{
	Name:      "export",
	Regex:     regexp.MustCompile(`\s*#\+(BEGIN_EXPORT|begin_export)\s+(.+)$`),
	Label:     "%[1]s",
	NeedParse: false,
}

var verse = &Block{
	Name:      "verse",
	Regex:     regexp.MustCompile(`\s*#\+(BEGIN_VERSE|begin_verse)$`),
	Label:     "<p class=\"org-verse\">\n%[1]s\n</p>",
	NeedParse: true,
}

var quote = &Block{
	Name:      "quote",
	Regex:     regexp.MustCompile(`\s*#\+(BEGIN_QUOTE|begin_quote)$`),
	Label:     "<blockquote>\n%[1]s\n</blockquote>",
	NeedParse: true,
}

var paragraph = &Block{
	Name:      "paragraph",
	NeedParse: true,
	Label:     "<p>\n%[1]s\n</p>",
}

var inlineblock = &Block{
	Name:      "inlineblock",
	NeedParse: true,
	Label:     "%[1]s",
}

var listitem = &Block{
	Name:      "listitem",
	Label:     "<li>%[1]s</li>",
	NeedParse: true,
}

var orderlist = &Block{
	Name:      "orderlist",
	Regex:     regexp.MustCompile(`(\s*)\d+(\.|\))\s+(.+)$`),
	Label:     "<ul>\n%[1]s\n</ul>",
	NeedParse: true,
}

var unorderlist = &Block{
	Name:      "unorderlist",
	Regex:     regexp.MustCompile(`(\s*)(-|\+)\s+(.+)$`),
	Label:     "<ul>\n%[1]s\n</ul>",
	NeedParse: true,
}

var table = &Block{
	Name:      "table",
	Regex:     regexp.MustCompile(`\s*\|(.+?)\|*$`),
	Label:     "<table>\n%[1]s\n</table>",
	NeedParse: true,
}

var tablerow = &Block{
	Name:      "tablerow",
	Regex:     regexp.MustCompile(`\s*\|(.+?)\|*$`),
	Label:     "<tr>\n%[1]s\n</tr>",
	NeedParse: true,
}

var tablecell = &Block{
	Name:      "tablecell",
	Regex:     regexp.MustCompile(`\s*\|(.+?)\|*$`),
	Label:     "<td>%[1]s</td>",
	NeedParse: true,
}

var blocks = []*Block{
	heading,
	table,
	unorderlist,
	orderlist,
	src,
	example,
	export,
	center,
	verse,
	quote,
}

// OrgToHTML ..
func OrgToHTML(text string) string {
	var buffer bytes.Buffer

	for _, str := range strings.Split(text, "\n") {
		org.Append(str)
	}
	for _, str := range org.Children {
		s := str.String()
		fmt.Println(s)
		buffer.WriteString(s)
	}
	// for _, i := range org.Debug() {
	//	fmt.Println(i)
	// }
	return buffer.String()
}
