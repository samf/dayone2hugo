package main

import (
	"bytes"
	"log"

	"github.com/BurntSushi/toml"
)

type FrontMatter struct {
	Date  string   `toml:"date"`
	Title string   `toml:"title"`
	Tags  []string `toml:"tags"`
}

type HugoCmd struct {
	CommonConvert

	KeepTitle bool `help:"keep the '# title' in your markdown"`
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

	err = hc.outBody(body, frontMatter.output())

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
