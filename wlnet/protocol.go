// Copyright (c) 2021 Wireleap

package wlnet

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"

	"github.com/wireleap/common/api/interfaces/clientrelay"
	"github.com/wireleap/common/api/sharetoken"
	"github.com/wireleap/common/api/status"
	"github.com/wireleap/common/api/texturl"

	"github.com/blang/semver"
)

// MAGIC is the Wireleap magic number.
const MAGIC uint8 = 42

// Init is the struct type encoding values passed while initializing the
// tunneled connection ("init payload").
type Init struct {
	Command  string          `json:"command"`
	Protocol string          `json:"protocol"`
	Remote   *texturl.URL    `json:"remote"`
	Token    *sharetoken.T   `json:"token"`
	Version  *semver.Version `json:"version"`
}

// WriteTo serializes the init payload to some io.Writer.
func (i *Init) WriteTo(w io.Writer) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(i)

	if err != nil {
		return err
	}

	err = writeHeader(w, true, buf.Len())

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

// ReadInit reads an init payload from a provided io.Reader.
func ReadInit(r io.Reader) (*Init, error) {
	oob, size, err := readHeader(r)

	if err != nil {
		return nil, err
	}

	if oob != true {
		return nil, fmt.Errorf("non-OOB packet where init payload expected")
	}

	p := new(Init)
	buf := make([]byte, size)
	_, err = io.ReadFull(r, buf)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buf, p)

	if err != nil {
		return nil, err
	}

	err = p.sanityCheck()

	if err != nil {
		return nil, err
	}

	return p, nil
}

// encodeHeader encodes the size of the following message and whether it
// contains out-of-band data. A negative size does not make sense so we can use
// the sign bit for marking OOB data via negative-signed size.
func encodeHeader(oob bool, size int) int64 {
	if oob {
		return int64(-size)
	} else {
		return int64(size)
	}
}

// decodeHeader decodes a int64 header into its sign bit (negative meaning
// OOB data) and number part (number of bytes to read following the header).
func decodeHeader(h int64) (bool, int) {
	if h < int64(0) {
		return true, int(-h)
	} else {
		return false, int(h)
	}
}

// readHeader reads the OOB bit and size of the following packet from an
// io.Reader or returns an error.
func readHeader(r io.Reader) (bool, int, error) {
	var m uint8
	var h int64

	// magic number
	if err := binary.Read(r, binary.LittleEndian, &m); err != nil {
		// we have no way of knowing how many bytes were read by binary.Read()!
		// https://codereview.appspot.com/3762041/
		// https://github.com/golang/go/issues/18585
		return false, 0, err
	}

	if m != MAGIC {
		// assume 1 byte was read by the previous binary.Read() call
		return false, 0, fmt.Errorf("no magic found, please contact thaumaturgist")
	}

	// OOB bit and packet size
	if err := binary.Read(r, binary.LittleEndian, &h); err != nil {
		// assume 1 byte was read by the previous binary.Read() call
		return false, 0, err
	}

	oob, size := decodeHeader(h)
	return oob, size, nil
}

// writeHeader writes the given OOB bit and size encoded as a packet header
// into the provided io.Writer.
func writeHeader(w io.Writer, oob bool, size int) error {
	err := binary.Write(w, binary.LittleEndian, MAGIC)

	if err != nil {
		return err
	}

	h := encodeHeader(oob, size)
	err = binary.Write(w, binary.LittleEndian, h)

	if err != nil {
		return err
	}

	return nil
}

// FragReadConn

// FragReadConn is the type of a connection reading from a network stream of
// wireleap-relay-fragmented data transparently.
type FragReadConn struct {
	net.Conn
	left int
	Errf func(error)
}

// readError(r, size) reads an status.T of size size from io.Reader r.
func readError(r io.Reader, size int) (*status.T, error) {
	e := new(status.T)
	buf := make([]byte, size)
	_, err := io.ReadFull(r, buf)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buf, e)

	if err != nil {
		return nil, err
	}

	return e, nil
}

// Read overrides the Read method of a net.Conn wrapped by a FragReadConn.
func (c *FragReadConn) Read(p []byte) (int, error) {
	var chunk int
	plen := len(p)

	if c.left > 0 {
		if c.left < plen {
			chunk = c.left
			c.left = 0
		} else { // c.left >= plen
			chunk = plen
			c.left -= plen
		}
	} else {
		oob, size, err := readHeader(c.Conn)

		if err != nil {
			return 0, err
		}

		if oob {
			je, err := readError(c.Conn, size)

			if err != nil {
				return 0, err
			}

			if c.Errf != nil {
				c.Errf(je)
			}

			return 0, je
		}

		if size > plen {
			chunk = plen
			c.left = size - plen
		} else {
			chunk = size
		}
	}

	return io.ReadFull(c.Conn, p[:chunk])
}

// FragWriteCloser

// FragWriteCloser is the type of a connection writing data to a network stream
// of wireleap-relay-fragmented data transparently.
type FragWriteCloser struct{ io.ReadWriteCloser }

// Write overrides the Write method of a net.Conn wrapped by FragWriteCloser.
func (c FragWriteCloser) Write(p []byte) (int, error) {
	err := writeHeader(c.ReadWriteCloser, false, len(p))

	if err != nil {
		return 0, err
	}

	err = binary.Write(c.ReadWriteCloser, binary.LittleEndian, p)

	if err != nil {
		// assume 1 byte magic + 8 bytes header were written by previous
		// binary.Write()
		return 9, err
	}

	return len(p), nil
}

// WriteStatus(w, err) serializes the error err to the io.Writer w.
func WriteStatus(w io.Writer, s *status.T) error {
	var buf bytes.Buffer
	_, err := s.WriteTo(&buf)

	if err != nil {
		return err
	}

	err = writeHeader(w, true, buf.Len())

	if err != nil {
		return err
	}

	_, err = buf.WriteTo(w)

	if err != nil {
		return err
	}

	return nil
}
