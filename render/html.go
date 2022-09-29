package render

import (
	"fmt"
	"strings"

	"github.com/honmaple/org-golang/parser"
)

type HTML struct {
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

func (s HTML) render(children []parser.Node, l int, sep string) string {
	return concat(s, children, sep, l)
}

func (s HTML) InlineText(n *parser.InlineText) string {
	return n.Content
}

func (s HTML) InlineLineBreak(n *parser.InlineLineBreak) string {
	return strings.Repeat("\n", n.Count)
}

func (s HTML) InlineFootnote(n *parser.InlineFootnote) string {
	return ""
}

func (s HTML) InlineTimestamp(n *parser.InlineTimestamp) string {
	return ""
}

func (s HTML) InlinePercent(n *parser.InlinePercent) string {
	return fmt.Sprintf("<code>[%s]</code>", n.Num)
}

func (s HTML) InlineLink(n *parser.InlineLink) string {
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

func (s HTML) InlineEmphasis(n *parser.InlineEmphasis) string {
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

func (s HTML) headline(n *parser.Headline, l int) string {
	var b strings.Builder

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
	return b.String()
}

func (s HTML) Headline(n *parser.Headline, l int) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("<h%[1]d>", n.Stars+s.Offset))
	b.WriteString(s.headline(n, l))
	b.WriteString(fmt.Sprintf("</h%[1]d>", n.Stars))
	if len(n.Children) > 0 {
		b.WriteString("\n")
	}
	b.WriteString(s.render(n.Children, l, "\n"))
	return b.String()
}

func (s HTML) Keyword(n *parser.Keyword, l int) string {
	return ""
}

func (s HTML) ListItem(n *parser.ListItem, l int) string {
	content := s.render(n.Children, l, "\n")
	if n.Status != "" {
		content = fmt.Sprintf("%[1]s %[2]s",
			fmt.Sprintf(listitemstatusElement, n.Status),
			content,
		)
	}
	return fmt.Sprintf(listitemElement, content)
}

func (s HTML) List(n *parser.List, l int) string {
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

func (s HTML) TableColumn(n *parser.TableColumn, l int) string {
	if n.IsHeader {
		return fmt.Sprintf(tableHeaderElement, "", s.render(n.Children, l, ""))
	}
	return fmt.Sprintf(tableColumnElement, "", s.render(n.Children, l, ""))
}

func (s HTML) TableRow(n *parser.TableRow, l int) string {
	if n.Separator {
		return ""
	}
	return fmt.Sprintf(tableRowElement, s.render(n.Children, l, "\n"))
}

func (s HTML) Table(n *parser.Table, l int) string {
	return fmt.Sprintf(tableElement, s.render(n.Children, l, "\n"))
}

func (s HTML) Block(n *parser.Block, l int) string {
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

func (s HTML) BlockResult(n *parser.BlockResult, l int) string {
	return s.render(n.Children, l, "\n")
}

func (s HTML) Drawer(n *parser.Drawer, l int) string {
	return s.render(n.Children, l, "\n")
}

func (s HTML) Blankline(n *parser.Blankline, l int) string {
	return ""
}

func (s HTML) Hr(n *parser.Hr, l int) string {
	return "<hr/>"
}

func (s HTML) Paragraph(n *parser.Paragragh, l int) string {
	return fmt.Sprintf(paragraghElement, s.render(n.Children, l, ""))
}

func (s HTML) Section(sections []*parser.Section) string {
	var b strings.Builder

	b.WriteString("<ul>\n")
	for _, section := range sections {
		b.WriteString(fmt.Sprintf(`<li><a href="%s">%s</a>`, section.Id(), s.headline(section.Headline, 0)))
		if len(section.Children) > 0 {
			b.WriteString("\n")
			b.WriteString(s.Section(section.Children))
		}
		b.WriteString("</li>\n")
	}
	b.WriteString("</ul>")
	return b.String()
}

func (s HTML) String() string {
	content := s.render(s.Document.Children, 0, "\n")
	if !s.Toc || s.Document.Get("toc") == "nil" || len(s.Document.Sections.Children) == 0 {
		return content
	}
	return s.Section(s.Document.Sections.Children) + content
}
