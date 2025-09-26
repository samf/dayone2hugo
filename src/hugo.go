package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/md"
)

type FrontMatter struct {
	Date  string   `toml:"date"`
	Title string   `toml:"title"`
	Tags  []string `toml:"tags"`
}

type HugoCmd struct {
	CommonConvert

	KeepTitle bool   `help:"keep the '# title' in your markdown"`
	UseFigure bool   `help:"use a <figure> shortcode for images" env:"HUGO_FIGURE"`
	FigureTag string `help:"shortcode to use for figure" default:"figure" env:"HUGO_FIGURE_TAG"`
}

func (hc *HugoCmd) Run(ctx *Context) error {
	defer hc.Input.Close()

	doz, _, entry, _, body, err := hc.getStuff(ctx)
	if err != nil {
		return err
	}

	frontMatter := hc.NewFrontMatter(entry, body)

	err = hc.GotoOutDir()
	if err != nil {
		return err
	}

	var renderer markdown.Renderer
	if hc.UseFigure {
		renderer = NewHugoRenderer()
	}

	err = hc.outBody(body, frontMatter.output(), renderer)

	err = hc.outPhotos(doz, entry)
	if err != nil {
		return err
	}

	return nil
}

func (hc *HugoCmd) NewFrontMatter(
	entry *Entry,
	body *Body,
) *FrontMatter {
	title := body.findTitle(!hc.KeepTitle)

	front := &FrontMatter{
		Date:  entry.CreationDate,
		Tags:  entry.Tags,
		Title: title,
	}

	return front
}

func (fm *FrontMatter) output() []byte {
	var out bytes.Buffer
	out.Write([]byte("+++\n"))
	encoder := toml.NewEncoder(&out)
	err := encoder.Encode(fm)
	if err != nil {
		log.Printf("cannot encode frontmatter: %w", err)
		panic(err)
	}
	out.Write([]byte("+++\n"))

	return out.Bytes()
}

type HugoRenderer struct {
	markdownRenderer markdown.Renderer
}

func NewHugoRenderer() markdown.Renderer {
	return &HugoRenderer{
		markdownRenderer: md.NewRenderer(md.WithRenderInFooter(true)),
	}
}

const figureFmt = `
{{< %s
  src=%q
>}}
`

func (hr *HugoRenderer) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	newNode := node
	switch node := node.(type) {
	case *ast.Image:
		newNode = &ast.HTMLBlock{
			Leaf: ast.Leaf{
				Parent:    node.Parent,
				Literal:   node.Literal,
				Content:   node.Content,
				Attribute: node.Attribute,
			},
		}
		var content string
		if !entering {
			content = fmt.Sprintf(figureFmt, "figure", node.Destination)
		}
		newNode.AsLeaf().Literal = []byte(content)
	}

	return hr.markdownRenderer.RenderNode(w, newNode, entering)
}

func (hr *HugoRenderer) RenderHeader(w io.Writer, ast ast.Node) {
	hr.markdownRenderer.RenderHeader(w, ast)
}

func (hr *HugoRenderer) RenderFooter(w io.Writer, ast ast.Node) {
	hr.markdownRenderer.RenderFooter(w, ast)
}
