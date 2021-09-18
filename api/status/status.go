// Copyright (c) 2021 Wireleap

package status

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// T is an error-compatible type for reporting errors to the API client
// in a consistent way.
type T struct {
	Code   int    `json:"code"`
	Desc   string `json:"description"`
	Origin string `json:"origin,omitempty"`
	Cause  Cause  `json:"cause,omitempty"`
}

var (
	OK = &T{
		Code: http.StatusOK,
		Desc: "OK",
	}

	ErrNotFound = &T{
		Code: http.StatusNotFound,
		Desc: "resource not found",
	}

	ErrRequest = &T{
		Code: http.StatusBadRequest,
		Desc: "bad request",
	}

	ErrMethod = &T{
		Code: http.StatusMethodNotAllowed,
		Desc: "HTTP method not allowed",
	}

	ErrInternal = &T{
		Code: http.StatusInternalServerError,
		Desc: "internal server error",
	}

	ErrUnpaid = &T{
		Code: http.StatusPaymentRequired,
		Desc: "action requires payment",
	}

	ErrForbidden = &T{
		Code: http.StatusForbidden,
		Desc: "action not allowed",
	}

	ErrConflict = &T{
		Code: http.StatusConflict,
		Desc: "action conflicts with existing resource",
	}

	ErrGateway = &T{
		Code: http.StatusBadGateway,
		Desc: "gateway is unreachable or down",
	}

	ErrChallenge = &T{
		Code: http.StatusAccepted,
		Desc: "repeat request with response to challenge in headers",
	}

	ErrUpgrade = &T{
		Code: http.StatusAccepted,
		Desc: "please upgrade your relay to the given version as soon as possible",
	}
)

func (t *T) Is(maybe error) bool {
	t2, ok := maybe.(*T)
	return ok && t.Code == t2.Code && t.Desc == t2.Desc
}

// TODO increase granularity
func IsRetryable(maybe error) (is bool) {
	t, ok := maybe.(*T)
	if !ok || (ok && (t.Is(ErrInternal) || t.Is(ErrGateway))) {
		is = true
	}
	return
}

func (t *T) Unwrap() error { return t.Cause }

func (t *T) Wrap(cause error) *T {
	tcopy := *t
	tcopy.Cause = Cause(cause.Error())
	return &tcopy
}

func (t *T) Error() string {
	b, err := json.Marshal(t)

	if err != nil {
		panic("should never happen")
	}

	return string(b)
}

type Cause string

func (c Cause) Error() string { return string(c) }

const (
	CauseInsufficientBalance      Cause = "insufficient balance to complete request"
	CauseExpiredPof               Cause = "expired proof of funding"
	CauseSneakyPof                Cause = "this pof has been seen/used already"
	CauseInvalidSig               Cause = "invalid signature"
	CauseMissingParam             Cause = "missing parameter in URL"
	CauseNoText                   Cause = "format=text is not supported for this field"
	CauseUnknownFormat            Cause = "unknown format"
	CauseVersionMismatch          Cause = "major version mismatch"
	CauseRequestExpired           Cause = "request expired"
	CauseSTRejected               Cause = "sharetoken submission rejected"
	CauseWithdrawalPending        Cause = "a withdrawal is already pending"
	CauseWithdrawalInvalid        Cause = "withdrawal amount is invalid (<= 0)"
	CauseSettlementNotOpen        Cause = "settlement window not yet open"
	CauseSettlementClosed         Cause = "settlement window already closed"
	CauseContractPubkeyMismatch   Cause = "contract public key mismatch"
	CausePaymentSystemUnreachable Cause = "payment system is unreachable or down, please try again later"
	CauseBadEnrollmentKey         Cause = "enrollment key is incorrect"
)

var (
	ErrInsufficientBalance      = ErrForbidden.Wrap(CauseInsufficientBalance)
	ErrExpiredPof               = ErrRequest.Wrap(CauseExpiredPof)
	ErrSneakyPof                = ErrRequest.Wrap(CauseSneakyPof)
	ErrInvalidSig               = ErrRequest.Wrap(CauseInvalidSig)
	ErrMissingParam             = ErrRequest.Wrap(CauseMissingParam)
	ErrNoText                   = ErrRequest.Wrap(CauseNoText)
	ErrUnknownFormat            = ErrRequest.Wrap(CauseUnknownFormat)
	ErrVersionMismatch          = ErrRequest.Wrap(CauseVersionMismatch)
	ErrRequestExpired           = ErrRequest.Wrap(CauseRequestExpired)
	ErrSTRejected               = ErrRequest.Wrap(CauseSTRejected)
	ErrWithdrawalPending        = ErrConflict.Wrap(CauseWithdrawalPending)
	ErrWithdrawalInvalid        = ErrRequest.Wrap(CauseWithdrawalInvalid)
	ErrSettlementNotOpen        = ErrRequest.Wrap(CauseSettlementNotOpen)
	ErrSettlementClosed         = ErrRequest.Wrap(CauseSettlementClosed)
	ErrContractPubkeyMismatch   = ErrRequest.Wrap(CauseContractPubkeyMismatch)
	ErrPaymentSystemUnreachable = ErrGateway.Wrap(CausePaymentSystemUnreachable)
	ErrBadEnrollmentKey         = ErrRequest.Wrap(CauseBadEnrollmentKey)
)

func IsCircuitError(maybe error) bool {
	var s *T
	return errors.As(maybe, &s) && s.Origin != "target" && s.Code != http.StatusBadRequest
}

func (t *T) WriteTo(w io.Writer) (int, error) {
	b, err := json.Marshal(t)

	if err != nil {
		return 0, err
	}

	rw, ok := w.(http.ResponseWriter)

	if ok {
		rw.Header().Set("Content-Type", "application/json")
		http.Error(rw, string(b), t.Code)
		return len(b), nil
	}

	return w.Write(b)
}

func (t *T) ToHeader(h http.Header) { h.Set("wl-status", t.Error()) }
