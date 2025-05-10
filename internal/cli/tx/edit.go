package tx

import (
	_ "embed"
	"fmt"

	"github.com/yannickkirschen/scrooge/pkg/editor"
	"github.com/yannickkirschen/scrooge/pkg/scrooge"
)

//go:embed edit.yaml.tpl
var EditTemplate string

type editCtx struct {
	newCtx
	Tx *Tx
}

func Edit(ctx *scrooge.Context, id uint64) error {
	tx, ok := ctx.GetTransaction(id)
	if !ok {
		return fmt.Errorf("transaction %d not found", id)
	}

	return editor.EditTplTempStructFile(EditTemplate, &tx, &editCtx{
		newCtx: newCtx{ctx.AccountRefs, ctx.Tags},
		Tx:     NewTx(tx),
	}, nil, ctx.TplFuncMap, ctx.SaveTransaction)
}
