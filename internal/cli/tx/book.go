package tx

import (
	"fmt"

	"github.com/yannickkirschen/scrooge/pkg/scrooge"
)

func Book(ctx *scrooge.Context, ids []uint64) error {
	for _, id := range ids {
		tx, ok := ctx.GetTransaction(id)
		if !ok {
			return fmt.Errorf("transaction %d not found", id)
		}

		if tx.Status != scrooge.Booked {
			tx.Status = scrooge.Booked
			if err := ctx.EditTransaction(id, tx); err != nil {
				return fmt.Errorf("error booking transaction %d: %s", id, err)
			}
		}
	}

	return nil
}
