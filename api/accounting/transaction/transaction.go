// Copyright (c) 2021 Wireleap

package transaction

import (
	"time"
)

// T is the type of a transaction.
type T struct {
	Time time.Time
	Desc string

	Posting []*Posting
}

// Posting is the type of a single posting line reflecting a debit or credit of
// an account.
type Posting struct {
	Account  string
	Amount   int64
	Currency string
}
