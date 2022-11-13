package atlassianwikirenderer

import (
	"bufio"
	"fmt"
	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"io"
	"strings"
)

type atlassianRenderer struct {
	nodeRendererFuncs map[ast.NodeKind]renderer.NodeRendererFunc
}

func New() renderer.Renderer {
	r := &atlassianRenderer{}
	r.nodeRendererFuncs = make(map[ast.NodeKind]renderer.NodeRendererFunc)
	r.nodeRendererFuncs[ast.KindDocument] = r.renderDocument
	r.nodeRendererFuncs[ast.KindHeading] = r.renderHeading
	r.nodeRendererFuncs[ast.KindFencedCodeBlock] = r.renderFencedCodeBlock
	r.nodeRendererFuncs[ast.KindBlockquote] = r.renderBlockQuote
	r.nodeRendererFuncs[ast.KindEmphasis] = r.renderEmphasis
	r.nodeRendererFuncs[east.KindTaskCheckBox] = r.renderTaskCheckBox
	r.nodeRendererFuncs[ast.KindText] = r.renderText
	r.nodeRendererFuncs[ast.KindParagraph] = r.renderParagraph
	r.nodeRendererFuncs[ast.KindCodeBlock] = r.renderCodeBlock
	r.nodeRendererFuncs[ast.KindCodeSpan] = r.renderCodeSpan
	r.nodeRendererFuncs[ast.KindAutoLink] = r.renderAutoLink
	r.nodeRendererFuncs[ast.KindLink] = r.renderLink
	r.nodeRendererFuncs[east.KindStrikethrough] = r.renderStrikeThrough
	r.nodeRendererFuncs[ast.KindList] = r.renderList
	r.nodeRendererFuncs[ast.KindListItem] = r.renderListItem
	r.nodeRendererFuncs[ast.KindTextBlock] = r.renderTextBlock
	r.nodeRendererFuncs[east.KindTableCell] = r.renderTableCell
	r.nodeRendererFuncs[east.KindTableRow] = r.renderTableRow
	r.nodeRendererFuncs[east.KindTableHeader] = r.renderTableHeader
	r.nodeRendererFuncs[east.KindTable] = r.renderTable
	r.nodeRendererFuncs[ast.KindImage] = r.renderImage
	r.nodeRendererFuncs[east.KindFootnote] = r.renderFootNote
	r.nodeRendererFuncs[east.KindFootnoteLink] = r.renderFootNoteLink
	r.nodeRendererFuncs[east.KindFootnoteList] = r.renderFootNoteList
	r.nodeRendererFuncs[east.KindDefinitionTerm] = r.renderDefinitionTerm
	r.nodeRendererFuncs[east.KindDefinitionDescription] = r.renderDefinitionDescription
	r.nodeRendererFuncs[ast.KindHTMLBlock] = r.renderHtmlBlock
	r.nodeRendererFuncs[ast.KindRawHTML] = r.renderRawHTML
	r.nodeRendererFuncs[ast.KindString] = r.renderString

	r.nodeRendererFuncs[ast.KindThematicBreak] = r.renderThematicBreak

	return r
}

/*func debugMessage(m string) {
	fmt.Println(m)
}*/
func (r *atlassianRenderer) renderTaskCheckBox(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*east.TaskCheckBox)
		if n.IsChecked {
			_, _ = w.WriteString("[x] ")
		} else {
			_, _ = w.WriteString("[  ] ")
		}
	}
	return ast.WalkContinue, nil
}
func (r *atlassianRenderer) renderThematicBreak(w util.BufWriter, _ []byte, _ ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	writeNewline(w)
	_, _ = renderHorizontalLine(w)
	writeNewline(w)
	return ast.WalkContinue, nil
}
func (r *atlassianRenderer) renderString(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.String)
	_, _ = w.Write(n.Value)
	return ast.WalkContinue, nil
}

const (
	emphasis       = "_"
	newline        = "\n"
	citation       = "??"
	inserted       = "+"
	deleted        = "-"
	monospaceStart = "{{"
	monospaceEnd   = "}}"
	strong         = "*"
	subscript      = "~"
	superscript    = "^"
	horizontalLine = "----"
)

var htmlWikiMapping = map[string]string{
	"<br>":      newline,
	"<cite>":    citation,
	"</cite>":   citation,
	"<code>":    "{{",
	"</code>":   "}}",
	"<del>":     deleted,
	"</del>":    deleted,
	"<s>":       deleted,
	"</s>":      deleted,
	"<ins>":     inserted,
	"</ins>":    inserted,
	"<em>":      emphasis,
	"</em>":     emphasis,
	"<dfn>":     emphasis,
	"</dfn>":    emphasis,
	"<i>":       emphasis,
	"</i>":      emphasis,
	"<kbd>":     monospaceStart,
	"</kbd>":    monospaceEnd,
	"<q>":       "\"",
	"</q>":      "\"",
	"<strong>":  strong,
	"</strong>": strong,
	"<sub>":     subscript,
	"</sub>":    subscript,
	"<sup>":     superscript,
	"</sup>":    superscript,
}

func (r *atlassianRenderer) renderRawHTML(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.RawHTML)
		lines := n.Segments
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			content := string(line.Value(source))
			if wiki, ok := htmlWikiMapping[content]; ok {
				_, _ = w.WriteString(wiki)
			} else {
				_, _ = w.WriteString(content)
			}
		}
	}
	return ast.WalkContinue, nil

}
func (r *atlassianRenderer) renderDefinitionDescription(w util.BufWriter, _ []byte, _ ast.Node, entering bool) (ast.WalkStatus, error) {

	if entering {
		writeNewline(w)
		_, _ = w.WriteString("--  ")
	}
	return ast.WalkContinue, nil

}
func (r *atlassianRenderer) renderDefinitionTerm(w util.BufWriter, _ []byte, _ ast.Node, entering bool) (ast.WalkStatus, error) {

	if entering {
		writeNewline(w)
		_, _ = w.WriteString("-  ")
	}
	return ast.WalkContinue, nil

}
func (r *atlassianRenderer) renderFootNoteList(w util.BufWriter, _ []byte, _ ast.Node, entering bool) (ast.WalkStatus, error) {

	if entering {
		writeNewline(w)
		_, _ = renderHorizontalLine(w)
	}
	return ast.WalkContinue, nil

}

func renderHorizontalLine(w util.BufWriter) (int, error) {
	return w.WriteString(horizontalLine)
}
func (r *atlassianRenderer) renderFootNoteLink(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*east.FootnoteLink)
		for c := n.OwnerDocument().LastChild(); c != nil; c = c.PreviousSibling() { //Find the list of footnotes, start from the end of the document.
			if c.Kind() == east.KindFootnoteList {
				list := c.(*east.FootnoteList)
				for fn := list.FirstChild(); fn != nil; fn = fn.NextSibling() { //Find the footnote the current link refers to
					if fn.Kind() == east.KindFootnote {
						footNote := fn.(*east.Footnote)
						if footNote.Index == n.Index {
							_, _ = w.WriteString(fmt.Sprintf("[#%s]", footNote.Ref))
						}
					}
				}
			}
		}
	}
	return ast.WalkContinue, nil

}
func (r *atlassianRenderer) renderFootNote(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {

	if entering {
		n := node.(*east.Footnote)
		if n.HasBlankPreviousLines() {
			writeNewline(w)
		}
		_, _ = w.WriteString(fmt.Sprintf("{anchor:%s}%s: ", n.Ref, n.Ref))
	}
	return ast.WalkContinue, nil

}

func (r *atlassianRenderer) renderImage(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Image)
	if entering {
		_, _ = w.WriteString(fmt.Sprintf("!%s!", n.Destination))
	}
	return ast.WalkSkipChildren, nil

}
func (r *atlassianRenderer) renderTableRow(w util.BufWriter, _ []byte, _ ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("|")
	} else {
		_, _ = w.WriteString("|")
		writeNewline(w)
	}
	return ast.WalkContinue, nil
}
func (r *atlassianRenderer) renderTableHeader(w util.BufWriter, _ []byte, _ ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("||")
	} else {
		_, _ = w.WriteString("||")
		writeNewline(w)
	}

	return ast.WalkContinue, nil
}
func (r *atlassianRenderer) renderTableCell(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering && !node.HasChildren() {
		//if this is an empty cell then add a space
		/*
			TableCell {
			   	RawText: ""
				HasBlankPreviousLines: false
			}
		*/
		_, _ = w.WriteString(" ")

	}
	if !entering && node.NextSibling() != nil {
		// Add separators between cells only, if the next sibling is nil then it's the last cell in the row.
		rowType := node.Parent().Kind()
		if rowType == east.KindTableRow {
			_, _ = w.WriteString("|")
		} else {
			_, _ = w.WriteString("||")
		}
	}
	return ast.WalkContinue, nil
}

func (r *atlassianRenderer) renderTable(_ util.BufWriter, _ []byte, _ ast.Node, _ bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}
func (r *atlassianRenderer) renderTextBlock(_ util.BufWriter, _ []byte, _ ast.Node, _ bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}
func (r *atlassianRenderer) renderListItem(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	listItem := node.(*ast.ListItem)

	if entering {
		level := 1
		list := listItem.Parent()

		for ancestor := list.Parent(); ancestor != nil; ancestor = ancestor.Parent().Parent() { //First parent is always List
			if _, ok := ancestor.(*ast.ListItem); ok {
				level = level + 1
			} else {
				break
			}
		}
		tag := "*"
		if listNode, ok := list.(*ast.List); ok && listNode.IsOrdered() {
			tag = "#"
		}

		_, _ = w.WriteString(strings.Repeat(tag, level))
		_, _ = w.WriteString(strings.Repeat(" ", listItem.Offset))

	} else {

		hasNestedList := false
		for c := listItem.FirstChild(); c != nil; c = c.NextSibling() {
			if _, ok := c.(*ast.List); ok {
				hasNestedList = true
			}
		}
		if !hasNestedList {
			writeNewline(w)
		}
	}
	return ast.WalkContinue, nil
}
func (r *atlassianRenderer) renderList(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {

	if entering {
		if node.Parent().Kind() == ast.KindListItem {
			writeNewline(w)
		}
	}
	return ast.WalkContinue, nil
}

func (r *atlassianRenderer) renderStrikeThrough(w util.BufWriter, _ []byte, _ ast.Node, entering bool) (ast.WalkStatus, error) {
	_, _ = w.WriteString(deleted)
	return ast.WalkContinue, nil
}
func (r *atlassianRenderer) renderAutoLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.AutoLink)
	address := string(n.URL(source))
	if n.AutoLinkType == ast.AutoLinkEmail {
		address = "mailto:" + address
	}
	if entering {
		_, _ = w.WriteString(fmt.Sprintf("[%v]", address))
	}

	return ast.WalkContinue, nil
}
func (r *atlassianRenderer) renderLink(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Link)
	if entering {
		_, _ = w.WriteString("[")
	} else {
		_, _ = w.WriteString(fmt.Sprintf("|%s]", n.Destination))
	}

	return ast.WalkContinue, nil
}
func (r *atlassianRenderer) renderHtmlBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("{code:html}")
		_, _ = w.Write(node.Text(source))
		renderLines(w, source, node.Lines())

	} else {
		_, _ = w.WriteString("{code}\n")
	}

	return ast.WalkContinue, nil
}
func (r *atlassianRenderer) renderCodeSpan(w util.BufWriter, _ []byte, _ ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString(monospaceStart)
	} else {
		_, _ = w.WriteString(monospaceEnd)
	}

	return ast.WalkContinue, nil
}
func (r *atlassianRenderer) renderCodeBlock(w util.BufWriter, _ []byte, _ ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("{code}")

	} else {
		_, _ = w.WriteString("{code}\n")
	}

	return ast.WalkContinue, nil
}
func (r *atlassianRenderer) renderDocument(_ util.BufWriter, _ []byte, _ ast.Node, _ bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}
func (r *atlassianRenderer) renderParagraph(w util.BufWriter, _ []byte, node ast.Node, _ bool) (ast.WalkStatus, error) {
	ignoreParents := []ast.NodeKind{ast.KindListItem, east.KindFootnote}
	for _, parent := range ignoreParents {
		if node.Parent().Kind() == parent {
			return ast.WalkContinue, nil
		}
	}
	writeNewline(w)
	return ast.WalkContinue, nil
}

func writeNewline(w util.BufWriter) {
	_ = w.WriteByte('\n')
}
func (r *atlassianRenderer) renderText(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.Text)
		_, _ = w.WriteString(string(n.Text(source)))
		if n.HardLineBreak() || n.SoftLineBreak() {
			_, _ = w.WriteString("\n\n")
		}
	}
	return ast.WalkContinue, nil
}

func (r *atlassianRenderer) renderEmphasis(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Emphasis)
	tag := emphasis
	if n.Level == 2 {
		tag = strong
	}
	if entering {
		_, _ = w.WriteString(tag)
	} else {
		_, _ = w.WriteString(tag)
	}
	return ast.WalkContinue, nil
}

func (r *atlassianRenderer) renderBlockQuote(w util.BufWriter, _ []byte, _ ast.Node, _ bool) (ast.WalkStatus, error) {
	_, _ = w.WriteString("{quote}")
	return ast.WalkContinue, nil
}

func (r *atlassianRenderer) renderFencedCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.FencedCodeBlock)
		lang := n.Language(source)
		if len(lang) == 0 {
			_, _ = w.WriteString("{code}")

		} else {
			_, _ = w.WriteString(fmt.Sprintf("{code:%s}", lang))
		}
		renderLines(w, source, n.Lines())
	} else {
		_, _ = w.WriteString("{code}\n")
	}
	return ast.WalkContinue, nil
}

func renderLines(w util.BufWriter, source []byte, lines *text.Segments) {
	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		content := string(line.Value(source))
		_, _ = w.WriteString(content)
	}
}

func (r *atlassianRenderer) renderHeading(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.Heading)
		_, _ = w.WriteString(fmt.Sprintf("h%v.", n.Level))
	} else {
		writeNewline(w)
		writeNewline(w)
	}
	return ast.WalkContinue, nil
}

func (r *atlassianRenderer) Render(w io.Writer, source []byte, n ast.Node) error {
	//debugMessage("Received request to render the document, staring the walk")
	//n.Dump(source, 2)
	writer := bufio.NewWriter(w)
	err := ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		s := ast.WalkContinue
		var err error
		f := r.nodeRendererFuncs[n.Kind()]

		if f != nil {
			s, err = f(writer, source, n, entering)
		} /* else {
			debugMessage(fmt.Sprintf("Node Kind %v, render func not found , source dump \n", n.Kind()))
			n.Dump(source, 2)

		}*/
		return s, err
	})

	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (r *atlassianRenderer) AddOptions(_ ...renderer.Option) {
	//no actions
}

// RegisterFuncs implements NodeRenderer.RegisterFuncs .
func (r *atlassianRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	for key, v := range r.nodeRendererFuncs {
		reg.Register(key, v)
	}
}
