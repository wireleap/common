// Copyright (c) 2021 Wireleap

package relaydir

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/wireleap/common/api/auth"
	"github.com/wireleap/common/api/client"
	"github.com/wireleap/common/api/status"
	"golang.org/x/crypto/bcrypt"
)

// EnrollHandshake performs the first enrollment request given a client and the
// prepared request itself. It completes the challenge-response PoW handshake
// which is required for new enrollments.
func EnrollHandshake(cl *client.Client, req *http.Request) (st *status.T, err error) {
	var (
		res *http.Response
		jm  json.RawMessage
	)

	res, err = cl.PerformRequestNoParse(req)

	if err != nil {
		err = fmt.Errorf(
			"could not perform pre-enrollment request for directory %s: %w",
			req.URL,
			err,
		)
		return
	}

	err = client.Refresh(req)

	if err != nil {
		err = fmt.Errorf(
			"could not refresh pre-enrollment request for directory %s: %w, body='%s'",
			req.URL,
			err,
			string(jm),
		)
		return
	}

	auth.DelHeader(res.Header, auth.API, auth.Version) // ignore API version for this
	err = client.ParseResponse(res, &jm)

	if err != nil {
		err = fmt.Errorf(
			"could not parse pre-enrollment response for directory %s: %w, body='%s'",
			req.URL,
			err,
			string(jm),
		)
		return
	}

	if err = json.Unmarshal(jm, &st); err != nil {
		err = fmt.Errorf(
			"could not unmarshal pre-enrollment response for directory %s: %w, body='%s'",
			req.URL,
			err,
			string(jm),
		)
		return
	}

	if st.Is(status.ErrChallenge) {
		var (
			cost int
			hash []byte

			challenge = auth.GetHeader(res.Header, auth.Directory, auth.Challenge)
		)

		if len(challenge) == 0 {
			err = fmt.Errorf(
				"response requested but no challenge provided by directory %s",
				req.URL,
			)
			return
		}

		cost, err = strconv.Atoi(strings.SplitN(challenge, "~", 2)[0])

		if err != nil {
			err = fmt.Errorf(
				"could not parse pre-enrollment challenge '%s' from directory %s: %w, body='%s'",
				challenge,
				req.URL,
				err,
				string(jm),
			)
			return
		}

		hash, err = bcrypt.GenerateFromPassword([]byte(challenge), cost)

		if err != nil {
			err = fmt.Errorf(
				"could not hash challenge from directory %s: %w, body='%s'",
				req.URL,
				err,
				string(jm),
			)
			return
		}

		auth.SetHeader(req.Header, auth.Directory, auth.Challenge, challenge)
		auth.SetHeader(req.Header, auth.Directory, auth.Response, string(hash))

		jm = nil // erase old body
		err = cl.PerformRequestOnce(req, &jm)

		if err != nil {
			err = fmt.Errorf(
				"could not perform enrollment request for directory %s: %w, body='%s'",
				req.URL,
				err,
				string(jm),
			)
			return
		}

		if err = client.Refresh(req); err != nil {
			err = fmt.Errorf("could not refresh enrollment request for directory %s: %w", req.URL, err)
		}
	}
	return
}
