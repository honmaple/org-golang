package render

import (
	"fmt"
	"strings"

	"github.com/honmaple/org-golang/parser"
)

type Renderer interface {
	RenderText(*parser.InlineText) string
	RenderTimestamp(*parser.InlineTimestamp) string
	RenderFootnote(*parser.InlineFootnote) string
	RenderPercent(*parser.InlinePercent) string
	RenderEmphasis(*parser.InlineEmphasis) string
	RenderLineBreak(*parser.InlineLineBreak) string
	RenderLink(*parser.InlineLink) string
	RenderSection(*parser.Section) string
	RenderHeadline(*parser.Headline, int) string
	RenderKeyword(*parser.Keyword, int) string
	RenderBlankline(*parser.Blankline, int) string
	RenderList(*parser.List, int) string
	RenderListItem(*parser.ListItem, int) string
	RenderTable(*parser.Table, int) string
	RenderTableRow(*parser.TableRow, int) string
	RenderTableColumn(*parser.TableColumn, int) string
	RenderBlock(*parser.Block, int) string
	RenderBlockResult(*parser.BlockResult, int) string
	RenderDrawer(*parser.Drawer, int) string
	RenderHr(*parser.Hr, int) string
	RenderParagraph(*parser.Paragragh, int) string
}

func concat(r Renderer, children []parser.Node, sep string, l int) string {
	cs := make([]string, len(children))
	for i, child := range children {
		cs[i] = render(r, child, l)
	}
	return strings.Join(cs, sep)
}

func render(r Renderer, n parser.Node, l int) string {
	switch node := n.(type) {
	case *parser.InlineText:
		return r.RenderText(node)
	case *parser.InlineLineBreak:
		return r.RenderLineBreak(node)
	case *parser.InlineLink:
		return r.RenderLink(node)
	case *parser.InlinePercent:
		return r.RenderPercent(node)
	case *parser.InlineEmphasis:
		return r.RenderEmphasis(node)
	case *parser.Section:
		return r.RenderSection(node)
	case *parser.Headline:
		return r.RenderHeadline(node, l)
	case *parser.Blankline:
		return r.RenderBlankline(node, l)
	case *parser.Keyword:
		return r.RenderKeyword(node, l)
	case *parser.Block:
		return r.RenderBlock(node, l)
	case *parser.BlockResult:
		return r.RenderBlockResult(node, l)
	case *parser.Table:
		return r.RenderTable(node, l)
	case *parser.TableRow:
		return r.RenderTableRow(node, l)
	case *parser.TableColumn:
		return r.RenderTableColumn(node, l)
	case *parser.List:
		return r.RenderList(node, l)
	case *parser.ListItem:
		return r.RenderListItem(node, l)
	case *parser.Drawer:
		return r.RenderDrawer(node, l)
	case *parser.Hr:
		return r.RenderHr(node, l)
	case *parser.Paragragh:
		return r.RenderParagraph(node, l)
	default:
		return ""
	}
}

type Debug struct {
	Document *parser.Document
}

func (s Debug) render(name string, children []parser.Node, sep string, l int) string {
	if len(children) > 0 {
		return fmt.Sprintf("%s%s\n%s", strings.Repeat(" ", l*2), name, concat(s, children, sep, l+1))
	}
	return fmt.Sprintf("%s%s", strings.Repeat(" ", l*2), name)
}

func (s Debug) RenderSection(n *parser.Section) string {
	return n.Name()
}

func (s Debug) RenderHeadline(n *parser.Headline, l int) string {
	return s.render(n.Name(), n.Children, "\n", l)
}

func (s Debug) RenderKeyword(n *parser.Keyword, l int) string {
	return fmt.Sprintf("%s%s", strings.Repeat(" ", l*2), n.Name())
}

func (s Debug) RenderBlock(n *parser.Block, l int) string {
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

func (s Debug) RenderBlockResult(n *parser.BlockResult, l int) string {
	return s.render(n.Name(), n.Children, "\n", l)
}

func (s Debug) RenderDrawer(n *parser.Drawer, l int) string {
	return s.render(n.Name(), n.Children, "\n", l)
}

func (s Debug) RenderListItem(n *parser.ListItem, l int) string {
	return s.render(n.Name(), n.Children, "\n", l)
}

func (s Debug) RenderList(n *parser.List, l int) string {
	return s.render(n.Name(), n.Children, "\n", l)
}

func (s Debug) RenderTableColumn(n *parser.TableColumn, l int) string {
	return s.render(n.Name(), n.Children, "\n", l)
}

func (s Debug) RenderTableRow(n *parser.TableRow, l int) string {
	return fmt.Sprintf("%s%s", strings.Repeat(" ", l*2), n.Name())
}

func (s Debug) RenderTable(n *parser.Table, l int) string {
	return s.render(n.Name(), n.Children, "\n", l)
}

func (s Debug) RenderBlankline(n *parser.Blankline, l int) string {
	return fmt.Sprintf("%s%s[%d]", strings.Repeat(" ", l*2), n.Name(), n.Count)
}

func (s Debug) RenderParagraph(n *parser.Paragragh, l int) string {
	return fmt.Sprintf("%s%s\n%s%s",
		strings.Repeat(" ", l*2), n.Name(),
		strings.Repeat(" ", (l+1)*2), concat(s, n.Children, ",", l+1))
}

func (s Debug) RenderHr(n *parser.Hr, l int) string {
	return fmt.Sprintf("%s%s", strings.Repeat(" ", l*2), n.Name())
}

func (s Debug) RenderText(n *parser.InlineText) string {
	return n.Name()
}

func (s Debug) RenderLineBreak(n *parser.InlineLineBreak) string {
	return n.Name()
}

func (s Debug) RenderFootnote(n *parser.InlineFootnote) string {
	return n.Name()
}

func (s Debug) RenderTimestamp(n *parser.InlineTimestamp) string {
	return n.Name()
}

func (s Debug) RenderPercent(n *parser.InlinePercent) string {
	return n.Name()
}

func (s Debug) RenderLink(n *parser.InlineLink) string {
	return n.Name()
}

func (s Debug) RenderEmphasis(n *parser.InlineEmphasis) string {
	return n.Name()
}

func (s Debug) String() string {
	return concat(s, s.Document.Children, "\n", 0)
}
