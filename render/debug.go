package render

import (
	"fmt"
	"strings"

	"github.com/honmaple/org-golang/parser"
)

type Renderer interface {
	RenderNode(parser.Node, bool) string
	RenderNodes([]parser.Node, string) string
	RenderInlineLink(*parser.InlineLink) string
	RenderInlineText(*parser.InlineText) string
	RenderInlinePercent(*parser.InlinePercent) string
	RenderInlineEmphasis(*parser.InlineEmphasis) string
	RenderInlineTimestamp(*parser.InlineTimestamp) string
	RenderInlineLineBreak(*parser.InlineLineBreak) string
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
	RenderFootnote(*parser.Footnote) string
	RenderParagraph(*parser.Paragragh) string
}

func DedentString(text string) string {
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

func RenderNodes(r Renderer, children []parser.Node, sep string) string {
	cs := make([]string, len(children))
	for i, child := range children {
		cs[i] = r.RenderNode(child, false)
	}
	return strings.Join(cs, sep)
}

func RenderNode(r Renderer, n parser.Node) string {
	switch node := n.(type) {
	case *parser.InlineText:
		return r.RenderInlineText(node)
	case *parser.InlineLineBreak:
		return r.RenderInlineLineBreak(node)
	case *parser.InlineLink:
		return r.RenderInlineLink(node)
	case *parser.InlinePercent:
		return r.RenderInlinePercent(node)
	case *parser.InlineEmphasis:
		return r.RenderInlineEmphasis(node)
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
	case *parser.Footnote:
		return r.RenderFootnote(node)
	case *parser.Paragragh:
		return r.RenderParagraph(node)
	default:
		return ""
	}
}

type Debug struct {
	Document *parser.Document
}

func (r *Debug) render(name string, children []parser.Node, sep string) string {
	text := name
	if len(children) > 0 {
		lines := strings.Split(r.RenderNodes(children, sep), "\n")
		for i, line := range lines {
			lines[i] = "  " + line
		}
		text = text + "\n" + strings.Join(lines, "\n")
	}
	return text
}

func (r *Debug) RenderNode(n parser.Node, def bool) string {
	return RenderNode(r, n)
}

func (r *Debug) RenderNodes(children []parser.Node, sep string) string {
	return RenderNodes(r, children, sep)
}

func (r *Debug) RenderInlineText(n *parser.InlineText) string {
	return n.Name()
}

func (r *Debug) RenderInlineLineBreak(n *parser.InlineLineBreak) string {
	return n.Name()
}

func (r *Debug) RenderInlineTimestamp(n *parser.InlineTimestamp) string {
	return n.Name()
}

func (r *Debug) RenderInlinePercent(n *parser.InlinePercent) string {
	return n.Name()
}

func (r *Debug) RenderInlineLink(n *parser.InlineLink) string {
	return n.Name()
}

func (r *Debug) RenderInlineEmphasis(n *parser.InlineEmphasis) string {
	return n.Name()
}

func (r *Debug) RenderFootnote(n *parser.Footnote) string {
	return n.Name()
}

func (r *Debug) RenderSection(n *parser.Section) string {
	return n.Name()
}

func (r *Debug) RenderHeading(n *parser.Heading) string {
	return r.render(n.Name(), n.Children, "\n")
}

func (r *Debug) RenderKeyword(n *parser.Keyword) string {
	return n.Name()
}

func (r *Debug) RenderBlock(n *parser.Block) string {
	name := fmt.Sprintf("%s[%s]", n.Name(), n.Type)
	return r.render(name, n.Children, "\n")
}

func (r *Debug) RenderBlockResult(n *parser.BlockResult) string {
	return r.render(n.Name(), n.Children, "\n")
}

func (r *Debug) RenderDrawer(n *parser.Drawer) string {
	return r.render(n.Name(), n.Children, "\n")
}

func (r *Debug) RenderListItem(n *parser.ListItem) string {
	return r.render(n.Name(), n.Children, "\n")
}

func (r *Debug) RenderList(n *parser.List) string {
	return r.render(n.Name(), n.Children, "\n")
}

func (r *Debug) RenderTableColumn(n *parser.TableColumn) string {
	return r.render(n.Name(), n.Children, "\n")
}

func (r *Debug) RenderTableRow(n *parser.TableRow) string {
	return n.Name()
}

func (r *Debug) RenderTable(n *parser.Table) string {
	return r.render(n.Name(), n.Children, "\n")
}

func (r *Debug) RenderBlankline(n *parser.Blankline) string {
	return fmt.Sprintf("%s[%d]", n.Name(), n.Count)
}

func (r *Debug) RenderParagraph(n *parser.Paragragh) string {
	return r.render(n.Name(), n.Children, ",")
}

func (r *Debug) RenderHr(n *parser.Hr) string {
	return n.Name()
}

func (r *Debug) String() string {
	return r.RenderNodes(r.Document.Children, "\n")
}
