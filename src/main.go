package main

import (
	"archive/zip"
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

	fi, err := cli.Input.Stat()
	if err != nil {
		err = fmt.Errorf("cannot stat input file: %w", err)
		return err
	}

	zippy, err := zip.NewReader(cli.Input, fi.Size())
	if err != nil {
		err = fmt.Errorf("cannot open zip reader: %w", err)
		return err
	}

	dir := map[string]*zip.File{}
	for _, f := range zippy.File {
		dir[f.Name] = f
	}

	jname := "Journal.json"
	journalFile := dir[jname]
	if journalFile == nil {
		return fmt.Errorf("no %v file found", jname)
	}
	journal, err := parseDOJson(journalFile)
	if err != nil {
		err = fmt.Errorf("cannot parse %v: %w", jname, err)
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
