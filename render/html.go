package render

import (
	"fmt"
	"strings"

	"github.com/honmaple/org-golang/parser"
)

type Render interface {
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
	Block(*parser.Block, int) string
	BlockResult(*parser.BlockResult, int) string
	Drawer(*parser.Drawer, int) string
	Hr(*parser.Hr, int) string
	Paragraph(*parser.Paragragh, int) string
}

type HTMLRender struct {
	Document  *parser.Document
	Toc       bool
	Offset    int
	Highlight func(text string, lang string) string
}

const (
	headlineElement         = "<h%[1]d>%[2]s</h%[1]d>"
	headlineKeywordElement  = "<span class=\"todo\">%[1]s</span>"
	headlinePriorityElement = "<span class=\"priority\">%[1]s</span>"
	headlineTagElement      = "<span class=\"tag\">%[1]s</span>"
)

const (
	centerElement = "<p style=\"text-align:center;\">%[1]s</p>"
	exportElement = "%[1]s"
	quoteElement  = "<blockquote>\n%[1]s\n</blockquote>"
	verseElement  = "<p>\n%[1]s\n</p>"
	srcElement    = "<pre class=\"src src-%[1]s\">\n%[2]s\n</pre>"
)

const (
	listitemElement         = "<li>\n%[1]s</li>"
	listitemstatusElement   = "<code>%[1]s</code>"
	orderlistElement        = "<ol>\n%[1]s\n</ol>"
	unorderlistElement      = "<ul>\n%[1]s\n</ul>"
	descriptiveElement      = "<dl>\n%[1]s\n</dl>"
	descriptiveDescElement  = "<dd>%[1]s</dd>"
	descriptiveTitleElement = "<dt>%[1]s</dt>"
)

const (
	tableElement       = "<table>\n%[1]s\n</table>"
	tableRowElement    = "<tr>\n%[1]s\n</tr>"
	tableHeaderElement = "<th class=\"align-%[1]s\">%[2]s</th>"
	tableColumnElement = "<td class=\"align-%[1]s\">%[2]s</td>"
)

const (
	paragraghElement = "<p>\n%[1]s\n</p>"
)

func (s HTMLRender) InlineText(n *parser.InlineText) string {
	return n.Content
}

func (s HTMLRender) InlineLineBreak(n *parser.InlineLineBreak) string {
	return strings.Repeat("\n", n.Count)
}

func (s HTMLRender) InlineFootnote(n *parser.InlineFootnote) string {
	return ""
}

func (s HTMLRender) InlineTimestamp(n *parser.InlineTimestamp) string {
	return ""
}

func (s HTMLRender) InlinePercent(n *parser.InlinePercent) string {
	return fmt.Sprintf("<code>[%s]</code>", n.Num)
}

func (s HTMLRender) InlineLink(n *parser.InlineLink) string {
	if n.IsImage() {
		return fmt.Sprintf("<img src=\"%s\"/>", n.URL)
	}
	if n.IsVideo() {
		return fmt.Sprintf("<video src=\"%s\">%s</video>", n.URL, n.URL)
	}
	if n.Desc == "" {
		return fmt.Sprintf("<a href=\"%s\">%s</a>", n.URL, n.URL)
	}
	return fmt.Sprintf("<a href=\"%s\">%s</a>", n.URL, n.Desc)
}

func (s HTMLRender) InlineEmphasis(n *parser.InlineEmphasis) string {
	text := s.render(n.Children, 0, "")
	switch n.Marker {
	case "=", "~", "`":
		return fmt.Sprintf("<code>%s</code>", text)
	case "*":
		return fmt.Sprintf("<bold>%s</bold>", text)
	case "_":
		return fmt.Sprintf("<span style=\"text-decoration:underline\">%s</span>", text)
	case "+":
		return fmt.Sprintf("<del>%s</del>", text)
	case "/":
		return fmt.Sprintf("<i>%s</i>", text)
	}
	return ""
}

func (s HTMLRender) Headline(n *parser.Headline, l int) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("<h%[1]d>", n.Stars))
	if n.Keyword != "" {
		b.WriteString(fmt.Sprintf(headlineKeywordElement, n.Keyword))
	}
	if n.Priority != "" {
		b.WriteString(fmt.Sprintf(headlinePriorityElement, n.Priority))
	}
	b.WriteString(s.render(n.Title, l, ""))
	for _, tag := range n.Tags {
		b.WriteString(fmt.Sprintf(headlineTagElement, tag))
	}
	b.WriteString(fmt.Sprintf("</h%[1]d>", n.Stars))
	if len(n.Children) > 0 {
		b.WriteString("\n")
	}
	b.WriteString(s.render(n.Children, l, "\n"))
	return b.String()
}

func (s HTMLRender) Keyword(n *parser.Keyword, l int) string {
	return ""
}

func (s HTMLRender) ListItem(n *parser.ListItem, l int) string {
	content := s.render(n.Children, l, "\n")
	if n.Status != "" {
		content = fmt.Sprintf("%[1]s %[2]s",
			fmt.Sprintf(listitemstatusElement, n.Status),
			content,
		)
	}
	return fmt.Sprintf(listitemElement, content)
}

func (s HTMLRender) List(n *parser.List, l int) string {
	content := s.render(n.Children, l, "\n")
	switch n.Type {
	case parser.OrderlistName:
		return fmt.Sprintf(orderlistElement, content)
	case parser.UnorderlistName:
		return fmt.Sprintf(unorderlistElement, content)
	case parser.DescriptiveName:
		return fmt.Sprintf(descriptiveElement, content)
	default:
		return ""
	}
}

func (s HTMLRender) TableColumn(n *parser.TableColumn, l int) string {
	if n.IsHeader {
		return fmt.Sprintf(tableHeaderElement, "", s.render(n.Children, l, ""))
	}
	return fmt.Sprintf(tableColumnElement, "", s.render(n.Children, l, ""))
}

func (s HTMLRender) TableRow(n *parser.TableRow, l int) string {
	if n.Separator {
		return ""
	}

	var b strings.Builder
	for _, child := range n.Children {
		b.WriteString(s.TableColumn(child, l))
	}
	return fmt.Sprintf(tableRowElement, b.String())
}

func (s HTMLRender) Table(n *parser.Table, l int) string {
	var b strings.Builder
	for _, child := range n.Children {
		b.WriteString(s.TableRow(child, l))
	}
	return fmt.Sprintf(tableElement, b.String())
}

func (s HTMLRender) Block(n *parser.Block, l int) string {
	switch n.Type {
	case "SRC":
		if s.Highlight == nil {
			return fmt.Sprintf(srcElement, n.Parameters[0], s.render(n.Children, l, "\n"))
		}
		return s.Highlight(n.Parameters[0], s.render(n.Children, l, "\n"))
	case "EXAMPLE":
		return fmt.Sprintf(srcElement, "example", s.render(n.Children, l, "\n"))
	case "CENTER":
		return fmt.Sprintf(centerElement, s.render(n.Children, l, "\n"))
	case "QUOTE":
		return fmt.Sprintf(quoteElement, s.render(n.Children, l, "\n"))
	case "EXPORT":
		return fmt.Sprintf(exportElement, s.render(n.Children, l, "\n"))
	case "VERSE":
		var b strings.Builder
		for _, child := range n.Children {
			if child.Name() == parser.InlineLineBreakName {
				b.WriteString("<br />\n")
				continue
			}
			b.WriteString(render(s, child, l))
		}
		return fmt.Sprintf(verseElement, b.String())
	}
	return s.render(n.Children, l, "\n")
}

func (s HTMLRender) BlockResult(n *parser.BlockResult, l int) string {
	return s.render(n.Children, l, "\n")
}

func (s HTMLRender) Drawer(n *parser.Drawer, l int) string {
	return s.render(n.Children, l, "\n")
}

func (s HTMLRender) Blankline(n *parser.Blankline, l int) string {
	return ""
}

func (s HTMLRender) Hr(n *parser.Hr, l int) string {
	return "<hr/>"
}

func (s HTMLRender) Paragraph(n *parser.Paragragh, l int) string {
	return fmt.Sprintf(paragraghElement, s.render(n.Children, l, ""))
}

func (s HTMLRender) Section(n *parser.Section) string {
	return ""
}

func (s HTMLRender) String() string {
	content := s.render(s.Document.Children, 0, "\n")
	if !s.Toc || s.Document.Get("toc") == "nil" {
		return content
	}
	return content
}

func (s HTMLRender) render(children []parser.Node, l int, sep string) string {
	return concat(s, children, sep, l)
}

func concat(r Render, children []parser.Node, sep string, l int) string {
	cs := make([]string, len(children))
	for i, child := range children {
		cs[i] = render(r, child, l)
	}
	return strings.Join(cs, sep)
}

func render(r Render, n parser.Node, l int) string {
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
