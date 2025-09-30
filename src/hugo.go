package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/barasher/go-exiftool"
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

	KeepTitle   bool   `help:"keep the '# title' in your markdown"`
	UseFigure   bool   `help:"use a <figure> shortcode for images" env:"HUGO_FIGURE"`
	FigureTag   string `help:"shortcode to use instead of <figure>" default:"figure" env:"HUGO_FIGURE_TAG"`
	LinkToSrc   bool   `help:"in a <figure> shortcode, add a link the same as the src" env:"HUGO_LINK_TO_SRC"`
	SetCaptions bool   `help:"in a <figure> shortcode, add a caption from exif data" env:"HUGO_SET_CAPTIONS"`
}

func (hc *HugoCmd) Run(ctx *Context) error {
	defer hc.Input.Close()

	doz, _, entry, wallet, body, err := hc.getStuff(ctx)
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
		renderer = hc.NewHugoRenderer(wallet)
	}

	err = hc.outPhotos(doz, entry)
	if err != nil {
		return err
	}

	if hc.SetCaptions {
		et, err := exiftool.NewExiftool()
		if err != nil {
			err = fmt.Errorf("cannot launch exiftool: %w", err)
			return err
		}
		defer et.Close()

		err = wallet.setCaptions(et)
		if err != nil {
			return err
		}
	}

	err = hc.outBody(body, frontMatter.output(), renderer)

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
		log.Printf("cannot encode frontmatter: %v", err)
		panic(err)
	}
	out.Write([]byte("+++\n"))

	return out.Bytes()
}

type HugoRenderer struct {
	markdownRenderer markdown.Renderer
	imgTag           string
	linkToSrc        bool
	wallet           *PhotoWallet
	setCaptions      bool
}

func (hc *HugoCmd) NewHugoRenderer(wallet *PhotoWallet) markdown.Renderer {
	return &HugoRenderer{
		markdownRenderer: md.NewRenderer(md.WithRenderInFooter(true)),
		imgTag:           hc.FigureTag,
		linkToSrc:        hc.LinkToSrc,
		wallet:           wallet,
		setCaptions:      hc.SetCaptions,
	}
}

const (
	figureStart = "{{< %s"
	figureKV    = "\n    %v=%q"
	figureEnd   = "\n>}}"
)

func (hr *HugoRenderer) RenderNode(
	w io.Writer,
	node ast.Node,
	entering bool,
) ast.WalkStatus {
	newNode := node
	switch node := node.(type) {
	case *ast.Image:
		if entering {
			return ast.GoToNext
		}

		newNode = &ast.HTMLBlock{
			Leaf: ast.Leaf{
				Parent:    node.Parent,
				Literal:   node.Literal,
				Content:   node.Content,
				Attribute: node.Attribute,
			},
		}
		content := fmt.Sprintf(figureStart, hr.imgTag)
		content += fmt.Sprintf(figureKV, "src", node.Destination)
		if hr.linkToSrc {
			content += fmt.Sprintf(figureKV, "link", node.Destination)
		}
		if hr.setCaptions {
			p := hr.wallet.byFName[string(node.Destination)]
			if p != nil && p.caption != "" {
				content += fmt.Sprintf(figureKV, "caption", p.caption)
			}
		}
		content += figureEnd
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
