package render

import (
	"fmt"
	"strings"

	"github.com/honmaple/org-golang/parser"
)

type DebugRender struct {
	Document *parser.Document
}

func (s DebugRender) String() string {
	return concat(s, s.Document.Children, "\n", 0)
}

func (s DebugRender) Headline(n *parser.Headline, l int) string {
	return s.render(n.Name(), n.Children, "\n", l)
}

func (s DebugRender) Keyword(n *parser.Keyword, l int) string {
	return fmt.Sprintf("%s%s", strings.Repeat(" ", l*2), n.Name())
}

func (s DebugRender) Block(n *parser.Block, l int) string {
	name := fmt.Sprintf("%s[%s]", n.Name(), n.Type)
	switch n.Type {
	case "SRC", "EXAMPLE", "VERSE":
		if len(n.Children) > 0 {
			return fmt.Sprintf("%s%s\n%s%s",
				strings.Repeat(" ", l*2), name,
				strings.Repeat(" ", (l+1)*2),
				concat(s, n.Children, ",", l+1))
		}
		return fmt.Sprintf("%s%s", strings.Repeat(" ", l*2), name)
	default:
		return s.render(name, n.Children, "\n", l)
	}
}

func (s DebugRender) BlockResult(n *parser.BlockResult, l int) string {
	return s.render(n.Name(), n.Children, "\n", l)
}

func (s DebugRender) Drawer(n *parser.Drawer, l int) string {
	return s.render(n.Name(), n.Children, "\n", l)
}

func (s DebugRender) ListItem(n *parser.ListItem, l int) string {
	return s.render(n.Name(), n.Children, "\n", l)
}

func (s DebugRender) List(n *parser.List, l int) string {
	return s.render(n.Name(), n.Children, "\n", l)
}

func (s DebugRender) TableColumn(n *parser.TableColumn, l int) string {
	return s.render(n.Name(), n.Children, "\n", l)
}

func (s DebugRender) TableRow(n *parser.TableRow, l int) string {
	return fmt.Sprintf("%s%s", strings.Repeat(" ", l*2), n.Name())
}

func (s DebugRender) Table(n *parser.Table, l int) string {
	var b strings.Builder
	b.WriteString(strings.Repeat(" ", l*2))
	b.WriteString(n.Name())
	b.WriteString("\n")
	for i, child := range n.Children {
		if i > 0 && i < len(n.Children) {
			b.WriteString("\n")
		}
		b.WriteString(s.TableRow(child, l+1))
	}
	return b.String()
}

func (s DebugRender) Blankline(n *parser.Blankline, l int) string {
	return fmt.Sprintf("%s%s[%d]", strings.Repeat(" ", l*2), n.Name(), n.Count)
}

func (s DebugRender) Paragraph(n *parser.Paragragh, l int) string {
	return fmt.Sprintf("%s%s\n%s%s",
		strings.Repeat(" ", l*2), n.Name(),
		strings.Repeat(" ", (l+1)*2), concat(s, n.Children, ",", l+1))
}

func (s DebugRender) Hr(n *parser.Hr, l int) string {
	return fmt.Sprintf("%s%s", strings.Repeat(" ", l*2), n.Name())
}

func (s DebugRender) InlineText(n *parser.InlineText) string {
	return n.Name()
}

func (s DebugRender) InlineLineBreak(n *parser.InlineLineBreak) string {
	return n.Name()
}

func (s DebugRender) InlineFootnote(n *parser.InlineFootnote) string {
	return n.Name()
}

func (s DebugRender) InlineTimestamp(n *parser.InlineTimestamp) string {
	return n.Name()
}

func (s DebugRender) InlinePercent(n *parser.InlinePercent) string {
	return n.Name()
}

func (s DebugRender) InlineLink(n *parser.InlineLink) string {
	return n.Name()
}

func (s DebugRender) InlineEmphasis(n *parser.InlineEmphasis) string {
	return n.Name()
}

func (s DebugRender) render(name string, children []parser.Node, sep string, l int) string {
	if len(children) > 0 {
		return fmt.Sprintf("%s%s\n%s", strings.Repeat(" ", l*2), name, concat(s, children, sep, l+1))
	}
	return fmt.Sprintf("%s%s", strings.Repeat(" ", l*2), name)
}
