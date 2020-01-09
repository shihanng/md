package md

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type rawParagraphParser struct {
}

// NewRawParagraphParser returns a new BlockParser that
// parses paragraphs without any white space trimming.
func NewRawParagraphParser() parser.BlockParser {
	return &rawParagraphParser{}
}

func (b *rawParagraphParser) Trigger() []byte {
	return nil
}

func (b *rawParagraphParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	_, segment := reader.PeekLine()
	if segment.IsEmpty() {
		return nil, parser.NoChildren
	}
	node := ast.NewParagraph()
	node.Lines().Append(segment)
	reader.Advance(segment.Len() - 1)
	return node, parser.NoChildren
}

func (b *rawParagraphParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	_, segment := reader.PeekLine()
	if segment.IsEmpty() {
		return parser.Close
	}
	node.Lines().Append(segment)
	reader.Advance(segment.Len() - 1)
	return parser.Continue | parser.NoChildren
}

func (b *rawParagraphParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
	parent := node.Parent()
	if parent == nil {
		// paragraph has been transformed
		return
	}
	lines := node.Lines()
	if lines.Len() != 0 {
		length := lines.Len()
		lastLine := node.Lines().At(length - 1)
		node.Lines().Set(length-1, lastLine)
	}
	if lines.Len() == 0 {
		node.Parent().RemoveChild(node.Parent(), node)
		return
	}
}

func (b *rawParagraphParser) CanInterruptParagraph() bool {
	return false
}

func (b *rawParagraphParser) CanAcceptIndentedLine() bool {
	return false
}

type RawRenderer struct{}

// RegisterFuncs implements github.com/yuin/goldmark/renderer NodeRenderer.RegisterFuncs.
func (r *RawRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// blocks

	reg.Register(ast.KindParagraph, r.renderParagraph)

	// inlines

	reg.Register(ast.KindText, r.renderText)
}

func (r *RawRenderer) renderParagraph(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		if _, ok := node.NextSibling().(ast.Node); ok && node.FirstChild() != nil {
			_, _ = w.WriteString("\n\n")
		}
	}
	return ast.WalkContinue, nil
}

func (r *RawRenderer) renderText(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
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
