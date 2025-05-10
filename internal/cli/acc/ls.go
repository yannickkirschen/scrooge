package acc

import (
	"os"
	"text/tabwriter"

	"github.com/yannickkirschen/scrooge/internal/filter"
	"github.com/yannickkirschen/scrooge/pkg/scrooge"
	"github.com/yannickkirschen/scrooge/pkg/tables"
)

var TableModel = tables.NewModel(os.Stdout, 1, 1, 2, ' ', tabwriter.AlignRight)

func List(ctx *scrooge.Context, ex string) error {
	headers := []string{"Id", "Name"}

	data := [][]string{}
	for _, acc := range ctx.Model.Accounts {
		env := ctx.ExprEnv
		env["acc"] = acc

		if matches, err := filter.Matches(env, ex); err != nil {
			return err
		} else if !matches {
			continue
		}

		data = append(data, []string{
			acc.Id,
			acc.Name,
		})
	}

	tables.Print(TableModel, headers, data)
	return nil
}
