package main

import (
	"fmt"
	"regexp"
	"strings"
)

// BlockType ..
type BlockType interface {
	Match(text string) bool
	MatchEnd(text string) bool
	Append(text string)
	HTML() string
	Open(text string) BlockType
	Close(text string)
	AddChild(block BlockType)
	SetParent(block BlockType)
	SetNeedParse(needparse bool)
	SetLabel(text string)

	GetName() string
	GetNeedParse() bool
	GetFirstLine() string
	GetParent() BlockType
	Debug() []string
}

// Block ..
type Block struct {
	Name      string
	FirstLine string
	LastLine  string
	Parent    BlockType
	Children  []BlockType
	Current   BlockType
	NeedParse bool
	Escape    bool
	Regex     *regexp.Regexp
	Label     string
	Attr      map[string]string
}

// Center ..
type Center struct {
	Block
}

// Verse ..
type Verse struct {
	Block
}

// Export ..
type Export struct {
	Block
}

// Quote ..
type Quote struct {
	Block
}

// checkInList ..
func checkInList(f []string, v string) bool {
	for i := range f {
		if f[i] == v {
			return true
		}
	}
	return false
}

// GetName ..
func (s *Block) GetName() string {
	return s.Name
}

// GetNeedParse ..
func (s *Block) GetNeedParse() bool {
	return s.NeedParse
}

// GetFirstLine ..
func (s *Block) GetFirstLine() string {
	return s.FirstLine
}

// GetParent ..
func (s *Block) GetParent() BlockType {
	return s.Parent
}

// SetParent ..
func (s *Block) SetParent(block BlockType) {
	s.Parent = block
}

// SetNeedParse ..
func (s *Block) SetNeedParse(needparse bool) {
	s.NeedParse = needparse
}

// SetLabel ..
func (s *Block) SetLabel(label string) {
	s.Label = label
}

// open ..
func (s *Block) open(firstline string) *Block {
	return &Block{
		Name:      s.Name,
		Label:     s.Label,
		Regex:     s.Regex,
		NeedParse: s.NeedParse,
		FirstLine: firstline,
	}
}

// Open ..
func (s *Block) Open(firstline string) BlockType {
	return s.open(firstline)
}

// Close ..
func (s *Block) Close(lastline string) {
	s.LastLine = lastline
}

// HTML ..
func (s *Block) HTML() string {
	strs := make([]string, 0)
	for _, child := range s.Children {
		strs = append(strs, child.HTML())
	}
	if s.Label == "" {
		return strings.Join(strs, "\n")
	}
	return fmt.Sprintf(s.Label, strings.Join(strs, "\n"))
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
		parent = parent.GetParent()
		count = count + "\t"
	}
	for _, str := range s.Children {
		d := fmt.Sprintf("%s%s%s", count, str.GetName(), strings.Join(str.Debug(), ""))
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
	if regex, ok := regexs[s.Name]; ok {
		return regex.MatchString(text)
	}
	return false
}

// AddChild ..
func (s *Block) AddChild(block BlockType) {
	block.SetParent(s)
	s.Children = append(s.Children, block)
}

// Append ..
func (s *Block) Append(text string) {
	if s.Current == nil {
		s.Current = s
	}

	p := s.Current

	for {
		if s.Current.GetName() != paragraph.GetName() {
			break
		}
		s.Current = s.Current.GetParent()
	}

	if s.Current.MatchEnd(text) {
		s.Current.Close(text)
		if !checkInList([]string{table.Name, unorderlist.Name, orderlist.Name}, s.Current.GetName()) {
			s.Current = s.Current.GetParent()
			return
		}
		s.Current = s.Current.GetParent()
	}

	if checkInList([]string{table.Name, unorderlist.Name, orderlist.Name}, s.Current.GetName()) {
		s.Current.Append(text)
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

	if p.GetName() == paragraph.Name {
		b := inlineblock.Open(text)
		b.SetNeedParse(p.GetNeedParse())
		p.AddChild(b)
		s.Current = p
		return
	}

	b := paragraph.Open(text)
	if !s.Current.GetNeedParse() {
		b.SetLabel("")
		b.SetNeedParse(false)
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

var center = &Center{
	Block: Block{
		Name:      "center",
		Regex:     regexp.MustCompile(`\s*#\+(BEGIN_CENTER|begin_center)$`),
		Label:     "<p class=\"org-center\">\n%[1]s\n</p>",
		NeedParse: true,
	},
}

var export = &Export{
	Block: Block{
		Name:      "export",
		Regex:     regexp.MustCompile(`\s*#\+(BEGIN_EXPORT|begin_export)\s+(.+)$`),
		Label:     "%[1]s",
		NeedParse: false,
	},
}

var verse = &Verse{
	Block: Block{
		Name:      "verse",
		Regex:     regexp.MustCompile(`\s*#\+(BEGIN_VERSE|begin_verse)$`),
		Label:     "<p class=\"org-verse\">\n%[1]s\n</p>",
		NeedParse: true,
	},
}

var quote = &Quote{
	Block: Block{
		Name:      "quote",
		Regex:     regexp.MustCompile(`\s*#\+(BEGIN_QUOTE|begin_quote)$`),
		Label:     "<blockquote>\n%[1]s\n</blockquote>",
		NeedParse: true,
	},
}

var blocks = []BlockType{
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

// Open ..
func (s *Center) Open(firstline string) BlockType {
	return &Center{Block: *s.open(firstline)}
}

// Open ..
func (s *Verse) Open(firstline string) BlockType {
	return &Verse{Block: *s.open(firstline)}
}

// Open ..
func (s *Export) Open(firstline string) BlockType {
	return &Export{Block: *s.open(firstline)}
}

// Open ..
func (s *Quote) Open(firstline string) BlockType {
	return &Quote{Block: *s.open(firstline)}
}
