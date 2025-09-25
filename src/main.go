package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/davecgh/go-spew/spew"
)

type Context struct{}

type CLI struct {
	Verbose bool `help:"be verbose"`

	Input *os.File `help:"zip file to be converted" arg:""`
}

func (cli *CLI) Run(ctx *Context) error {
	defer cli.Input.Close()

	doz, err := NewDOZip(cli.Input)
	if err != nil {
		return err
	}

	journal, err := doz.getJournal()
	if err != nil {
		return err
	}

	fmt.Println(spew.Sdump(journal))

	photoWallet := NewPhotoWallet(journal.Entries[0])

	body := NewBody(journal.Entries[0].Text)
	body.fixImages(photoWallet)

	body.render(os.Stdout)

	return nil
}

func main() {
	ctx := kong.Parse(&CLI{})
	err := ctx.Run(&Context{})
	ctx.FatalIfErrorf(err)
}
