package tx

import (
	_ "embed"

	"github.com/yannickkirschen/scrooge/pkg/editor"
	"github.com/yannickkirschen/scrooge/pkg/scrooge"
)

//go:embed new.yaml.tpl
var NewTemplate string

type newCtx struct {
	AccountRefs []string
	Tags        []string
}

func New(ctx *scrooge.Context) error {
	var tx *scrooge.Transaction
	return editor.EditTplTempStructFile(NewTemplate, &tx, &newCtx{ctx.AccountRefs, ctx.Tags}, nil, ctx.TplFuncMap, ctx.AddTransaction)
}
