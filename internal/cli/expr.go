package cli

import (
	"github.com/yannickkirschen/scrooge/internal/cli/tx"
	"github.com/yannickkirschen/scrooge/pkg/scrooge"
)

func InitExprEnv(ctx *scrooge.Context) {
	ctx.ExprEnv["lsTx"] = func(ex string) error { return tx.List(ctx, ex) }
}
