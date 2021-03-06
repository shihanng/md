package md

import (
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type list struct {
	marker    byte
	number    int
	isOrdered bool
	offset    int
}

type Renderer struct {
	blockquoteDepth int

	lists []list
}

func (r *Renderer) currentList() list {
	return r.lists[len(r.lists)-1]
}

func (r *Renderer) listString(node ast.Node) string {
	if parent := node.Parent(); parent != nil && parent.Kind() == ast.KindListItem {
		cl := r.currentList()

		if node.PreviousSibling() == nil {
			l := fmt.Sprintf("%c", cl.marker)
			if cl.isOrdered {
				l = fmt.Sprintf("%d%s", cl.number, l)
			}
			return l + strings.Repeat(" ", cl.offset-len(l))
		} else {
			return strings.Repeat(" ", cl.offset)
		}
	}

	return ""
}

// RegisterFuncs implements github.com/yuin/goldmark/renderer NodeRenderer.RegisterFuncs.
func (r *Renderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// blocks

	reg.Register(ast.KindDocument, RenderNoop)
	reg.Register(ast.KindHeading, RenderHeading)
	reg.Register(ast.KindBlockquote, r.RenderBlockquote)
	reg.Register(ast.KindCodeBlock, r.RenderCodeBlock)
	reg.Register(ast.KindFencedCodeBlock, r.RenderFencedCodeBlock)
	reg.Register(ast.KindHTMLBlock, RenderNoop)
	reg.Register(ast.KindList, r.RenderList)
	reg.Register(ast.KindListItem, r.RenderListItem)
	reg.Register(ast.KindParagraph, r.RenderParagraph)
	reg.Register(ast.KindTextBlock, RenderNoop)
	reg.Register(ast.KindThematicBreak, RenderNoop)

	// inlines

	reg.Register(ast.KindAutoLink, RenderAutoLink)
	reg.Register(ast.KindCodeSpan, RenderCodeSpan)
	reg.Register(ast.KindEmphasis, RenderEmphasis)
	reg.Register(ast.KindImage, RenderImage)
	reg.Register(ast.KindLink, RenderLink)
	reg.Register(ast.KindRawHTML, RenderNoop)
	reg.Register(ast.KindText, r.RenderText)
	reg.Register(ast.KindString, RenderNoop)
}

// RenderNoop does nothing.
func RenderNoop(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Renderer) RenderBlockquote(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString(r.listString(node))
		r.blockquoteDepth++
	} else {
		r.blockquoteDepth--

		if n := node.NextSibling(); n != nil && n.Type() == ast.TypeBlock {
			_, _ = w.WriteString(quotes(r.blockquoteDepth))
			_, _ = w.WriteString("\n")
		}
	}

	return ast.WalkContinue, nil
}

func (r *Renderer) RenderCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString(r.listString(node))
		_, _ = w.WriteString("```\n")
		r.writeLines(w, source, node)
	} else {
		_, _ = w.WriteString(r.listString(node))
		_, _ = w.WriteString("```\n")

		if n := node.NextSibling(); n != nil && n.Type() == ast.TypeBlock {
			_, _ = w.WriteString("\n")
		}
	}

	return ast.WalkContinue, nil
}

func (r *Renderer) RenderFencedCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.FencedCodeBlock)
	if entering {
		language := n.Language(source)
		_, _ = w.WriteString(r.listString(node))
		_, _ = w.WriteString(fmt.Sprintf("```%s\n", language))
		r.writeLines(w, source, node)
	} else {
		_, _ = w.WriteString("```\n")

		if n := node.NextSibling(); n != nil && n.Type() == ast.TypeBlock {
			_, _ = w.WriteString("\n")
		}
	}
	return ast.WalkContinue, nil
}

func RenderHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Heading)

	if entering {
		if len(n.Text(source)) == 0 || n.Level > 2 {
			_, _ = w.WriteString(strings.Repeat("#", n.Level))

			if len(n.Text(source)) > 0 {
				_ = w.WriteByte(' ')
			} else {
				_, _ = w.WriteString("\n\n")
			}
		}
	} else if l := len(n.Text(source)); l > 0 {
		switch n.Level {
		case 1:
			_ = w.WriteByte('\n')
			_, _ = w.WriteString(strings.Repeat("=", l))
		case 2:
			_ = w.WriteByte('\n')
			_, _ = w.WriteString(strings.Repeat("-", l))
		}
		_, _ = w.WriteString("\n\n")
	}

	return ast.WalkContinue, nil
}

func (r *Renderer) RenderList(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	l := node.(*ast.List)
	if entering {
		r.lists = append(r.lists, list{
			marker:    l.Marker,
			number:    l.Start,
			isOrdered: l.IsOrdered(),
		})
	} else {
		r.lists = r.lists[:len(r.lists)]
	}

	return ast.WalkContinue, nil
}

func (r *Renderer) RenderListItem(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	li := node.(*ast.ListItem)
	if entering {
		r.lists[len(r.lists)-1].offset = li.Offset
	}

	return ast.WalkContinue, nil
}

func (r *Renderer) RenderParagraph(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		if n := node.NextSibling(); n != nil && n.Type() == ast.TypeBlock {
			_, _ = w.WriteString("\n")
			_, _ = w.WriteString(quotes(r.blockquoteDepth))
			_, _ = w.WriteString("\n")
		} else {
			_, _ = w.WriteString("\n")
		}
	} else {
		_, _ = w.WriteString(r.listString(node))

		if r.blockquoteDepth > 0 {
			_, _ = w.WriteString(quotes(r.blockquoteDepth))
			_, _ = w.WriteString(" ")
		}
	}
	return ast.WalkContinue, nil
}

func RenderAutoLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*ast.AutoLink)
	url := n.URL(source)

	_ = w.WriteByte('<')
	_, _ = w.WriteString(string(url))
	_ = w.WriteByte('>')

	return ast.WalkContinue, nil
}

func RenderCodeSpan(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	_ = w.WriteByte('`')

	return ast.WalkContinue, nil
}

func RenderEmphasis(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Emphasis)

	star := "*"

	if n.Level == 2 {
		star += "*"
	}

	if entering {
		_, _ = w.WriteString(star)
	} else {
		_, _ = w.WriteString(star)
	}

	return ast.WalkContinue, nil
}

func RenderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.Image)
	_, _ = w.WriteString(`![`)
	_, _ = w.Write(n.Text(source))
	_, _ = w.WriteString(`](`)
	_, _ = w.Write(n.Destination)
	if n.Title != nil {
		_, _ = w.WriteString(` "`)
		_, _ = w.Write(n.Title)
		_ = w.WriteByte('"')
	}
	_, _ = w.WriteString(`)`)
	return ast.WalkSkipChildren, nil
}

func RenderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.Link)
	_, _ = w.WriteString(`[`)
	_, _ = w.Write(n.Text(source))
	_, _ = w.WriteString(`](`)
	_, _ = w.Write(n.Destination)
	if n.Title != nil {
		_, _ = w.WriteString(` "`)
		_, _ = w.Write(n.Title)
		_ = w.WriteByte('"')
	}
	_, _ = w.WriteString(`)`)
	return ast.WalkSkipChildren, nil
}

func (r *Renderer) RenderText(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*ast.Text)
	segment := n.Segment
	_, _ = w.Write(segment.Value(source))

	if n.IsRaw() {
		return ast.WalkContinue, nil
	}

	switch {
	case n.HardLineBreak():
		_, _ = w.WriteString("\\\n")
		if r.blockquoteDepth > 0 {
			_, _ = w.WriteString(quotes(r.blockquoteDepth))
			_, _ = w.WriteString(" ")
		}
	case n.SoftLineBreak():
		_ = w.WriteByte(' ')
	}

	return ast.WalkContinue, nil
}

func quotes(level int) string {
	return strings.TrimSuffix(strings.Repeat("> ", level), " ")
}

func (r *Renderer) writeLines(w util.BufWriter, source []byte, n ast.Node) {
	l := n.Lines().Len()
	for i := 0; i < l; i++ {
		line := n.Lines().At(i)
		_, _ = w.WriteString(r.listString(n))
		_, _ = w.WriteString(string(line.Value(source)))
	}
}
