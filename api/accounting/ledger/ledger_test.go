// Copyright (c) 2021 Wireleap

package ledger

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/wireleap/common/api/accounting/transaction"
)

func TestWriteTransaction(t *testing.T) {
	tmpd, err := ioutil.TempDir("", "wltest.*")

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		os.RemoveAll(tmpd)
	})

	tmpf, err := ioutil.TempFile(tmpd, "tally*")

	if err != nil {
		t.Fatal(err)
	}

	tmpf.Close()

	l := T{Filename: tmpf.Name()}
	err = l.WriteTransaction(&transaction.T{})

	if err != nil {
		t.Fatal(err)
	}

	err = l.WriteTransaction(&transaction.T{
		Time: time.Now(),
		Desc: "test",
	})

	if err != nil {
		t.Fatal(err)
	}

	err = l.WriteTransaction(&transaction.T{
		Time: time.Now(),
		Desc: "test",
		Posting: []*transaction.Posting{
			{
				Account:  "foo",
				Amount:   1,
				Currency: "test",
			},
		},
	})

	if err != nil {
		t.Fatal(err)
	}
}
