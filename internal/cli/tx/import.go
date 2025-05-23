package tx

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/shopspring/decimal"
	"github.com/yannickkirschen/scrooge/pkg/scrooge"
	"golang.org/x/text/currency"
)

func Import(ctx *scrooge.Context, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening file %s: %s", path, err)
	}

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading records from file %s: %s", path, err)
	}

	txs := make([]*scrooge.Transaction, len(records))
	for i, record := range records {
		if len(record) < 7 {
			return fmt.Errorf("format error in line %d: must be at least 7 elements, not %d", i, len(record))
		}

		date, err := time.Parse(time.DateOnly, record[0])
		if err != nil {
			return fmt.Errorf("parse error for date %s in line %d: %s", record[0], i, err)
		}

		amount, err := decimal.NewFromString(record[5])
		if err != nil {
			return fmt.Errorf("parse error for amount %s in line %d: %s", record[5], i, err)
		}

		unit, err := currency.ParseISO(record[6])
		if err != nil {
			return fmt.Errorf("parse error for currency %s in line %d: %s", record[6], i, err)
		}

		tx := &scrooge.Transaction{
			Date:        &date,
			AccountRef:  record[1],
			Type:        scrooge.TransactionType(record[2]),
			Tags:        record[7:],
			Status:      scrooge.TransactionStatus(record[3]),
			Description: record[4],
			Amount:      &amount,
			Currency:    &unit,
		}

		ctx.SetIdFor(tx)
		txs[i] = tx
	}

	return ctx.AddTransactions(txs)
}
