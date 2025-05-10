package scrooge

import (
	"slices"
	"time"

	"github.com/shopspring/decimal"
	"golang.org/x/text/currency"
)

type Transaction struct {
	Id          uint64            `json:"id" yaml:"id"`
	Date        *time.Time        `json:"date" yaml:"date"`
	Account     *Account          `json:"-" yaml:"-"`
	AccountRef  string            `json:"account" yaml:"account"`
	Type        TransactionType   `json:"type" yaml:"type"`
	Tags        []string          `json:"tags,omitempty" yaml:"tags"`
	Status      TransactionStatus `json:"status" yaml:"status"`
	Description string            `json:"description" yaml:"description"`
	Amount      *decimal.Decimal  `json:"amount" yaml:"amount"`
	Currency    *currency.Unit    `json:"currency" yaml:"currency"`
}

func (tx *Transaction) Equal(o *Transaction) bool {
	return tx.Id == o.Id &&
		tx.Date.Equal(*o.Date) &&
		tx.Type == o.Type &&
		slices.Equal(tx.Tags, o.Tags) &&
		tx.Status == o.Status &&
		tx.Amount.Equal(*o.Amount) &&
		tx.Currency.String() == o.Currency.String()
}

type TransactionType string
type TransactionStatus string

const (
	Receipt  TransactionType = "receipt"
	Spending TransactionType = "spending"
	Balance  TransactionType = "balance"

	Booked  TransactionStatus = "booked"
	Planned TransactionStatus = "planned"
)
