package md

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func TestRenderer(t *testing.T) {
	source := []byte(`test`)

	var buf bytes.Buffer
	reader := text.NewReader(source)
	node := goldmark.DefaultParser().Parse(reader)
	r := renderer.NewRenderer(renderer.WithNodeRenderers(util.Prioritized(&Renderer{}, 1000)))

	assert.NoError(t, r.Render(&buf, source, node))
	assert.Empty(t, buf.Bytes())
}
