package expr

import (
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/yannickkirschen/scrooge/pkg/scrooge"
)

func Run(ctx *scrooge.Context, name string) error {
	ex, ok := ctx.Model.Expressions[name]
	if !ok {
		return fmt.Errorf("expr %s not found", name)
	}

	program, err := expr.Compile(ex, expr.Env(ctx.ExprEnv))
	if err != nil {
		return fmt.Errorf("compilation of expression %s failed: %s", ex, err)
	}

	output, err := expr.Run(program, ctx.ExprEnv)
	if err != nil {
		return fmt.Errorf("expression failed: %s", err)
	}

	if output != nil {
		fmt.Println(output)
	}

	return nil
}
