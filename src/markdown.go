package main

import (
	"log"
	"net/url"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

type Body struct {
	ast ast.Node
}

func NewMD(input string) *Body {
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
				log.Printf("mapping %q to %q",
					dest,
					pw.fixPhotoSrc(durl.Host),
				)
				img.Destination = []byte(pw.fixPhotoSrc(durl.Host))
			}
		}
		return ast.GoToNext
	}

	ast.WalkFunc(body.ast, trav)
}
