package md

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type Renderer struct{}

// RegisterFuncs implements github.com/yuin/goldmark/renderer NodeRenderer.RegisterFuncs.
func (r *Renderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// blocks

	reg.Register(ast.KindDocument, r.renderNoop)
	reg.Register(ast.KindHeading, r.renderNoop)
	reg.Register(ast.KindBlockquote, r.renderNoop)
	reg.Register(ast.KindCodeBlock, r.renderNoop)
	reg.Register(ast.KindFencedCodeBlock, r.renderNoop)
	reg.Register(ast.KindHTMLBlock, r.renderNoop)
	reg.Register(ast.KindList, r.renderNoop)
	reg.Register(ast.KindListItem, r.renderNoop)
	reg.Register(ast.KindParagraph, r.renderParagraph)
	reg.Register(ast.KindTextBlock, r.renderNoop)
	reg.Register(ast.KindThematicBreak, r.renderNoop)

	// inlines

	reg.Register(ast.KindAutoLink, r.renderNoop)
	reg.Register(ast.KindCodeSpan, r.renderNoop)
	reg.Register(ast.KindEmphasis, r.renderEmphasis)
	reg.Register(ast.KindImage, r.renderImage)
	reg.Register(ast.KindLink, r.renderNoop)
	reg.Register(ast.KindRawHTML, r.renderNoop)
	reg.Register(ast.KindText, r.renderText)
	reg.Register(ast.KindString, r.renderNoop)
}

func (r *Renderer) renderNoop(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Renderer) renderParagraph(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// TODO: Handle list depth
	if !entering {
		if _, ok := node.NextSibling().(ast.Node); ok && node.FirstChild() != nil {
			_, _ = w.WriteString("\n\n")
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderEmphasis(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
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

func (r *Renderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
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

func (r *Renderer) renderText(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*ast.Text)
	segment := n.Segment
	_, _ = w.Write(segment.Value(source))

	switch {
	case n.HardLineBreak():
		_, _ = w.WriteString(`  `)
		_ = w.WriteByte('\n')
	case n.SoftLineBreak():
		_ = w.WriteByte('\n')
	}

	return ast.WalkContinue, nil
}
