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
	headingElement         = "<h%[1]d>%[2]s</h%[1]d>"
	headingKeywordElement  = "<span class=\"todo\">%[1]s</span>"
	headingPriorityElement = "<span class=\"priority\">%[1]s</span>"
	headingTagElement      = "<span class=\"tag\">%[1]s</span>"
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

func (s HTML) render(children []parser.Node, sep string) string {
	return concat(s, children, sep, 0)
}

func (s HTML) RenderText(n *parser.InlineText) string {
	return n.Content
}

func (s HTML) RenderLineBreak(n *parser.InlineLineBreak) string {
	return strings.Repeat("\n", n.Count)
}

func (s HTML) RenderFootnote(n *parser.InlineFootnote) string {
	return ""
}

func (s HTML) RenderTimestamp(n *parser.InlineTimestamp) string {
	return ""
}

func (s HTML) RenderPercent(n *parser.InlinePercent) string {
	return fmt.Sprintf("<code>[%s]</code>", n.Num)
}

func (s HTML) RenderLink(n *parser.InlineLink) string {
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

func (s HTML) RenderEmphasis(n *parser.InlineEmphasis) string {
	text := s.render(n.Children, "")
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
	default:
		return fmt.Sprintf("%[1]s%[2]s%[1]s", n.Marker, text)
	}
}

func (s HTML) heading(n *parser.Heading) string {
	var b strings.Builder

	if n.Keyword != "" {
		b.WriteString(fmt.Sprintf(headingKeywordElement, n.Keyword))
	}
	if n.Priority != "" {
		b.WriteString(fmt.Sprintf(headingPriorityElement, n.Priority))
	}
	b.WriteString(s.render(n.Title, ""))
	for _, tag := range n.Tags {
		b.WriteString(fmt.Sprintf(headingTagElement, tag))
	}
	return b.String()
}

func (s HTML) RenderHeading(n *parser.Heading) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("<h%[1]d>", n.Stars+s.Offset))
	b.WriteString(s.heading(n))
	b.WriteString(fmt.Sprintf("</h%[1]d>", n.Stars))
	if len(n.Children) > 0 {
		b.WriteString("\n")
	}
	b.WriteString(s.render(n.Children, "\n"))
	return b.String()
}

func (s HTML) RenderKeyword(n *parser.Keyword) string {
	return ""
}

func (s HTML) RenderListItem(n *parser.ListItem) string {
	content := s.render(n.Children, "\n")
	if n.Status != "" {
		content = fmt.Sprintf("%[1]s %[2]s",
			fmt.Sprintf(listitemstatusElement, n.Status),
			content,
		)
	}
	return fmt.Sprintf(listitemElement, content)
}

func (s HTML) RenderList(n *parser.List) string {
	content := s.render(n.Children, "\n")
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

func (s HTML) RenderTableColumn(n *parser.TableColumn) string {
	if n.IsHeader {
		return fmt.Sprintf(tableHeaderElement, "", s.render(n.Children, ""))
	}
	return fmt.Sprintf(tableColumnElement, "", s.render(n.Children, ""))
}

func (s HTML) RenderTableRow(n *parser.TableRow) string {
	if n.Separator {
		return ""
	}
	return fmt.Sprintf(tableRowElement, s.render(n.Children, "\n"))
}

func (s HTML) RenderTable(n *parser.Table) string {
	return fmt.Sprintf(tableElement, s.render(n.Children, "\n"))
}

func (s HTML) RenderBlock(n *parser.Block) string {
	switch n.Type {
	case "SRC":
		if s.Highlight == nil {
			return fmt.Sprintf(srcElement, n.Parameters[0], s.render(n.Children, "\n"))
		}
		return s.Highlight(s.render(n.Children, "\n"), n.Parameters[0])
	case "EXAMPLE":
		return fmt.Sprintf(srcElement, "example", s.render(n.Children, "\n"))
	case "CENTER":
		return fmt.Sprintf(centerElement, s.render(n.Children, "\n"))
	case "QUOTE":
		return fmt.Sprintf(quoteElement, s.render(n.Children, "\n"))
	case "EXPORT":
		return fmt.Sprintf(exportElement, s.render(n.Children, "\n"))
	case "VERSE":
		var b strings.Builder
		for _, child := range n.Children {
			if child.Name() == parser.InlineLineBreakName {
				b.WriteString("<br />\n")
				continue
			}
			b.WriteString(render(s, child))
		}
		return fmt.Sprintf(verseElement, b.String())
	}
	return s.render(n.Children, "\n")
}

func (s HTML) RenderBlockResult(n *parser.BlockResult) string {
	return s.render(n.Children, "\n")
}

func (s HTML) RenderDrawer(n *parser.Drawer) string {
	return s.render(n.Children, "\n")
}

func (s HTML) RenderBlankline(n *parser.Blankline) string {
	return ""
}

func (s HTML) RenderHr(n *parser.Hr) string {
	return "<hr/>"
}

func (s HTML) RenderParagraph(n *parser.Paragragh) string {
	return fmt.Sprintf(paragraghElement, s.render(n.Children, ""))
}

func (s HTML) RenderSection(n *parser.Section) string {
	if len(n.Children) == 0 {
		return ""
	}

	var b strings.Builder

	b.WriteString("<ul>\n")
	for _, section := range n.Children {
		b.WriteString(fmt.Sprintf(`<li><a href="%s">%s</a>`, section.Id(), s.heading(section.Heading)))
		if len(section.Children) > 0 {
			b.WriteString("\n")
			b.WriteString(s.RenderSection(section))
		}
		b.WriteString("</li>\n")
	}
	b.WriteString("</ul>")
	return b.String()
}

func (s HTML) String() string {
	content := s.render(s.Document.Children, "\n")
	if !s.Toc || s.Document.Get("toc") == "nil" {
		return content
	}
	if toc := s.RenderSection(s.Document.Sections); toc != "" {
		return toc + "\n" + content
	}
	return content
}
