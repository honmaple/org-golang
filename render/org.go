package render

import (
	"strings"

	"github.com/honmaple/org-golang/parser"
)

type Org struct {
	*parser.Document
}

func (s Org) render(children []parser.Node, sep string) string {
	return concat(s, children, sep)
}

func (s Org) RenderHeading(n *parser.Heading) string {
	var b strings.Builder

	b.WriteString(strings.Repeat("*", n.Stars))
	if n.Keyword != "" {
		b.WriteString(" ")
		b.WriteString(n.Keyword)
	}
	if n.Priority != "" {
		b.WriteString(" [#")
		b.WriteString(n.Priority)
		b.WriteString("]")
	}
	b.WriteString(" ")
	b.WriteString(s.render(n.Title, ""))
	if len(n.Tags) > 0 {
		b.WriteString(" :")
		for _, tag := range n.Tags {
			b.WriteString(tag)
			b.WriteString(":")
		}
	}
	b.WriteString("\n")
	b.WriteString(s.render(n.Children, "\n"))
	return b.String()
}

func (s Org) RenderListItem(n *parser.ListItem) string {
	var b strings.Builder
	if n.Status != "" {
		b.WriteString("[")
		b.WriteString(n.Status)
		b.WriteString("]")
	}
	b.WriteString(s.render(n.Children, "\n"))
	return b.String()
}

func (s Org) RenderList(n *parser.List) string {
	return s.render(n.Children, "\n")
}

func (s Org) RenderTableColumn(n *parser.TableColumn) string {
	return ""
}

func (s Org) RenderTableRow(n *parser.TableRow) string {
	return ""
}

func (s Org) RenderTable(n *parser.Table) string {
	return ""
}

func (s Org) RenderBlock(n *parser.Block) string {
	var b strings.Builder

	b.WriteString("#+begin_")
	b.WriteString(strings.ToLower(n.Type))
	for _, param := range n.Parameters {
		b.WriteString(" ")
		b.WriteString(param)
	}
	b.WriteString("\n")
	if len(n.Children) > 0 {
		if n.Type == "VERSE" {
			b.WriteString(s.render(n.Children, ""))
		} else {
			b.WriteString(s.render(n.Children, "\n"))
		}
		b.WriteString("\n")
	}
	b.WriteString("#+end_")
	b.WriteString(strings.ToLower(n.Type))
	return b.String()
}

func (s Org) RenderBlockResult(n *parser.BlockResult) string {
	return ""
}

func (s Org) RenderDrawer(n *parser.Drawer) string {
	return ""
}

func (s Org) RenderParagraph(n *parser.Paragragh) string {
	return s.render(n.Children, "")
}

func (s Org) RenderBlankline(n *parser.Blankline) string {
	return ""
}

func (s Org) RenderEmphasis(n *parser.InlineEmphasis) string {
	var b strings.Builder

	b.WriteString(n.Marker)
	b.WriteString(s.render(n.Children, ""))
	b.WriteString(n.Marker)
	return b.String()
}

func (s Org) RenderFootnote(*parser.InlineFootnote) string {
	return ""
}

func (s Org) RenderHr(*parser.Hr) string {
	return "-----"
}

func (s Org) RenderKeyword(*parser.Keyword) string {
	return ""
}

func (s Org) RenderLineBreak(*parser.InlineLineBreak) string {
	return "\n"
}

func (s Org) RenderLink(*parser.InlineLink) string {
	return ""
}

func (s Org) RenderPercent(*parser.InlinePercent) string {
	return ""
}

func (s Org) RenderSection(*parser.Section) string {
	return ""
}

func (s Org) RenderText(n *parser.InlineText) string {
	return n.Content
}

func (s Org) RenderTimestamp(*parser.InlineTimestamp) string {
	return ""
}

func (s Org) String() string {
	content := s.render(s.Document.Children, "\n")
	return content
}
