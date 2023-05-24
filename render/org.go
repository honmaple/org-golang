package render

import (
	"strings"

	"github.com/honmaple/org-golang/parser"
)

type Org struct {
	Document       *parser.Document
	RenderNodeFunc func(r Renderer, n parser.Node) string
}

func (r *Org) RenderNode(n parser.Node, def bool) string {
	if def || r.RenderNodeFunc == nil {
		return RenderNode(r, n)
	}
	return r.RenderNodeFunc(r, n)
}

func (r *Org) RenderNodes(children []parser.Node, sep string) string {
	return RenderNodes(r, children, sep)
}

func (r *Org) RenderInlineLink(*parser.InlineLink) string {
	return ""
}

func (r *Org) RenderInlineText(n *parser.InlineText) string {
	return n.Content
}

func (r *Org) RenderInlinePercent(*parser.InlinePercent) string {
	return ""
}

func (r *Org) RenderInlineEmphasis(n *parser.InlineEmphasis) string {
	var b strings.Builder

	b.WriteString(n.Marker)
	b.WriteString(r.RenderNodes(n.Children, ""))
	b.WriteString(n.Marker)
	return b.String()
}

func (r *Org) RenderInlineTimestamp(*parser.InlineTimestamp) string {
	return ""
}

func (r *Org) RenderInlineLineBreak(n *parser.InlineLineBreak) string {
	return strings.Repeat("\\", n.Count)
}

func (r *Org) RenderInlineBackSlash(n *parser.InlineBackSlash) string {
	return strings.Repeat("\\", n.Count)
}

func (r *Org) RenderInlineFootnote(*parser.Footnote) string {
	return ""
}

func (r *Org) RenderFootnote(*parser.Footnote) string {
	return ""
}

func (r *Org) RenderHeading(n *parser.Heading) string {
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
	b.WriteString(r.RenderNodes(n.Title, ""))
	if len(n.Tags) > 0 {
		b.WriteString(" :")
		for _, tag := range n.Tags {
			b.WriteString(tag)
			b.WriteString(":")
		}
	}
	b.WriteString("\n")
	b.WriteString(r.RenderNodes(n.Children, "\n"))
	return b.String()
}

func (r *Org) RenderListItem(n *parser.ListItem) string {
	var b strings.Builder
	if n.Status != "" {
		b.WriteString("[")
		b.WriteString(n.Status)
		b.WriteString("]")
	}
	b.WriteString(r.RenderNodes(n.Children, "\n"))
	return b.String()
}

func (r *Org) RenderList(n *parser.List) string {
	return r.RenderNodes(n.Children, "\n")
}

func (r *Org) RenderTableColumn(n *parser.TableColumn) string {
	return ""
}

func (r *Org) RenderTableRow(n *parser.TableRow) string {
	return ""
}

func (r *Org) RenderTable(n *parser.Table) string {
	return ""
}

func (r *Org) RenderBlock(n *parser.Block) string {
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
			b.WriteString(r.RenderNodes(n.Children, ""))
		} else {
			b.WriteString(r.RenderNodes(n.Children, "\n"))
		}
		b.WriteString("\n")
	}
	b.WriteString("#+end_")
	b.WriteString(strings.ToLower(n.Type))
	return b.String()
}

func (r *Org) RenderBlockResult(n *parser.BlockResult) string {
	return ""
}

func (r *Org) RenderDrawer(n *parser.Drawer) string {
	return ""
}

func (r *Org) RenderParagraph(n *parser.Paragragh) string {
	return r.RenderNodes(n.Children, "")
}

func (r *Org) RenderBlankline(n *parser.Blankline) string {
	return ""
}

func (r *Org) RenderHr(*parser.Hr) string {
	return "-----"
}

func (r *Org) RenderKeyword(*parser.Keyword) string {
	return ""
}

func (r *Org) RenderSection(*parser.Section) string {
	return ""
}

func (r *Org) String() string {
	return r.RenderNodes(r.Document.Children, "\n")
}
