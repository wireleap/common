// Copyright (c) 2021 Wireleap

// Package canned provides tools for working with canned HTTP responses.
package canned

import (
	"io/ioutil"
	"net/http"
)

type T struct {
	h    http.Header
	st   int
	body []byte
}

func Can(res *http.Response) (t T, err error) {
	b, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return
	}

	t = T{h: res.Header.Clone(), st: res.StatusCode, body: b}
	return
}

func (t T) Uncan(w http.ResponseWriter) {
	for k, vs := range t.h {
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}

	w.WriteHeader(t.st)
	w.Write(t.body)
}
