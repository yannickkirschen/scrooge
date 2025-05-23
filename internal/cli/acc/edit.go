package acc

import (
	_ "embed"
	"fmt"

	"github.com/yannickkirschen/scrooge/pkg/editor"
	"github.com/yannickkirschen/scrooge/pkg/scrooge"
)

//go:embed edit.yaml.tpl
var EditTemplate string

type editCtx struct {
	Acc *scrooge.Account
}

func Edit(ctx *scrooge.Context, id string) error {
	acc, ok := ctx.GetAccount(id)
	if !ok {
		return fmt.Errorf("account %s not found", id)
	}

	return editor.EditTplTempStructFile(EditTemplate, &acc, &editCtx{acc}, nil, ctx.TplFuncMap, func(acc *scrooge.Account) error { return ctx.EditAccount(id, acc) })
}
