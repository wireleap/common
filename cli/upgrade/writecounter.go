// Copyright (c) 2021 Wireleap

package upgrade

import "fmt"

type WriteCounter struct {
	url, name string
	n, total  int64
}

func (wc *WriteCounter) Write(p []byte) (n int, err error) {
	n = len(p)
	wc.n += int64(n)
	perc := int64(float64(wc.n) / float64(wc.total) * 100)
	fmt.Printf("Downloading %s to %s... %d%%\r", wc.url, wc.name, perc)
	return
}
