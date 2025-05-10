package expr

import (
	"github.com/yannickkirschen/scrooge/pkg/editor"
	"github.com/yannickkirschen/scrooge/pkg/scrooge"
)

func New(ctx *scrooge.Context, name string) error {
	return editor.EditTempStringFile("", func(ex string) error { return ctx.SaveExpression(name, ex) })
}
