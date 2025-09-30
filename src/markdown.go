package main

import (
	"io"
	"log"
	"net/url"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/md"
	"github.com/gomarkdown/markdown/parser"
)

type Body struct {
	ast ast.Node
}

func NewBody(entry *Entry) *Body {
	input := entry.Text
	extensions := +parser.CommonExtensions
	mdParser := parser.NewWithExtensions(extensions)

	ast := markdown.Parse([]byte(input), mdParser)
	body := &Body{
		ast: ast,
	}

	body.findTitle(false)

	return body
}

func (body *Body) findTitle(deleteTitle bool) string {
	var (
		title     string
		titleNode ast.Node
	)

	trav := func(node ast.Node, entering bool) ast.WalkStatus {
		if !entering {
			return ast.GoToNext
		}

		if title != "" {
			return ast.GoToNext
		}

		if head, ok := node.(*ast.Heading); ok {
			if head.Level == 1 {
				cont := head.AsContainer()
				if len(cont.Children) >= 1 {
					titleNode = node
					child := cont.Children[0]
					title = string(child.AsLeaf().Literal)
				}
			}
		}

		return ast.GoToNext
	}

	ast.WalkFunc(body.ast, trav)

	if deleteTitle && titleNode != nil {
		ast.RemoveFromTree(titleNode)
	}

	return title
}

func (body *Body) fixImages(pw *PhotoWallet) {
	trav := func(node ast.Node, entering bool) ast.WalkStatus {
		if !entering {
			return ast.GoToNext
		}

		if img, ok := node.(*ast.Image); ok {
			dest := string(img.Destination)
			durl, err := url.Parse(dest)
			if err != nil {
				log.Printf("bad destination on image %q: %v", dest, err)
			}

			if durl.Scheme == "dayone-moment" {
				img.Destination = []byte(pw.fixPhotoSrc(durl.Host))
			}
		}
		return ast.GoToNext
	}

	ast.WalkFunc(body.ast, trav)
}

func (body *Body) renderMarkdown(
	out io.Writer,
	mdRender markdown.Renderer,
) {
	if mdRender == nil {
		mdRender = md.NewRenderer(md.WithRenderInFooter(true))
	}
	stuff := markdown.Render(body.ast, mdRender)
	out.Write(stuff)
}
