package scrooge

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"golang.org/x/text/currency"
	"gopkg.in/yaml.v3"
)

func (tx *Transaction) MarshalJSON() ([]byte, error) {
	return tx.Marshal(json.Marshal)
}

func (tx *Transaction) MarshalYAML() (any, error) {
	return tx.Marshal(yaml.Marshal)
}

func (tx *Transaction) Marshal(f func(any) ([]byte, error)) ([]byte, error) {
	type alias Transaction
	return f(&struct {
		*alias
		AccountRef string `json:"account"`
		Date       string `json:"date"`
		Currency   string `json:"currency"`
	}{
		alias:      (*alias)(tx),
		AccountRef: tx.Account.Id,
		Date:       tx.Date.Format(time.DateOnly),
		Currency:   tx.Currency.String(),
	})
}

func (tx *Transaction) unmarshalId(v any) error {
	id, err := strconv.ParseUint(fmt.Sprintf("%v", v), 10, 64)
	if err != nil {
		return fmt.Errorf("cannot unmarshal string into Go struct field Transaction.Id of type uint64")
	}
	tx.Id = id
	return nil
}

func (tx *Transaction) unmarshalDate(v any) error {
	t, err := time.Parse(time.DateOnly, fmt.Sprintf("%v", v))
	if err != nil {
		return fmt.Errorf("cannot unmarshal string into Go struct field Transaction.Date of type *time.Time: %s", err)
	}
	tx.Date = &t
	return nil
}

func (tx *Transaction) unmarshalAccountRef(v any) error {
	tx.AccountRef = fmt.Sprintf("%v", v)
	return nil
}

func (tx *Transaction) unmarshalType(v any) error {
	switch fmt.Sprintf("%v", v) {
	case string(Receipt):
		tx.Type = Receipt
	case string(Spending):
		tx.Type = Spending
	case string(Balance):
		tx.Type = Balance
	default:
		return fmt.Errorf("cannot unmarshal string into Go struct field Transaction.Type of type scrooge.TransactionType: type %s is unknown", v)
	}

	return nil
}

func (tx *Transaction) unmarshalTags(v any) error {
	if v == nil {
		return nil
	}

	tags, ok := v.([]any)
	if !ok {
		return fmt.Errorf("cannot unmarshal string into Go struct field Transaction.Tags of type []string: %s", v)
	}

	tx.Tags = []string{}
	for _, tag := range tags {
		tx.Tags = append(tx.Tags, fmt.Sprintf("%v", tag))
	}

	return nil
}

func (tx *Transaction) unmarshalDescription(v any) error {
	tx.Description = fmt.Sprintf("%v", v)
	return nil
}

func (tx *Transaction) unmarshalStatus(v any) error {
	switch fmt.Sprintf("%v", v) {
	case string(Booked):
		tx.Status = Booked
	case string(Planned):
		tx.Status = Planned
	default:
		return fmt.Errorf("cannot unmarshal string into Go struct field Transaction.Status of type scrooge.TransactionStatus: %s", v)
	}

	return nil
}

func (tx *Transaction) unmarshalAmount(v any) error {
	amount, err := decimal.NewFromString(fmt.Sprintf("%v", v))
	if err != nil {
		return fmt.Errorf("cannot unmarshal string into Go struct field Transaction.Amount of type *decimal.Decimal: %s", err)
	}

	if amount.IsNegative() {
		return fmt.Errorf("validation error for amount %s: amount must not be negative", amount)
	}

	tx.Amount = &amount
	return nil
}

func (tx *Transaction) unmarshalCurrency(v any) error {
	currency, err := currency.ParseISO(fmt.Sprintf("%v", v))
	if err != nil {
		return fmt.Errorf("cannot unmarshal string into Go struct field Transaction.Currency of type *currency.Unit: %s", err)
	}
	tx.Currency = &currency
	return nil
}

func (tx *Transaction) UnmarshalJSON(data []byte) error {
	raw := map[string]any{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("cannot unmarshal byte array into type map[string]string: %s", err)
	}
	return tx.Unmarshal(raw, json.Unmarshal)
}

func (tx *Transaction) UnmarshalYAML(value *yaml.Node) error {
	raw := map[string]any{}
	if err := value.Decode(&raw); err != nil {
		return fmt.Errorf("cannot decode YAML node into type map[string]string: %s", err)
	}
	return tx.Unmarshal(raw, yaml.Unmarshal)
}

func (tx *Transaction) Unmarshal(raw map[string]any, f func([]byte, any) error) error {
	var unmarshalFuncs = map[string]func(any) error{
		"id":          tx.unmarshalId,
		"date":        tx.unmarshalDate,
		"account":     tx.unmarshalAccountRef,
		"type":        tx.unmarshalType,
		"tags":        tx.unmarshalTags,
		"status":      tx.unmarshalStatus,
		"description": tx.unmarshalDescription,
		"amount":      tx.unmarshalAmount,
		"currency":    tx.unmarshalCurrency,
	}

	for k, v := range raw {
		f, ok := unmarshalFuncs[k]
		if !ok {
			return fmt.Errorf("cannot unmarshal value of key %s into type scrooge.Transaction: key unknown", k)
		}

		if err := f(v); err != nil {
			return err
		}
	}

	return nil
}
