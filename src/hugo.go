package main

import (
	"bytes"
	"fmt"

	"github.com/BurntSushi/toml"
)

type FrontMatter struct {
	Date  string   `toml:"date"`
	Title string   `toml:"title"`
	Tags  []string `toml:"tags"`
}

type HugoCmd struct {
	CommonConvert
}

func (hc *HugoCmd) Run(ctx *Context) error {
	defer hc.Input.Close()

	doz, _, entry, _, body, err := hc.getStuff(ctx)
	if err != nil {
		return err
	}

	frontMatter := hc.NewFrontMatter(entry, body)
	front, err := frontMatter.output()
	if err != nil {
		return err
	}

	err = hc.GotoOutDir()
	if err != nil {
		return err
	}

	err = hc.outBody(body, front)

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
	front := &FrontMatter{
		Date:  entry.CreationDate,
		Tags:  entry.Tags,
		Title: body.findTitle(),
	}

	return front
}

func (fm *FrontMatter) output() ([]byte, error) {
	var out bytes.Buffer
	out.Write([]byte("+++\n"))
	encoder := toml.NewEncoder(&out)
	err := encoder.Encode(fm)
	if err != nil {
		err = fmt.Errorf("cannot encode frontmatter: %w", err)
		return nil, err
	}
	out.Write([]byte("+++\n"))

	return out.Bytes(), nil
}
