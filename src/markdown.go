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

func NewBody(input string) *Body {
	extensions := +parser.CommonExtensions
	mdParser := parser.NewWithExtensions(extensions)

	ast := markdown.Parse([]byte(input), mdParser)
	md := &Body{
		ast: ast,
	}

	return md
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
				log.Printf("bad destination on image %q: %w", dest, err)
			}

			if durl.Scheme == "dayone-moment" {
				img.Destination = []byte(pw.fixPhotoSrc(durl.Host))
			}
		}
		return ast.GoToNext
	}

	ast.WalkFunc(body.ast, trav)
}

func (body *Body) render(out io.Writer) {
	mdRender := md.NewRenderer(md.WithRenderInFooter(true))
	stuff := markdown.Render(body.ast, mdRender)
	out.Write(stuff)
}
