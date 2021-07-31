// Copyright (c) 2021 Wireleap

package interfaces

import "github.com/blang/semver"

type Component string

const (
	Client   Component = "client"
	Relay              = "relay"
	Dir                = "dir"
	Contract           = "contract"
	Auth               = "auth"
	PS                 = "ps"
)

func (c Component) String() string { return string(c) }

type T struct {
	Consumer Component      `json:"consumer"`
	Provider Component      `json:"provider"`
	Version  semver.Version `json:"version"`
}

func (t T) String() string { return string(t.Consumer) + "-" + string(t.Provider) }
