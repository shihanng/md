package md

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var debug = flag.Bool("debug", false, "Debug with ast.Walk")

func TestRenderer(t *testing.T) {
	testCases := []string{
		"standard_renderer",
		"blockquotes",
		"codeblocks",
		"list",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			input, err := ioutil.ReadFile("testdata/" + tc + ".md")
			assert.NoError(t, err)

			var buf bytes.Buffer
			reader := text.NewReader(input)
			node := goldmark.DefaultParser().Parse(reader)
			r := renderer.NewRenderer(renderer.WithNodeRenderers(util.Prioritized(&Renderer{}, 1000)))

			printNodes(node, input)

			assert.NoError(t, r.Render(&buf, input, node))

			g := goldie.New(t)
			g.Assert(t, tc, buf.Bytes())
		})
	}
}

func TestReRenderer(t *testing.T) {
	testCases := []string{
		"standard_renderer",
		"blockquotes",
		"codeblocks",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			input, err := ioutil.ReadFile("testdata/" + tc + ".golden")
			assert.NoError(t, err)

			var buf bytes.Buffer
			reader := text.NewReader(input)
			node := goldmark.DefaultParser().Parse(reader)
			r := renderer.NewRenderer(renderer.WithNodeRenderers(util.Prioritized(&Renderer{}, 1000)))

			printNodes(node, input)

			assert.NoError(t, r.Render(&buf, input, node))

			g := goldie.New(t)
			g.Assert(t, tc, buf.Bytes())
		})
	}
}

func printNodes(node ast.Node, source []byte) {
	if !*debug {
		return
	}

	space := 0
	if err := ast.Walk(node, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			text := node.Text(source)
			fmt.Printf("%s%s: '%.10s'\n", strings.Repeat(" ", space), node.Kind(), text)
			space++
		} else {
			space--
			fmt.Printf("%s%s\n", strings.Repeat(" ", space), node.Kind())
		}

		return ast.WalkContinue, nil
	}); err != nil {
		fmt.Println(err)
	}
}
