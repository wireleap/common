// Copyright (c) 2022 Wireleap

package balance

import (
	"encoding/json"
	"errors"
	"math/big"
	"sync"
)

type T struct{ i *internal }

type internal struct {
	sync.Mutex

	Currency string   `json:"currency"`
	Value    *big.Rat `json:"value"`
	Pending  int64    `json:"pending,omitempty"`
}

func (t *T) MarshalJSON() ([]byte, error) { return json.Marshal(t.i) }

func (t *T) UnmarshalJSON(b []byte) error {
	*t = T{&internal{}}
	return json.Unmarshal(b, t.i)
}

func New(cur string) *T {
	return &T{i: &internal{
		Currency: cur,
		Value:    big.NewRat(0, 1),
		Pending:  0,
	}}
}

func (t *T) Add(x *big.Rat) {
	t.i.Lock()
	t.i.Value.Add(t.i.Value, x)
	t.i.Unlock()
}

func (t *T) Book(delta int64) error {
	t.i.Lock()

	if t.i.Pending != 0 {
		t.i.Unlock()
		return errors.New("balance already has pending transaction")
	}

	d := big.NewRat(delta, 1)
	d.Abs(d)

	if t.i.Value.Cmp(d) == -1 {
		t.i.Unlock()
		return errors.New("insufficient balance available")
	}

	t.i.Pending = delta
	t.i.Unlock()
	return nil
}

func (t *T) Commit() (n int64) {
	t.i.Lock()
	t.i.Value.Add(t.i.Value, big.NewRat(t.i.Pending, 1))
	n = t.i.Pending
	t.i.Pending = 0
	t.i.Unlock()
	return
}

func (t *T) Cancel() {
	t.i.Lock()
	t.i.Pending = 0
	t.i.Unlock()
}

type Exported struct {
	Currency  string `json:"currency"`
	Available int64  `json:"available"`
	Pending   int64  `json:"pending,omitempty"`
}

func ratFloor(x *big.Rat) int64 {
	r := &big.Int{}
	r.Div(x.Num(), x.Denom())
	return r.Int64()
}

func (t *T) Export() (r Exported) {
	t.i.Lock()
	r.Currency = t.i.Currency
	r.Available = ratFloor(t.i.Value)
	r.Pending = t.i.Pending
	t.i.Unlock()
	return
}
