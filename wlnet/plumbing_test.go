// Copyright (c) 2022 Wireleap

package wlnet

import (
	"bytes"
	"context"
	"net"
	"testing"
	"time"
)

var test = []byte{'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'}

const bufsize int = 2048

func TestRetransmit(t *testing.T) {
	ec := make(chan error)

	r := bytes.NewReader(test)

	var w bytes.Buffer
	go retransmit(r, &w, ec, bufsize)

	e := <-ec

	if e != nil {
		t.Error(e)
	}
}

func TestSplice(t *testing.T) {
	c1, c2 := net.Pipe()
	ctx := context.Background()

	go Splice(ctx, c1, c2, time.Second*0, bufsize)

	_, err := c1.Write(test)

	if err != nil {
		t.Error(err)
	}

	_, err = c2.Write(test)

	if err != nil {
		t.Error(err)
	}

	b := make([]byte, len(test))

	_, err = c1.Read(b)

	if err != nil {
		t.Error(err)
	}

	n := bytes.Compare(b, test)

	if n != 0 {
		t.Errorf("bytes.Compare returned %d for b vs test", n)
	}

	_, err = c2.Read(b)

	if err != nil {
		t.Error(err)
	}

	n = bytes.Compare(b, test)

	if n != 0 {
		t.Errorf("bytes.Compare returned %d for b vs test", n)
	}
}
