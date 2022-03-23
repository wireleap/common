// Copyright (c) 2022 Wireleap

package pofrequest

import (
	"fmt"

	"github.com/wireleap/common/api/accounting"
)

type T struct {
	Quantity   *int64        `json:"quantity,omitempty"`
	Type       *string       `json:"type,omitempty"`
	Duration   *int64        `json:"duration,omitempty"`
	Accounting *accounting.T `json:"accounting,omitempty"`
}

func (r *T) Validate() error {
	switch {
	case r == nil:
		return fmt.Errorf("pof request is null or missing")
	case r.Quantity == nil:
		return fmt.Errorf("pof quantity is null or missing")
	case r.Type == nil:
		return fmt.Errorf("pof type is null or missing")
	case r.Duration == nil:
		return fmt.Errorf("pof duration is null or missing")
	}

	return nil
}
