package render

import (
	"fmt"
	"strings"

	"github.com/honmaple/org-golang/parser"
)

type Renderer interface {
	InlineText(*parser.InlineText) string
	InlineTimestamp(*parser.InlineTimestamp) string
	InlineFootnote(*parser.InlineFootnote) string
	InlinePercent(*parser.InlinePercent) string
	InlineEmphasis(*parser.InlineEmphasis) string
	InlineLineBreak(*parser.InlineLineBreak) string
	InlineLink(*parser.InlineLink) string
	Headline(*parser.Headline, int) string
	Keyword(*parser.Keyword, int) string
	Blankline(*parser.Blankline, int) string
	List(*parser.List, int) string
	ListItem(*parser.ListItem, int) string
	Table(*parser.Table, int) string
	TableRow(*parser.TableRow, int) string
	TableColumn(*parser.TableColumn, int) string
	Block(*parser.Block, int) string
	BlockResult(*parser.BlockResult, int) string
	Drawer(*parser.Drawer, int) string
	Hr(*parser.Hr, int) string
	Paragraph(*parser.Paragragh, int) string
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
		return r.InlineText(node)
	case *parser.InlineLineBreak:
		return r.InlineLineBreak(node)
	case *parser.InlineLink:
		return r.InlineLink(node)
	case *parser.InlinePercent:
		return r.InlinePercent(node)
	case *parser.InlineEmphasis:
		return r.InlineEmphasis(node)
	case *parser.Headline:
		return r.Headline(node, l)
	case *parser.Blankline:
		return r.Blankline(node, l)
	case *parser.Keyword:
		return r.Keyword(node, l)
	case *parser.Block:
		return r.Block(node, l)
	case *parser.BlockResult:
		return r.BlockResult(node, l)
	case *parser.Table:
		return r.Table(node, l)
	case *parser.TableRow:
		return r.TableRow(node, l)
	case *parser.TableColumn:
		return r.TableColumn(node, l)
	case *parser.List:
		return r.List(node, l)
	case *parser.ListItem:
		return r.ListItem(node, l)
	case *parser.Drawer:
		return r.Drawer(node, l)
	case *parser.Hr:
		return r.Hr(node, l)
	case *parser.Paragragh:
		return r.Paragraph(node, l)
	default:
		return ""
	}
	// switch n.Name() {
	// case parser.InlineTextName:
	//	return r.InlineText(n.(*parser.InlineText))
	// case parser.InlineLineBreakName:
	//	return r.InlineLineBreak(n.(*parser.InlineLineBreak))
	// case parser.HrName:
	//	return r.Hr(n.(*parser.Hr))
	// case parser.InlineLinkName:
	//	return r.InlineLink(n.(*parser.InlineLink))
	// case parser.InlinePercentName:
	//	return r.InlinePercent(n.(*parser.InlinePercent))
	// case parser.InlineEmphasisName:
	//	return r.InlineEmphasis(n.(*parser.InlineEmphasis))
	// case parser.HeadlineName:
	//	return r.Headline(n.(*parser.Headline))
	// case parser.KeywordName:
	//	return r.Keyword(n.(*parser.Keyword))
	// case parser.BlockName:
	//	return r.Block(n.(*parser.Block))
	// case parser.BlockResultName:
	//	return r.BlockResult(n.(*parser.BlockResult))
	// case parser.TableName:
	//	return r.Table(n.(*parser.Table))
	// case parser.ListName:
	//	return r.List(n.(*parser.List))
	// case parser.ListItemName:
	//	return r.ListItem(n.(*parser.ListItem))
	// case parser.DrawerName:
	//	return r.Drawer(n.(*parser.Drawer))
	// case parser.ParagraghName:
	//	return r.Paragraph(n.(*parser.Paragragh))
	// default:
	//	return ""
	// }
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

func (s Debug) String() string {
	return concat(s, s.Document.Children, "\n", 0)
}

func (s Debug) Headline(n *parser.Headline, l int) string {
	return s.render(n.Name(), n.Children, "\n", l)
}

func (s Debug) Keyword(n *parser.Keyword, l int) string {
	return fmt.Sprintf("%s%s", strings.Repeat(" ", l*2), n.Name())
}

func (s Debug) Block(n *parser.Block, l int) string {
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

func (s Debug) BlockResult(n *parser.BlockResult, l int) string {
	return s.render(n.Name(), n.Children, "\n", l)
}

func (s Debug) Drawer(n *parser.Drawer, l int) string {
	return s.render(n.Name(), n.Children, "\n", l)
}

func (s Debug) ListItem(n *parser.ListItem, l int) string {
	return s.render(n.Name(), n.Children, "\n", l)
}

func (s Debug) List(n *parser.List, l int) string {
	return s.render(n.Name(), n.Children, "\n", l)
}

func (s Debug) TableColumn(n *parser.TableColumn, l int) string {
	return s.render(n.Name(), n.Children, "\n", l)
}

func (s Debug) TableRow(n *parser.TableRow, l int) string {
	return fmt.Sprintf("%s%s", strings.Repeat(" ", l*2), n.Name())
}

func (s Debug) Table(n *parser.Table, l int) string {
	return s.render(n.Name(), n.Children, "\n", l)
}

func (s Debug) Blankline(n *parser.Blankline, l int) string {
	return fmt.Sprintf("%s%s[%d]", strings.Repeat(" ", l*2), n.Name(), n.Count)
}

func (s Debug) Paragraph(n *parser.Paragragh, l int) string {
	return fmt.Sprintf("%s%s\n%s%s",
		strings.Repeat(" ", l*2), n.Name(),
		strings.Repeat(" ", (l+1)*2), concat(s, n.Children, ",", l+1))
}

func (s Debug) Hr(n *parser.Hr, l int) string {
	return fmt.Sprintf("%s%s", strings.Repeat(" ", l*2), n.Name())
}

func (s Debug) InlineText(n *parser.InlineText) string {
	return n.Name()
}

func (s Debug) InlineLineBreak(n *parser.InlineLineBreak) string {
	return n.Name()
}

func (s Debug) InlineFootnote(n *parser.InlineFootnote) string {
	return n.Name()
}

func (s Debug) InlineTimestamp(n *parser.InlineTimestamp) string {
	return n.Name()
}

func (s Debug) InlinePercent(n *parser.InlinePercent) string {
	return n.Name()
}

func (s Debug) InlineLink(n *parser.InlineLink) string {
	return n.Name()
}

func (s Debug) InlineEmphasis(n *parser.InlineEmphasis) string {
	return n.Name()
}
