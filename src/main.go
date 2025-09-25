package main

import (
	"fmt"
	"io"
	"os"

	"github.com/alecthomas/kong"
)

type Context struct{}

type CLI struct {
	Verbose bool `help:"be verbose"`

	Input *os.File `help:"zip file to be converted" arg:""`

	OutDir     string   `help:"directory for output files; will create if necessary" short:"d" type:"path" default:"."`
	Entry      int      `help:"which entry to export" default:"0"`
	FileName   string   `help:"file name of the markdown file we will produce" default:"index.md"`
	PhotoNames []string `help:"friendly names for photos"`
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

	entry, err := journal.getEntry(cli.Entry)
	if err != nil {
		return err
	}

	photoWallet := NewPhotoWallet(entry)
	err = photoWallet.setFriendlyNames(cli.PhotoNames)

	body := NewBody(entry)
	body.fixImages(photoWallet)

	err = cli.GotoOutDir()
	if err != nil {
		return err
	}

	err = cli.outBody(body)
	if err != nil {
		return err
	}

	err = cli.outPhotos(doz, entry)
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

func (cli *CLI) outBody(body *Body) error {
	file, err := os.OpenFile(
		cli.FileName,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0666,
	)
	if err != nil {
		err = fmt.Errorf("cannot open %q: %w", err)
		return err
	}
	defer file.Close()

	body.renderMarkdown(file)

	return nil
}

func (cli *CLI) outPhotos(doz *DOZip, entry *Entry) error {
	for _, p := range entry.Photos {
		zf, err := doz.findPhotoFile(p.MD5)
		if err != nil {
			return err
		}
		infile, err := zf.Open()
		if err != nil {
			err = fmt.Errorf("cannot open zip part for %q: %w",
				p.MD5,
				err,
			)
			return err
		}
		defer infile.Close()

		outname := p.getFName()
		outfile, err := os.OpenFile(
			outname,
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
			0666,
		)
		if err != nil {
			err = fmt.Errorf("cannot open output for photo %q: %w",
				outname,
				err,
			)
			return err
		}
		defer outfile.Close()
		_, err = io.Copy(outfile, infile)
		if err != nil {
			err = fmt.Errorf("problem writing image file %q: %w",
				outname,
				err,
			)
			return err
		}
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
