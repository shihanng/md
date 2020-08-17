package md

import (
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type Renderer struct {
	blockquoteDepth int
}

// RegisterFuncs implements github.com/yuin/goldmark/renderer NodeRenderer.RegisterFuncs.
func (r *Renderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// blocks

	reg.Register(ast.KindDocument, RenderNoop)
	reg.Register(ast.KindHeading, RenderHeading)
	reg.Register(ast.KindBlockquote, r.RenderBlockquote)
	reg.Register(ast.KindCodeBlock, RenderCodeBlock)
	reg.Register(ast.KindFencedCodeBlock, RenderFencedCodeBlock)
	reg.Register(ast.KindHTMLBlock, RenderNoop)
	reg.Register(ast.KindList, RenderNoop)
	reg.Register(ast.KindListItem, RenderNoop)
	reg.Register(ast.KindParagraph, r.RenderParagraph)
	reg.Register(ast.KindTextBlock, RenderNoop)
	reg.Register(ast.KindThematicBreak, RenderNoop)

	// inlines

	reg.Register(ast.KindAutoLink, RenderNoop)
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

func RenderCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("```\n")
		writeLines(w, source, node)
	} else {
		_, _ = w.WriteString("```\n")

		if n := node.NextSibling(); n != nil && n.Type() == ast.TypeBlock {
			_, _ = w.WriteString("\n")
		}
	}

	return ast.WalkContinue, nil
}

func RenderFencedCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.FencedCodeBlock)
	if entering {
		language := n.Language(source)
		_, _ = w.WriteString(fmt.Sprintf("```%s\n", language))
		writeLines(w, source, node)
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

func (r *Renderer) RenderParagraph(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// TODO: Handle list depth
	if !entering {
		if n := node.NextSibling(); n != nil && n.Type() == ast.TypeBlock {
			_, _ = w.WriteString("\n")
			_, _ = w.WriteString(quotes(r.blockquoteDepth))
			_, _ = w.WriteString("\n")
		} else {
			_, _ = w.WriteString("\n")
		}
	} else {
		if r.blockquoteDepth > 0 {
			_, _ = w.WriteString(quotes(r.blockquoteDepth))
			_, _ = w.WriteString(" ")
		}
	}
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

func writeLines(w util.BufWriter, source []byte, n ast.Node) {
	l := n.Lines().Len()
	for i := 0; i < l; i++ {
		line := n.Lines().At(i)
		_, _ = w.WriteString(string(line.Value(source)))
	}
}
