package acc

import (
	_ "embed"

	"github.com/yannickkirschen/scrooge/pkg/editor"
	"github.com/yannickkirschen/scrooge/pkg/scrooge"
)

//go:embed new.yaml.tpl
var NewTemplate string

func New(ctx *scrooge.Context) error {
	var acc *scrooge.Account
	return editor.EditTplTempStructFile(NewTemplate, &acc, nil, nil, nil, ctx.SaveAccount)
}
