package tx

import (
	"bytes"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/shopspring/decimal"
	"github.com/yannickkirschen/scrooge/internal/filter"
	"github.com/yannickkirschen/scrooge/pkg/scrooge"
	"github.com/yannickkirschen/scrooge/pkg/tables"
)

var TimeFormat = "02 Jan 2006"
var TableModel = tables.NewModel(os.Stdout, 1, 1, 2, ' ', tabwriter.AlignRight)

func List(ctx *scrooge.Context, ex string) error {
	headers := []string{"Id", "Date", "Account", "Status", "Amount", "Balance", "Tags"}

	balance := decimal.NewFromInt(0)
	data := [][]string{}
	for _, tx := range ctx.Model.Transactions {
		// Calculate the balance throughout the loop even if the transaction is
		// not being displayed.
		balance = UpdateBalance(balance, *tx.Amount, tx.Type)

		env := ctx.ExprEnv
		env["tx"] = NewTx(tx)

		if matches, err := filter.Matches(env, ex); err != nil {
			return err
		} else if !matches {
			continue
		}

		var tags bytes.Buffer
		for _, tag := range tx.Tags {
			tags.WriteString(tag)
		}

		var amount string
		if tx.Type == scrooge.Spending {
			amount = fmt.Sprintf("-%s %s", tx.Amount.StringFixedBank(2), tx.Currency)
		} else {
			amount = fmt.Sprintf("+%s %s", tx.Amount.String(), tx.Currency)
		}

		data = append(data, []string{
			fmt.Sprintf("%d", tx.Id),
			tx.Date.Format(TimeFormat),
			tx.Account.Id,
			string(tx.Status),
			amount,
			fmt.Sprintf("%4s %s", balance.String(), tx.Currency),
			tags.String(),
		})
	}

	tables.Print(TableModel, headers, data)
	return nil
}

func UpdateBalance(balance, amount decimal.Decimal, typ scrooge.TransactionType) decimal.Decimal {
	switch typ {
	case scrooge.Receipt:
		return balance.Add(amount)
	case scrooge.Spending:
		return balance.Sub(amount)
	case scrooge.Balance:
		return amount
	default:
		return balance
	}
}
