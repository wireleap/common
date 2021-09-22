// Copyright (c) 2021 Wireleap

package wlnet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/wireleap/common/api/interfaces/clientrelay"
	"github.com/wireleap/common/api/sharetoken"
	"github.com/wireleap/common/api/texturl"

	"github.com/blang/semver"
)

const PayloadHeader = "wl-payload"

// Init is the struct type encoding values passed while initializing the
// tunneled connection ("init payload").
type Init struct {
	Command  string          `json:"command"`
	Protocol string          `json:"protocol,omitempty"`
	Remote   *texturl.URL    `json:"remote,omitempty"`
	Token    *sharetoken.T   `json:"token,omitempty"`
	Version  *semver.Version `json:"version,omitempty"`
}

func (i *Init) Headers() map[string]string {
	b, err := json.Marshal(i)
	if err != nil {
		log.Printf("error when marshaling payload %+v: %s", i, err)
		return map[string]string{}
	}

	return map[string]string{PayloadHeader: string(b)}
}

func InitFromHeaders(h http.Header) (i *Init, err error) {
	i = &Init{}
	if err = json.Unmarshal([]byte(h.Get(PayloadHeader)), &i); err != nil {
		return nil, fmt.Errorf("could not parse payload header: %w", err)
	}
	return
}

// WriteTo serializes the init payload to some io.Writer.
func (i *Init) WriteTo(w io.Writer) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(i)

	if err != nil {
		return err
	}

	_, err = buf.WriteTo(w)

	if err != nil {
		return err
	}

	return nil
}

// sanityCheck performs sanity checks on an init payload.
func (i *Init) sanityCheck() error {
	if i.Command == "PING" {
		return nil
	}

	switch i.Protocol {
	case "tcp", "tcp4", "tcp6", "udp", "udp4", "udp6":
		// OK
	default:
		return fmt.Errorf("unknown protocol in payload: %s", i.Protocol)
	}

	switch i.Remote.Scheme {
	case "wireleap", "https", "target":
		// OK
	default:
		return fmt.Errorf("unknown URL scheme in payload: %s", i.Remote.Scheme)
	}

	switch i.Command {
	case "CONNECT":
		// OK
	default:
		return fmt.Errorf("unknown command in payload: %s", i.Command)
	}

	if clientrelay.T.Version.Minor != i.Version.Minor {
		return fmt.Errorf("expecting version 0.%d.x, got %s", clientrelay.T.Version.Minor, i.Version)
	}

	return nil
}
