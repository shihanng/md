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
	reg.Register(ast.KindParagraph, r.renderNoop)
	reg.Register(ast.KindTextBlock, r.renderNoop)
	reg.Register(ast.KindThematicBreak, r.renderNoop)

	// inlines

	reg.Register(ast.KindAutoLink, r.renderNoop)
	reg.Register(ast.KindCodeSpan, r.renderNoop)
	reg.Register(ast.KindEmphasis, r.renderNoop)
	reg.Register(ast.KindImage, r.renderNoop)
	reg.Register(ast.KindLink, r.renderNoop)
	reg.Register(ast.KindRawHTML, r.renderNoop)
	reg.Register(ast.KindText, r.renderNoop)
	reg.Register(ast.KindString, r.renderNoop)
}

func (r *Renderer) renderNoop(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}
