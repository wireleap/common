// Copyright (c) 2021 Wireleap

package withdrawal

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/wireleap/common/api/withdrawalrequest"
)

// field rename
type WR = withdrawalrequest.T

type T struct {
	ID           string `json:"id,omitempty"`
	State        string `json:"state,omitempty"`
	StateChanged int64  `json:"state_changed,omitempty"`

	*WR

	Receipt json.RawMessage `json:"receipt,omitempty"`
}

func (t *T) Validate() error {
	if t.ID == "" {
		return fmt.Errorf("withdrawal id is missing")
	}

	switch t.State {
	case "failed", "pending", "complete":
		// OK
	case "":
		return fmt.Errorf("withdrawal state is missing")
	default:
		return fmt.Errorf("withdrawal state is invalid: `%s`", t.State)
	}

	if t.StateChanged == 0 {
		return errors.New("withdrawal state_changed is missing")
	}

	return nil
}
