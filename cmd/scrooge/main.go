package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/yannickkirschen/scrooge/internal/cli"
	"github.com/yannickkirschen/scrooge/internal/filter"
	kong_addon "github.com/yannickkirschen/scrooge/pkg/kong"
	"github.com/yannickkirschen/scrooge/pkg/scrooge"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: scrooge <file>")
		os.Exit(1)
	}

	if os.Args[1] == "--help" {
		kong.Parse(&cli.CLI)
	}

	scroogeCtx, err := scrooge.NewContextFromFile(os.Args[1], cli.TplBaseFunctions, filter.BaseEnv)
	if err != nil {
		panic(err)
	}

	cli.InitExprEnv(scroogeCtx)

	if len(os.Args) == 2 {
		if err := cli.InputLoop(scroogeCtx); err != nil {
			panic(err)
		}

		return
	}

	kongCtx, err := kong_addon.ParseString(&cli.CLI, os.Args[2:])
	if err != nil {
		panic(err)
	}

	if err := cli.HandleCommand(scroogeCtx, kongCtx); err != nil {
		panic(err)
	}
}
