package acc

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/yannickkirschen/scrooge/pkg/scrooge"
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

	accounts := make([]*scrooge.Account, len(records))
	for i, record := range records {
		if len(record) != 2 {
			return fmt.Errorf("format error in line %d: must be 2 elements, not %d", i, len(record))
		}

		accounts[i] = &scrooge.Account{
			Id:   record[0],
			Name: record[1],
		}
	}

	return ctx.AddAccounts(accounts)
}
