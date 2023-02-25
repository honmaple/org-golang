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
	CheckLink func(string) parser.LinkType
}

func (s HTML) render(children []parser.Node, sep string) string {
	return concat(s, children, sep)
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
	var typ parser.LinkType
	if s.CheckLink != nil {
		typ = s.CheckLink(n.URL)
	} else {
		typ = n.Type()
	}
	switch typ {
	case parser.ImageLink:
		return fmt.Sprintf("<img src=\"%s\"/>", n.URL)
	case parser.VedioLink:
		return fmt.Sprintf("<video src=\"%s\">%s</video>", n.URL, n.URL)
	default:
		if n.Desc == "" {
			return fmt.Sprintf("<a href=\"%s\">%s</a>", n.URL, n.URL)
		}
		return fmt.Sprintf("<a href=\"%s\">%s</a>", n.URL, n.Desc)
	}
}

func (s HTML) RenderEmphasis(n *parser.InlineEmphasis) string {
	text := s.render(n.Children, "")
	switch n.Marker {
	case "=", "~", "`":
		return fmt.Sprintf("<code>%s</code>", text)
	case "*":
		return fmt.Sprintf("<b>%s</b>", text)
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
		b.WriteString(fmt.Sprintf("<span class=\"todo\">%[1]s</span>", n.Keyword))
	}
	if n.Priority != "" {
		b.WriteString(fmt.Sprintf("<span class=\"priority\">%[1]s</span>", n.Priority))
	}
	b.WriteString(s.render(n.Title, ""))
	for _, tag := range n.Tags {
		b.WriteString(fmt.Sprintf("<span class=\"tag\">%[1]s</span>", tag))
	}
	return b.String()
}

func (s HTML) RenderHeading(n *parser.Heading) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("<h%[1]d id=\"%s\">", n.Stars+s.Offset, n.Id()))
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
		content = fmt.Sprintf("<code>%[1]s</code>", n.Status) + " " + content
	}
	return fmt.Sprintf("<li>\n%[1]s</li>", content)
}

func (s HTML) RenderList(n *parser.List) string {
	content := s.render(n.Children, "\n")
	switch n.Type {
	case parser.OrderlistName:
		return fmt.Sprintf("<ol>\n%[1]s\n</ol>", content)
	case parser.UnorderlistName:
		return fmt.Sprintf("<ul>\n%[1]s\n</ul>", content)
	case parser.DescriptiveName:
		return fmt.Sprintf("<dl>\n%[1]s\n</dl>", content)
	default:
		return ""
	}
}

func (s HTML) RenderTableColumn(n *parser.TableColumn) string {
	if n.IsHeader {
		return fmt.Sprintf("<th class=\"align-%[1]s\">%[2]s</th>", "", s.render(n.Children, ""))
	}
	return fmt.Sprintf("<td class=\"align-%[1]s\">%[2]s</td>", "", s.render(n.Children, ""))
}

func (s HTML) RenderTableRow(n *parser.TableRow) string {
	if n.Separator {
		return ""
	}
	return fmt.Sprintf("<tr>\n%[1]s\n</tr>", s.render(n.Children, "\n"))
}

func (s HTML) RenderTable(n *parser.Table) string {
	return fmt.Sprintf("<table>\n%[1]s\n</table>", s.render(n.Children, "\n"))
}

func (s HTML) RenderBlock(n *parser.Block) string {
	switch n.Type {
	case "SRC":
		lang := ""
		if len(n.Parameters) > 0 {
			lang = n.Parameters[0]
		}
		text := dedent(s.render(n.Children, "\n"))
		if s.Highlight == nil {
			return fmt.Sprintf("<pre class=\"src src-%[1]s\">%[2]s</pre>", lang, text)
		}
		return s.Highlight(text, lang)
	case "EXAMPLE":
		text := dedent(s.render(n.Children, "\n"))
		return fmt.Sprintf("<pre class=\"src src-example\">%[1]s</pre>", text)
	case "CENTER":
		return fmt.Sprintf("<div style=\"text-align:center;\">\n%[1]s\n</div>", s.render(n.Children, "\n"))
	case "QUOTE":
		return fmt.Sprintf("<blockquote>\n%[1]s\n</blockquote>", s.render(n.Children, "\n"))
	case "EXPORT":
		return s.render(n.Children, "\n")
	case "VERSE":
		var b strings.Builder
		for _, child := range n.Children {
			if child.Name() == parser.InlineLineBreakName {
				b.WriteString("<br />")
				continue
			}
			b.WriteString(render(s, child))
		}
		return fmt.Sprintf("<p>\n%[1]s\n</p>", b.String())
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
	return fmt.Sprintf("<p>\n%[1]s\n</p>", s.render(n.Children, ""))
}

func (s HTML) RenderSection(n *parser.Section) string {
	if len(n.Children) == 0 {
		return ""
	}

	var b strings.Builder

	b.WriteString("<ul>\n")
	for _, section := range n.Children {
		b.WriteString(fmt.Sprintf(`<li><a href="#%s">%s</a>`, section.Id(), s.heading(section.Heading)))
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
		toc = fmt.Sprintf(`<div id="table-of-contents"><h2>Table of Contents</h2><div id="text-table-of-contents">%s</div></div>`, toc)
		return toc + "\n" + content
	}
	return content
}
