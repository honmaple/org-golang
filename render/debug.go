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
	RenderHeading(*parser.Heading) string
	RenderKeyword(*parser.Keyword) string
	RenderBlankline(*parser.Blankline) string
	RenderList(*parser.List) string
	RenderListItem(*parser.ListItem) string
	RenderTable(*parser.Table) string
	RenderTableRow(*parser.TableRow) string
	RenderTableColumn(*parser.TableColumn) string
	RenderBlock(*parser.Block) string
	RenderBlockResult(*parser.BlockResult) string
	RenderDrawer(*parser.Drawer) string
	RenderHr(*parser.Hr) string
	RenderParagraph(*parser.Paragragh) string
}

func dedent(text string) string {
	min := -1

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		newline := strings.TrimLeft(line, " ")
		if newline == "" {
			continue
		}
		if indent := len(line) - len(newline); min == -1 || indent < min {
			min = indent
		}
	}
	if min == -1 {
		return text
	}
	for i, line := range lines {
		lines[i] = strings.TrimPrefix(line, strings.Repeat(" ", min))
	}
	return strings.Join(lines, "\n")
}

func concat(r Renderer, children []parser.Node, sep string) string {
	cs := make([]string, len(children))
	for i, child := range children {
		cs[i] = render(r, child)
	}
	return strings.Join(cs, sep)
}

func render(r Renderer, n parser.Node) string {
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
	case *parser.Heading:
		return r.RenderHeading(node)
	case *parser.Blankline:
		return r.RenderBlankline(node)
	case *parser.Keyword:
		return r.RenderKeyword(node)
	case *parser.Block:
		return r.RenderBlock(node)
	case *parser.BlockResult:
		return r.RenderBlockResult(node)
	case *parser.Table:
		return r.RenderTable(node)
	case *parser.TableRow:
		return r.RenderTableRow(node)
	case *parser.TableColumn:
		return r.RenderTableColumn(node)
	case *parser.List:
		return r.RenderList(node)
	case *parser.ListItem:
		return r.RenderListItem(node)
	case *parser.Drawer:
		return r.RenderDrawer(node)
	case *parser.Hr:
		return r.RenderHr(node)
	case *parser.Paragragh:
		return r.RenderParagraph(node)
	default:
		return ""
	}
}

type Debug struct {
	Document *parser.Document
}

func (s Debug) render(name string, children []parser.Node, sep string) string {
	text := name
	if len(children) > 0 {
		lines := strings.Split(concat(s, children, sep), "\n")
		for i, line := range lines {
			lines[i] = "  " + line
		}
		text = text + "\n" + strings.Join(lines, "\n")
	}
	return text
}

func (s Debug) RenderSection(n *parser.Section) string {
	return n.Name()
}

func (s Debug) RenderHeading(n *parser.Heading) string {
	return s.render(n.Name(), n.Children, "\n")
}

func (s Debug) RenderKeyword(n *parser.Keyword) string {
	return n.Name()
}

func (s Debug) RenderBlock(n *parser.Block) string {
	name := fmt.Sprintf("%s[%s]", n.Name(), n.Type)
	return s.render(name, n.Children, "\n")
}

func (s Debug) RenderBlockResult(n *parser.BlockResult) string {
	return s.render(n.Name(), n.Children, "\n")
}

func (s Debug) RenderDrawer(n *parser.Drawer) string {
	return s.render(n.Name(), n.Children, "\n")
}

func (s Debug) RenderListItem(n *parser.ListItem) string {
	return s.render(n.Name(), n.Children, "\n")
}

func (s Debug) RenderList(n *parser.List) string {
	return s.render(n.Name(), n.Children, "\n")
}

func (s Debug) RenderTableColumn(n *parser.TableColumn) string {
	return s.render(n.Name(), n.Children, "\n")
}

func (s Debug) RenderTableRow(n *parser.TableRow) string {
	return n.Name()
}

func (s Debug) RenderTable(n *parser.Table) string {
	return s.render(n.Name(), n.Children, "\n")
}

func (s Debug) RenderBlankline(n *parser.Blankline) string {
	return fmt.Sprintf("%s[%d]", n.Name(), n.Count)
}

func (s Debug) RenderParagraph(n *parser.Paragragh) string {
	return s.render(n.Name(), n.Children, ",")
}

func (s Debug) RenderHr(n *parser.Hr) string {
	return n.Name()
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
	return concat(s, s.Document.Children, "\n")
}
