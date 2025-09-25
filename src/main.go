package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/davecgh/go-spew/spew"
)

type Context struct{}

type CLI struct {
	Verbose bool `help:"be verbose"`

	Input *os.File `help:"zip file to be converted" arg:""`

	OutDir   string `help:"directory for output files; will create if necessary" short:"d" type:"path" default:"."`
	Entry    int    `help:"which entry to export" default:"0"`
	FileName string `help:"file name of the markdown file we will produce" default:"index.md"`
}

func (cli *CLI) Run(ctx *Context) error {
	defer cli.Input.Close()

	doz, err := cli.getDOZip()
	if err != nil {
		return err
	}

	journal, err := doz.getJournal()
	if err != nil {
		return err
	}

	fmt.Println(spew.Sdump(journal))

	entry, err := journal.getEntry(cli.Entry)
	if err != nil {
		return err
	}

	photoWallet := NewPhotoWallet(entry)

	body := NewBody(entry)

	body.fixImages(photoWallet)
	body.renderMarkdown(os.Stdout)

	log.Printf("outdir: %q", cli.OutDir)

	err = cli.GotoOutDir()
	if err != nil {
		return err
	}

	return nil
}

func (cli *CLI) GotoOutDir() error {
	err := os.MkdirAll(cli.OutDir, 0777)
	if err != nil {
		err = fmt.Errorf("cannot make directory %q: %w",
			cli.OutDir,
			err,
		)
		return err
	}
	err = os.Chdir(cli.OutDir)
	if err != nil {
		err = fmt.Errorf("cannot access directory %q: %w",
			cli.OutDir,
			err,
		)
		return err
	}

	return nil
}

func (cli *CLI) getDOZip() (*DOZip, error) {
	doz, err := NewDOZip(cli.Input)
	if err != nil {
		return nil, err
	}

	return doz, nil
}

func main() {
	ctx := kong.Parse(&CLI{})
	err := ctx.Run(&Context{})
	ctx.FatalIfErrorf(err)
}
