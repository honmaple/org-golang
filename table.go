package org

import (
	"fmt"
	"regexp"
	"strings"
)

// Table ..
type Table struct {
	Block
}

// TableHeader ..
type TableHeader struct {
	Block
}

// TableRow ..
type TableRow struct {
	Block
}

// TableCell ..
type TableCell struct {
	Block
}

var tablesep = regexp.MustCompile(`^(\s*)\|((?:\+|-)*?)\|?$`)

var table = &Table{
	Block: Block{
		Name:      "table",
		Regex:     regexp.MustCompile(`\s*\|(.+?)\|*$`),
		Label:     "<table>\n%[1]s\n</table>",
		NeedParse: true,
	},
}

var tableheader = &TableHeader{
	Block: Block{
		Name:      "tableheader",
		Regex:     regexp.MustCompile(`\s*\|(.+?)\|*$`),
		Label:     "<th>%[1]s</th>",
		NeedParse: true,
	},
}

var tablerow = &TableRow{
	Block: Block{
		Name:      "tablerow",
		Regex:     regexp.MustCompile(`\s*\|(.+?)\|*$`),
		Label:     "<tr>\n%[1]s\n</tr>",
		NeedParse: true,
	},
}

var tablecell = &TableCell{
	Block: Block{
		Name:      "tablecell",
		Regex:     regexp.MustCompile(`\s*\|(.+?)\|*$`),
		Label:     "<td>%[1]s</td>",
		NeedParse: true,
	},
}

// HTML ..
func (s *TableCell) HTML() string {
	inlinetext := &InlineText{
		Text:      s.FirstLine,
		NeedParse: s.NeedParse,
	}
	if s.Label == "" {
		return inlinetext.HTML()
	}
	return fmt.Sprintf(s.Label, inlinetext.HTML())
}

// Open ..
func (s *TableHeader) Open(firstline string) BlockType {
	return &TableHeader{Block: *s.open(firstline)}
}

// Open ..
func (s *TableCell) Open(firstline string) BlockType {
	return &TableCell{Block: *s.open(firstline)}
}

// Open ..
func (s *TableRow) Open(firstline string) BlockType {
	b := &TableRow{
		Block: *s.open(firstline),
	}
	match := b.Regex.FindStringSubmatch(b.FirstLine)
	for _, i := range strings.Split(match[1], "|") {
		b.AddChild(tablecell.Open(i))
	}
	return b
}

// Open ..
func (s *Table) Open(firstline string) BlockType {
	b := &Table{
		Block: *s.open(firstline),
	}
	b.Append(firstline)
	return b
}

// Append ..
func (s *Table) Append(text string) {
	if tablesep.MatchString(text) {
		for _, child := range s.Children {
			for _, cell := range child.(*TableRow).Children {
				cell.SetLabel(tableheader.Label)
			}
		}
		return
	}
	s.AddChild(tablerow.Open(text))
}

// MatchEnd ..
func (s *Table) MatchEnd(text string) bool {
	if s.Match(text) {
		return false
	}
	return true
}
