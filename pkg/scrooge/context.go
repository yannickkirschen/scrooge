package scrooge

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"text/template"

	"github.com/expr-lang/expr"
)

var Perm os.FileMode = 0644

type Context struct {
	path string

	Model       *Model
	AccountRefs []string
	Tags        []string
	NextId      uint64
	TplFuncMap  template.FuncMap
	ExprEnv     map[string]any
}

func NewContext(model *Model, tplFuncMap template.FuncMap, exprEnv map[string]any) (*Context, error) {
	ctx := &Context{
		"",
		model,
		[]string{},
		[]string{},
		0,
		tplFuncMap,
		exprEnv,
	}

	if err := ctx.UpdateRefs(); err != nil {
		return nil, fmt.Errorf("parsing error while updating references: %s", err)
	}

	return ctx, nil
}

func NewContextFromFile(path string, tplFuncMap template.FuncMap, exprEnv map[string]any) (*Context, error) {
	if err := InitFile(path); err != nil {
		return nil, fmt.Errorf("error initializing file %s: %s", path, err)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening %s: %s", path, err)
	}

	var model *Model
	if err := json.NewDecoder(f).Decode(&model); err != nil {
		return nil, fmt.Errorf("error decoding json from %s into *scrooge.Model: %s", path, err)
	}

	ctx, err := NewContext(model, tplFuncMap, exprEnv)
	if err != nil {
		return nil, fmt.Errorf("error creating context from file %s: %s", path, err)
	}

	ctx.path = path
	return ctx, nil
}

func InitFile(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, Perm)
		if err != nil {
			return fmt.Errorf("error creating file %s: %s", path, err)
		}
		defer f.Close()

		if err := json.NewEncoder(f).Encode(&Model{}); err != nil {
			return fmt.Errorf("error parsing empty model to %s: %s", path, err)
		}
	}

	return nil
}

func (ctx *Context) UpdateRefs() error {
	for _, tx := range ctx.Model.Transactions {
		if err := ctx.UpdateAccount(tx); err != nil {
			return err
		}

		ctx.UpdateAccountRef()
		ctx.UpdateTags(tx)
		ctx.UpdateNextId(tx)
	}

	return nil
}

func (ctx *Context) UpdateAccount(tx *Transaction) error {
	acc, ok := ctx.GetAccount(tx.AccountRef)
	if !ok {
		return fmt.Errorf("account %s in transaction %d not found", tx.AccountRef, tx.Id)
	}

	tx.Account = acc
	return nil
}

func (ctx *Context) UpdateAccountRef() {
	for _, acc := range ctx.Model.Accounts {
		if !slices.Contains(ctx.AccountRefs, acc.Id) {
			ctx.AccountRefs = append(ctx.AccountRefs, acc.Id)
		}
	}
}

func (ctx *Context) UpdateTags(tx *Transaction) {
	for _, tag := range tx.Tags {
		if tag != "" && !slices.Contains(ctx.Tags, tag) {
			ctx.Tags = append(ctx.Tags, tag)
		}
	}
}

func (ctx *Context) UpdateNextId(tx *Transaction) {
	if tx.Id == ctx.NextId {
		ctx.NextId++
	} else if tx.Id > ctx.NextId {
		ctx.NextId += tx.Id + 1
	}
}

func (ctx *Context) SetIdFor(tx *Transaction) {
	tx.Id = ctx.NextId
	ctx.NextId++
}

func (ctx *Context) GetAccount(id string) (*Account, bool) {
	acc, _, ok := ctx.Model.GetAccount(id)
	return acc, ok
}

func (ctx *Context) GetTransaction(id uint64) (*Transaction, bool) {
	tx, _, ok := ctx.Model.GetTransaction(id)
	return tx, ok
}

func (ctx *Context) SaveToFile() error {
	if err := ctx.UpdateRefs(); err != nil {
		return fmt.Errorf("error updating refs: %s", err)
	}

	f, err := os.OpenFile(ctx.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("cannot open file %s for writing: %s", ctx.path, err)
	}

	if err := json.NewEncoder(f).Encode(ctx.Model); err != nil {
		return fmt.Errorf("cannot encode into JSON file %s: %s", ctx.path, err)
	}

	return nil
}

func (ctx *Context) AddAccounts(accounts []*Account) error {
	ctx.Model.Accounts = append(ctx.Model.Accounts, accounts...)
	return ctx.SaveToFile()
}

func (ctx *Context) AddAccount(new *Account) error {
	if acc, ok := ctx.GetAccount(new.Id); !ok {
		ctx.Model.Accounts = append(ctx.Model.Accounts, new)
	} else {
		id := acc.Id
		acc = new
		acc.Id = id
	}

	return ctx.SaveToFile()
}

func (ctx *Context) EditAccount(id string, acc *Account) error {
	_, i, ok := ctx.Model.GetAccount(acc.Id)
	if !ok {
		return fmt.Errorf("account %s not found", acc.Id)
	}

	acc.Id = id
	ctx.Model.Accounts[i] = acc
	return ctx.SaveToFile()
}

func (ctx *Context) AddTransactions(txs []*Transaction) error {
	ctx.Model.Transactions = append(ctx.Model.Transactions, txs...)
	return ctx.SaveToFile()
}

func (ctx *Context) AddTransaction(tx *Transaction) error {
	ctx.SetIdFor(tx)
	ctx.Model.Transactions = append(ctx.Model.Transactions, tx)
	return ctx.SaveToFile()
}

func (ctx *Context) EditTransaction(id uint64, tx *Transaction) error {
	_, i, ok := ctx.Model.GetTransaction(tx.Id)
	if !ok {
		return fmt.Errorf("transaction %d not found", tx.Id)
	}

	tx.Id = id
	ctx.Model.Transactions[i] = tx
	return ctx.SaveToFile()
}

func (ctx *Context) SaveExpression(name string, ex string) error {
	_, err := expr.Compile(ex, expr.Env(ctx.ExprEnv))
	if err != nil {
		return fmt.Errorf("compilation of expression %s failed: %s", ex, err)
	}

	ctx.Model.Expressions[name] = ex
	return ctx.SaveToFile()
}

func (ctx *Context) DeleteAccount(id string) error {
	_, i, ok := ctx.Model.GetAccount(id)
	if !ok {
		return fmt.Errorf("account %s not found", id)
	}

	txId, found := ctx.HasAccountRef(id)
	if found {
		return fmt.Errorf("cannot delete account %s as it is referenced in transaction %d", id, txId)
	}

	ctx.Model.Accounts = slices.Delete(ctx.Model.Accounts, i, i+1)
	return ctx.SaveToFile()
}

func (ctx *Context) HasAccountRef(id string) (uint64, bool) {
	for _, tx := range ctx.Model.Transactions {
		if tx.AccountRef == id {
			return tx.Id, true
		}
	}

	return 0, false
}

func (ctx *Context) DeleteTransaction(id uint64) error {
	_, i, ok := ctx.Model.GetTransaction(id)
	if !ok {
		return fmt.Errorf("transaction %d not found", id)
	}

	ctx.Model.Transactions = slices.Delete(ctx.Model.Transactions, i, i+1)
	return ctx.SaveToFile()
}
