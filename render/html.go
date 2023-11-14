package render

import (
	"fmt"
	"net/url"
	"path/filepath"
	"sort"
	"strings"

	"github.com/honmaple/org-golang/parser"
)

type HTML struct {
	Document           *parser.Document
	Toc                bool
	HeadingOffset      int
	RenderNodeFunc     func(Renderer, parser.Node) string
	RenderFootnoteFunc func(Renderer, []*parser.Footnote, map[string]bool) string

	fnList []*parser.Footnote
	fnUsed map[string]bool
}

var htmlEscaper = strings.NewReplacer(
	`<`, "&lt;",
	`>`, "&gt;",
)

func htmlEscape(s string) string {
	return htmlEscaper.Replace(s)
}

// If def is true, use default RenderNode
func (r *HTML) RenderNode(n parser.Node, def bool) string {
	if def || r.RenderNodeFunc == nil {
		return RenderNode(r, n)
	}
	return r.RenderNodeFunc(r, n)
}

func (r *HTML) RenderNodes(children []parser.Node, sep string) string {
	return RenderNodes(r, children, sep)
}

func (r *HTML) RenderInlineLink(n *parser.InlineLink) string {
	rawURL := n.URL
	if n.Protocol != "" && n.Protocol != "file" {
		rawURL = n.Protocol + "://" + n.URL
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	switch n.Type() {
	case parser.ImageLink:
		return fmt.Sprintf("<img src=\"%s\" alt=\"%s\"/>", rawURL, filepath.Base(parsedURL.Path))
	case parser.VideoLink:
		return fmt.Sprintf("<video src=\"%s\">%s</video>", rawURL, filepath.Base(parsedURL.Path))
	default:
		desc := rawURL
		if n.Desc != "" {
			desc = n.Desc
		}
		return fmt.Sprintf("<a href=\"%s\">%s</a>", rawURL, desc)
	}
}

func (r *HTML) RenderInlineText(n *parser.InlineText) string {
	if n.Raw {
		return n.Content
	}
	return htmlEscape(n.Content)
}

func (r *HTML) RenderInlineEmphasis(n *parser.InlineEmphasis) string {
	text := r.RenderNodes(n.Children, "")
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

func (r *HTML) RenderInlineBackSlash(n *parser.InlineBackSlash) string {
	if n.Break {
		return "<br />"
	}
	return ""
}

func (r *HTML) RenderInlineLineBreak(n *parser.InlineLineBreak) string {
	return strings.Repeat("\n", n.Count)
}

func (r *HTML) RenderInlineTimestamp(n *parser.InlineTimestamp) string {
	return ""
}

func (r *HTML) RenderInlinePercent(n *parser.InlinePercent) string {
	return fmt.Sprintf("<code>[%s]</code>", n.Num)
}

func (r *HTML) heading(n *parser.Heading) string {
	var b strings.Builder

	if n.Keyword != "" {
		b.WriteString(fmt.Sprintf("<span class=\"todo\">%[1]s</span>", n.Keyword))
	}
	if n.Priority != "" {
		b.WriteString(fmt.Sprintf("<span class=\"priority\">%[1]s</span>", n.Priority))
	}
	b.WriteString(r.RenderNodes(n.Title, ""))
	for _, tag := range n.Tags {
		b.WriteString(fmt.Sprintf("<span class=\"tag\">%[1]s</span>", tag))
	}
	return b.String()
}

func (r *HTML) RenderHeading(n *parser.Heading) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("<h%[1]d id=\"%s\">", n.Stars+r.HeadingOffset, n.Id()))
	b.WriteString(r.heading(n))
	b.WriteString(fmt.Sprintf("</h%[1]d>", n.Stars))
	if len(n.Children) > 0 {
		b.WriteString("\n")
	}
	b.WriteString(r.RenderNodes(n.Children, "\n"))
	return b.String()
}

func (r *HTML) RenderKeyword(n *parser.Keyword) string {
	return ""
}

func (r *HTML) RenderListItem(n *parser.ListItem) string {
	content := r.RenderNodes(n.Children, "\n")
	if n.Status != "" {
		if strings.HasPrefix(content, "<p>") {
			content = fmt.Sprintf("<p><code>[%[1]s]</code>", n.Status) + content[3:]
		} else {
			content = fmt.Sprintf("<code>[%[1]s]</code>", n.Status) + content
		}
	}
	return fmt.Sprintf("<li>\n%[1]s</li>", content)
}

func (r *HTML) RenderList(n *parser.List) string {
	content := r.RenderNodes(n.Children, "\n")
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

func (r *HTML) RenderTableColumn(n *parser.TableColumn) string {
	align := ""
	if n.Align != "" {
		align = fmt.Sprintf(` align="%s"`, n.Align)
	}
	if n.IsHeader {
		return fmt.Sprintf("<th%[1]s>%[2]s</th>", align, r.RenderNodes(n.Children, ""))
	}
	return fmt.Sprintf("<td%[1]s>%[2]s</td>", align, r.RenderNodes(n.Children, ""))
}

func (r *HTML) RenderTableRow(n *parser.TableRow) string {
	if n.Separator || len(n.Infos) > 0 {
		return ""
	}
	return fmt.Sprintf("<tr>\n%[1]s\n</tr>", r.RenderNodes(n.Children, "\n"))
}

func (r *HTML) RenderTable(n *parser.Table) string {
	return fmt.Sprintf("<table>\n%[1]s\n</table>", r.RenderNodes(n.Children, "\n"))
}

func (r *HTML) RenderBlock(n *parser.Block) string {
	switch n.Type {
	case "SRC":
		lang := "unknown"
		if len(n.Parameters) > 0 {
			lang = n.Parameters[0]
		}
		text := htmlEscape(DedentString(r.RenderNodes(n.Children, "\n")))
		return fmt.Sprintf("<pre class=\"src src-%[1]s\">%[2]s</pre>", lang, text)
	case "EXAMPLE":
		text := htmlEscape(DedentString(r.RenderNodes(n.Children, "\n")))
		return fmt.Sprintf("<pre class=\"src src-example\">%[1]s</pre>", text)
	case "CENTER":
		return fmt.Sprintf("<div style=\"text-align:center;\">\n%[1]s\n</div>", r.RenderNodes(n.Children, "\n"))
	case "QUOTE":
		return fmt.Sprintf("<blockquote>\n%[1]s\n</blockquote>", r.RenderNodes(n.Children, "\n"))
	case "EXPORT":
		return r.RenderNodes(n.Children, "\n")
	case "VERSE":
		var b strings.Builder
		for _, child := range n.Children {
			if child.Name() == parser.InlineLineBreakName {
				b.WriteString(strings.Repeat("<br />", child.(*parser.InlineLineBreak).Count))
				continue
			}
			b.WriteString(r.RenderNode(child, false))
		}
		return fmt.Sprintf("<p>\n%[1]s\n</p>", b.String())
	}
	return fmt.Sprintf("<div class=\"%[1]s-block\">\n%[2]s\n</div>", strings.ToLower(n.Type), r.RenderNodes(n.Children, "\n"))
}

func (r *HTML) RenderBlockResult(n *parser.BlockResult) string {
	return r.RenderNodes(n.Children, "\n")
}

func (r *HTML) RenderDrawer(n *parser.Drawer) string {
	return r.RenderNodes(n.Children, "\n")
}

func (r *HTML) RenderBlankline(n *parser.Blankline) string {
	return ""
}

func (r *HTML) RenderHr(n *parser.Hr) string {
	return "<hr/>"
}

func (r *HTML) RenderParagraph(n *parser.Paragragh) string {
	return fmt.Sprintf("<p>\n%[1]s\n</p>", r.RenderNodes(n.Children, ""))
}

func (r *HTML) RenderFootnote(n *parser.Footnote) string {
	if len(n.Definition) > 0 {
		r.fnList = append(r.fnList, n)
	}
	if n.Inline && n.Label != "" {
		r.fnUsed[n.Label] = true
		return fmt.Sprintf(`<sup><a id="fnr.%[1]s" href="#fn.%[1]s">[%[1]s]</a></sup>`, n.Label)
	}
	return ""
}

func (r *HTML) RenderFootnotes(fns []*parser.Footnote, used map[string]bool) string {
	if r.RenderFootnoteFunc != nil {
		return r.RenderFootnoteFunc(r, fns, used)
	}
	if len(fns) == 0 || len(used) == 0 {
		return ""
	}
	sort.SliceStable(fns, func(i, j int) bool {
		return fns[i].Label < fns[j].Label
	})
	var b strings.Builder

	b.WriteString("<div id=\"footnotes\"><h2 class=\"footnotes\">Footnotes</h2>\n<ol id=\"text-footnotes\">\n")
	for _, fn := range fns {
		if len(fn.Definition) == 0 || !used[fn.Label] {
			continue
		}
		b.WriteString("<li>")
		b.WriteString(fmt.Sprintf(`<sup><a id="fn.%[1]s" href="#fnr.%[1]s">%[1]s</a></sup>`, fn.Label))
		b.WriteString("<div style=\"display: inline-grid;\">")
		b.WriteString(r.RenderNodes(fn.Definition, "\n"))
		b.WriteString("</div>")
		b.WriteString("\n</li>\n")
	}
	b.WriteString("</ol></div>")
	return b.String()
}

func (r *HTML) RenderSection(n *parser.Section) string {
	if len(n.Children) == 0 {
		return ""
	}

	var b strings.Builder

	b.WriteString("<ul>\n")
	for _, section := range n.Children {
		b.WriteString(fmt.Sprintf(`<li><a href="#%s">%s</a>`, section.Id(), r.heading(section.Heading)))
		if len(section.Children) > 0 {
			b.WriteString("\n")
			b.WriteString(r.RenderSection(section))
		}
		b.WriteString("</li>\n")
	}
	b.WriteString("</ul>")
	return b.String()
}

func (r *HTML) String() string {
	r.fnUsed = make(map[string]bool)

	content := r.RenderNodes(r.Document.Children, "\n")
	if fn := r.RenderFootnotes(r.fnList, r.fnUsed); fn != "" {
		content = content + fn
	}

	if !r.Toc || r.Document.Get("toc") == "nil" {
		return content
	}
	if toc := r.RenderNode(r.Document.Sections, false); toc != "" {
		toc = fmt.Sprintf(`<div id="table-of-contents"><h2>Table of Contents</h2><div id="text-table-of-contents">%s</div></div>`, toc)
		return toc + "\n" + content
	}
	return content
}
