package md

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
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
