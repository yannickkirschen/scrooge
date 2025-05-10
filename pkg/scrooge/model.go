package scrooge

type Model struct {
	Accounts     []*Account        `json:"accounts"`
	Expressions  map[string]string `json:"expressions"`
	Transactions []*Transaction    `json:"transactions"`
}

func (model *Model) GetAccount(id string) (*Account, int, bool) {
	for i, account := range model.Accounts {
		if account.Id == id {
			return account, i, true
		}
	}

	return nil, -1, false
}

func (model *Model) GetTransaction(id uint64) (*Transaction, int, bool) {
	for i, tx := range model.Transactions {
		if tx.Id == id {
			return tx, i, true
		}
	}

	return nil, -1, false
}
