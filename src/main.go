package main

import (
	"fmt"

	"github.com/alecthomas/kong"
)

type Context struct{}

type CLI struct {
	Verbose bool `help:"be verbose"`
}

func (cli *CLI) Run(ctx *Context) error {
	fmt.Println("vim-go")
	return nil
}

func main() {
	var cli CLI
	ctx := kong.Parse(&cli)
	err := ctx.Run(&Context{})
	ctx.FatalIfErrorf(err)
}
