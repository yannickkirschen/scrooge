package scrooge_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/yannickkirschen/scrooge/pkg/scrooge"
	"golang.org/x/text/currency"
)

var validTransaction *scrooge.Transaction

func init() {
	ts, _ := time.Parse(time.RFC3339, "2025-04-25T20:00:00.000+02:00")
	amount, _ := decimal.NewFromString("0.00")
	currency, _ := currency.ParseISO("EUR")

	validTransaction = &scrooge.Transaction{
		Id:       1,
		Date:     &ts,
		Type:     scrooge.Spending,
		Tags:     []string{"fixed-cost"},
		Status:   "booked",
		Amount:   &amount,
		Currency: &currency,
	}
}

func parse(filename string) (*scrooge.Transaction, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	var tx *scrooge.Transaction
	if err := json.NewDecoder(f).Decode(&tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func assertUnmarshalOk(filename string, t *testing.T) {
	if _, err := parse(filename); err != nil {
		t.Fatalf("didn't expected an error but got one: %s", err)
	}
}

func assertUnmarshalNotOk(filename string, t *testing.T) {
	if _, err := parse(filename); err == nil {
		t.Fatal("expected error, but didn't get one")
	}
}

func TestUnmarshalValid(t *testing.T) {
	assertUnmarshalOk("../../test/valid.json", t)
}

func TestUnmarshalEmptyTags(t *testing.T) {
	assertUnmarshalOk("../../test/empty-tags.json", t)
}

func TestUnmarshalNegativeId(t *testing.T) {
	assertUnmarshalNotOk("../../test/negative-id.json", t)
}

func TestUnmarshalNoIsoDate(t *testing.T) {
	assertUnmarshalNotOk("../../test/no-iso-date.json", t)
}

func TestUnmarshalUnknownType(t *testing.T) {
	assertUnmarshalNotOk("../../test/unknown-type.json", t)
}

func TestUnmarshalNegativeAmount(t *testing.T) {
	assertUnmarshalNotOk("../../test/negative-amount.json", t)
}

func TestUnmarshalUnknownCurrency(t *testing.T) {
	assertUnmarshalNotOk("../../test/unknown-currency.json", t)
}

func TestContent(t *testing.T) {
	tx, _ := parse("../../test/valid.json")
	if !tx.Equal(validTransaction) {
		t.Fatalf("expected transaction to be %v but got %v", validTransaction, tx)
	}
}
