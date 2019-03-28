package main

import (
	"fmt"
)

// Paragraph ..
type Paragraph struct {
	Block
}

// InlineBlock ..
type InlineBlock struct {
	Block
}

var paragraph = &Paragraph{
	Block: Block{
		Name:      "paragraph",
		NeedParse: true,
		Label:     "<p>\n%[1]s\n</p>",
	},
}

var inlineblock = &InlineBlock{
	Block: Block{
		Name:      "inlineblock",
		NeedParse: true,
		Label:     "%[1]s",
	},
}

// HTML ..
func (s *InlineBlock) HTML() string {
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
func (s *InlineBlock) Open(firstline string) BlockType {
	return &InlineBlock{Block: *s.open(firstline)}
}

// Open ..
func (s *Paragraph) Open(firstline string) BlockType {
	b := &Paragraph{
		Block: *s.open(firstline),
	}
	p := inlineblock.Open(firstline)
	p.SetNeedParse(s.NeedParse)
	b.AddChild(p)
	return b
}
