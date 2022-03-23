// Copyright (c) 2022 Wireleap

package withdrawalrequest

import "fmt"

type T struct {
	Amount      int64  `json:"amount,omitempty"`
	Type        string `json:"type,omitempty"`
	Destination string `json:"destination,omitempty"`
}

func (t *T) Validate() error {
	if t.Amount <= 0 {
		return fmt.Errorf("withdrawal request amount is invalid (<= 0): %d", t.Amount)
	}

	if t.Type == "" {
		return fmt.Errorf("withdrawal request type is missing")
	}

	if t.Destination == "" {
		return fmt.Errorf("withdrawal destination is missing")
	}

	return nil
}
