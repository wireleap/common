// Copyright (c) 2022 Wireleap

package accounting

import "fmt"

type T struct {
	Price *int64 `json:"price,omitempty"`
	// Currency is the lowercase ISO code of the currency to use.
	Currency *string `json:"currency,omitempty"`
}

func (a *T) Validate() error {
	switch {
	case a.Price == nil:
		return fmt.Errorf("accounting price is null or missing")
	case a.Currency == nil:
		return fmt.Errorf("accounting currency is null or missing")
	}

	return nil
}
