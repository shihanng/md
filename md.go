package md

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
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
	reg.Register(ast.KindImage, r.renderNoop)
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

type emphasisParser struct {
}

var defaultEmphasisParser = &emphasisParser{}

// NewEmphasisParser return a new InlineParser that parses emphasises.
func NewEmphasisParser() parser.InlineParser {
	return defaultEmphasisParser
}

func (s *emphasisParser) Trigger() []byte {
	return []byte{'*', '_'}
}

func (s *emphasisParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	before := block.PrecendingCharacter()
	line, segment := block.PeekLine()
	node := parser.ScanDelimiter(line, before, 1, defaultEmphasisDelimiterProcessor)
	if node == nil {
		return nil
	}
	node.Segment = segment.WithStop(segment.Start + node.OriginalLength)
	block.Advance(node.OriginalLength)
	pc.PushDelimiter(node)
	return node
}

type emphasisDelimiterProcessor struct {
}

func (p *emphasisDelimiterProcessor) IsDelimiter(b byte) bool {
	return b == '*' || b == '_'
}

func (p *emphasisDelimiterProcessor) CanOpenCloser(opener, closer *parser.Delimiter) bool {
	return opener.Char == closer.Char
}

func (p *emphasisDelimiterProcessor) OnMatch(consumes int) ast.Node {
	return ast.NewEmphasis(consumes)
}

var defaultEmphasisDelimiterProcessor = &emphasisDelimiterProcessor{}
