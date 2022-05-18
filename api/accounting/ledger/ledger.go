// Copyright (c) 2022 Wireleap

package ledger

import (
	"fmt"
	"os"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/wireleap/common/api/accounting/transaction"
)

type T struct {
	Filename string
	Currency string

	mut sync.Mutex
}

func (t *T) WriteTransaction(tr *transaction.T) error {
	t.mut.Lock()
	defer t.mut.Unlock()

	// date description
	//     acct amount
	//     acct amount
	//     ...
	//     acct amount
	f, err := os.OpenFile(t.Filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		return err
	}

	defer f.Close()

	w := tabwriter.NewWriter(f, 8, 4, 2, ' ', 0)

	if tr.Time.IsZero() {
		tr.Time = time.Now()
	}

	_, err = fmt.Fprintf(
		w,
		"%s %s \t; @%d\t\n",
		tr.Time.Format("2006-01-02"),
		tr.Desc,
		tr.Time.Unix(),
	)

	if err != nil {
		return err
	}

	for _, p := range tr.Posting {
		bucks := p.Amount / 100
		cents := p.Amount % 100

		sign := " "

		if p.Amount < 0 {
			sign = "-"
			bucks = -bucks
			cents = -cents
		}

		if len(p.Currency) == 0 {
			p.Currency = t.Currency
		}

		if len(p.Comment) > 0 {
			p.Comment = " ; " + p.Comment
		}

		_, err = fmt.Fprintf(
			w,
			"    %s\t %s%d.%02d %s%s\n",
			p.Account,
			sign,
			bucks,
			cents,
			p.Currency,
			p.Comment,
		)

		if err != nil {
			return err
		}
	}

	fmt.Fprintln(w)
	w.Flush()
	f.Sync()

	return nil
}
