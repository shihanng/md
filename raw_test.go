package md

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func TestRaw(t *testing.T) {
	input, err := ioutil.ReadFile(`testdata/raw_input.md`)
	assert.NoError(t, err)

	expected, err := ioutil.ReadFile(`testdata/raw_expected.md`)
	assert.NoError(t, err)

	reader := text.NewReader(input)

	p := parser.NewParser(
		parser.WithBlockParsers(
			[]util.PrioritizedValue{
				util.Prioritized(NewRawParagraphParser(), 100),
			}...),
	)
	node := p.Parse(reader)
	r := renderer.NewRenderer(renderer.WithNodeRenderers(util.Prioritized(&RawRenderer{}, 1000)))

	var buf bytes.Buffer
	assert.NoError(t, r.Render(&buf, input, node))
	assert.Equal(t, string(expected), buf.String())
}
