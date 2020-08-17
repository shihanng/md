package md

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func TestRenderer(t *testing.T) {
	testCases := []string{
		"standard_renderer",
		"blockquotes",
		"codeblocks",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			input, err := ioutil.ReadFile("testdata/" + tc + ".md")
			assert.NoError(t, err)

			var buf bytes.Buffer
			reader := text.NewReader(input)
			node := goldmark.DefaultParser().Parse(reader)
			r := renderer.NewRenderer(renderer.WithNodeRenderers(util.Prioritized(&Renderer{}, 1000)))

			assert.NoError(t, r.Render(&buf, input, node))

			g := goldie.New(t)
			g.Assert(t, tc, buf.Bytes())
		})
	}
}
