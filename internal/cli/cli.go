package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/yannickkirschen/scrooge/internal/cli/acc"
	"github.com/yannickkirschen/scrooge/internal/cli/expr"
	"github.com/yannickkirschen/scrooge/internal/cli/tx"
	kong_addon "github.com/yannickkirschen/scrooge/pkg/kong"
	"github.com/yannickkirschen/scrooge/pkg/scrooge"
)

var CLI struct {
	Acc struct {
		Ls struct {
			Filter string `arg:"" optional:""`
		} `cmd:"" help:"List all accounts."`
		New struct {
		} `cmd:"" help:"Add a new account."`
		Edit struct {
			Id string `arg:""`
		} `cmd:"" help:"Edit an account by its Id."`
		Rm struct {
			Id string `arg:""`
		} `cmd:"" help:"Delete an account by its Id."`
		Import struct {
			Path string `arg:""`
		} `cmd:"" help:"Import accounts from a CSV file."`
	} `cmd:"" help:"Work on accounts."`
	Tx struct {
		Ls struct {
			Filter string `arg:"" optional:""`
			Opaque bool
		} `cmd:"" help:"List all transactions."`
		New struct {
		} `cmd:"" help:"Add a new transaction."`
		Edit struct {
			Id uint64 `arg:""`
		} `cmd:"" help:"Edit a transaction by its Id."`
		Rm struct {
			Id uint64 `arg:""`
		} `cmd:"" help:"Delete a transaction by its Id."`
		Import struct {
			Path string `arg:""`
		} `cmd:"" help:"Import transactions from a CSV file."`
	} `cmd:"" help:"Work on transactions."`
	Expr struct {
		New struct {
			Name string `arg:""`
		} `cmd:"" help:"Write a new expr."`
		Run struct {
			Name string `arg:""`
		} `cmd:"" help:"Run an expression."`
	} `cmd:"" help:"Work on expressions."`
}

func InputLoop(ctx *scrooge.Context) error {
	for {
		fmt.Print("scrooge> ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()

		if input == ".exit" {
			return nil
		}

		commands := strings.Split(input, " ")
		kongCtx, err := kong_addon.ParseString(&CLI, commands)
		if err != nil {
			fmt.Printf("error parsing command %s: %s\n", commands, err)
			continue
		}

		if err := HandleCommand(ctx, kongCtx); err != nil {
			fmt.Printf("error executing command %s: %s\n", commands, err)
		}
	}
}

func HandleCommand(scroogeCtx *scrooge.Context, kongCtx *kong.Context) error {
	var err error
	switch kongCtx.Command() {
	case "acc ls":
		err = acc.List(scroogeCtx, "")
	case "acc ls <filter>":
		err = acc.List(scroogeCtx, CLI.Acc.Ls.Filter)
	case "acc new":
		err = acc.New(scroogeCtx)
	case "acc edit <id>":
		err = acc.Edit(scroogeCtx, CLI.Acc.Edit.Id)
	case "acc rm <id>":
		err = scroogeCtx.DeleteAccount(CLI.Acc.Rm.Id)
	case "acc import <path>":
		err = acc.Import(scroogeCtx, CLI.Acc.Import.Path)
	case "tx ls":
		err = tx.List(scroogeCtx, "", false)
	case "tx ls <filter>":
		err = tx.List(scroogeCtx, CLI.Tx.Ls.Filter, CLI.Tx.Ls.Opaque)
	case "tx new":
		err = tx.New(scroogeCtx)
	case "tx edit <id>":
		err = tx.Edit(scroogeCtx, CLI.Tx.Edit.Id)
	case "tx rm <id>":
		err = scroogeCtx.DeleteTransaction(CLI.Tx.Rm.Id)
	case "tx import <path>":
		err = tx.Import(scroogeCtx, CLI.Tx.Import.Path)
	case "expr new <name>":
		err = expr.New(scroogeCtx, CLI.Expr.New.Name)
	case "expr run <name>":
		err = expr.Run(scroogeCtx, CLI.Expr.Run.Name)
	default:
		return fmt.Errorf("unknown command %s", kongCtx.Command())
	}

	return err
}
