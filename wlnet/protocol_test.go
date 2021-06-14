// Copyright (c) 2021 Wireleap

package wlnet

import (
	"bytes"
	"net"
	"reflect"
	"testing"
)

func TestCodeHeader(t *testing.T) {
	i := encodeHeader(true, 1)

	if i != -1 {
		t.Error("incorrect sign for oob encoded header")
	}

	oob, size := decodeHeader(i)

	if oob != true {
		t.Error("incorrect oob value for decoded header")
	}

	if size != 1 {
		t.Error("incorrect size for decoded header")
	}

	i = encodeHeader(false, 1)

	if i != 1 {
		t.Error("incorrect sign for non-oob header of size 1")
	}

	oob, size = decodeHeader(i)

	if oob != false {
		t.Error("incorrect oob value for decoded header")
	}

	if size != 1 {
		t.Error("incorrect size for decoded header")
	}
}

func TestRWHeader(t *testing.T) {
	var buf bytes.Buffer

	woob, wsize := true, 1
	err := writeHeader(&buf, woob, wsize)

	if err != nil {
		t.Fatal(err)
	}

	oob, size, err := readHeader(&buf)

	if err != nil {
		t.Fatal(err)
	}

	if oob != woob {
		t.Error("read different value for oob than written")
	}

	if size != wsize {
		t.Error("read different value for size than written")
	}
}

func TestFragConn(t *testing.T) {
	c1, c2 := net.Pipe()

	frc := &FragReadConn{Conn: c1, Errf: func(e error) { t.Error(e) }}
	fwc := &FragWriteCloser{c2}

	b1 := []byte{'h', 'e', 'l', 'l', 'o'}

	t.Run("write", func(t *testing.T) {
		t.Parallel()

		n, err := fwc.Write(b1)

		if err != nil {
			t.Error(err)
		}

		if n != len(b1) {
			t.Errorf("only %d out of %d bytes written", n, len(b1))
		}
	})

	t.Run("read", func(t *testing.T) {
		t.Parallel()

		b2 := make([]byte, 5)
		n, err := frc.Read(b2)

		if err != nil {
			t.Error(err)
		}

		if n != len(b1) {
			t.Errorf("only %d out of %d bytes read", n, len(b1))
		}

		if !reflect.DeepEqual(b1, b2) {
			t.Errorf("data mismatch; wrote %v, read %v", b1, b2)
		}
	})
}
