// Copyright (c) 2021 Wireleap

package duration

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
	"unicode"
)

// https://github.com/golang/go/issues/25705

type T time.Duration

func (d *T) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)

	if err != nil {
		return err
	}

	last := -1
	var out string

	// convert every instance of days to 24h
	for i, c := range s {
		if unicode.IsDigit(c) {
			if last == -1 {
				last = i
			}
		} else if c == 'd' && last != -1 {
			// days
			n := s[last:i]

			days, err := strconv.Atoi(n)

			if err != nil {
				return errors.New("invalid value of days in duration")
			}

			h := days * 24
			hrs := strconv.Itoa(h)

			out += hrs + "h"
			last = -1
		} else if last != -1 {
			out += s[last : i+1]
			last = -1
		}
	}

	tmpd, err := time.ParseDuration(out)

	if err != nil {
		return err
	}

	*d = T(tmpd)
	return nil
}

func (d T) MarshalJSON() ([]byte, error) {
	// NOTE this encodes days as 24h
	return json.Marshal(time.Duration(d).String())
}

func (d T) String() (r string) {
	dur := time.Duration(d)
	days := int64(dur.Hours()) / 24

	if days >= 1 {
		r = strconv.FormatInt(days, 10) + "d"
		dur -= time.Duration(days*24) * time.Hour

		if dur.Seconds() > 1 {
			r += time.Duration(dur).String()
		}
	} else {
		r = dur.String()
	}

	return
}
