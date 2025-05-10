package tx

import (
	"time"

	"github.com/yannickkirschen/scrooge/pkg/scrooge"
)

type Tx struct {
	Id          uint64
	Date        string
	Account     *scrooge.Account
	Type        string
	Tags        []string
	Status      string
	Description string
	Amount      float64
	Currency    string
}

func NewTx(tx *scrooge.Transaction) *Tx {
	amount, _ := tx.Amount.Float64()
	return &Tx{
		Id:          tx.Id,
		Date:        tx.Date.Format(time.DateOnly),
		Account:     tx.Account,
		Type:        string(tx.Type),
		Tags:        tx.Tags,
		Status:      string(tx.Status),
		Description: tx.Description,
		Amount:      amount,
		Currency:    tx.Currency.String(),
	}
}
