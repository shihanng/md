package md

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func TestRaw(t *testing.T) {
	input, err := ioutil.ReadFile(`testdata/raw.md`)
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

	g := goldie.New(t)
	g.Assert(t, "raw", buf.Bytes())
}
