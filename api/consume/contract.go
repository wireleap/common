// Copyright (c) 2021 Wireleap

package consume

import (
	"crypto/ed25519"
	"fmt"
	"net/http"

	"github.com/wireleap/common/api/auth"
	"github.com/wireleap/common/api/client"
	"github.com/wireleap/common/api/contractinfo"
	"github.com/wireleap/common/api/dirinfo"
	"github.com/wireleap/common/api/relaylist"
	"github.com/wireleap/common/api/texturl"
)

func ContractInfo(cl *client.Client, sc *texturl.URL) (info *contractinfo.T, err error) {
	infourl := sc.String() + "/info"
	if err = cl.Perform(http.MethodGet, infourl, nil, &info); err != nil {
		err = fmt.Errorf("could not get contract info from %s: %s", infourl, err)
	}
	return
}

// returns pubkey of this sc
func ContractPubkey(cl *client.Client, sc *texturl.URL) (ed25519.PublicKey, error) {
	info, err := ContractInfo(cl, sc)
	if err != nil {
		return nil, err
	}
	return info.Pubkey.T(), nil
}

func DirectoryData(cl *client.Client, sc *texturl.URL) (ddata *contractinfo.Directory, err error) {
	info, err := ContractInfo(cl, sc)
	if err != nil {
		return nil, err
	}
	return &info.Directory, nil
}

func DirectoryInfo(cl *client.Client, sc *texturl.URL) (dinfo dirinfo.T, err error) {
	ddata, err := DirectoryData(cl, sc)
	if err != nil {
		return
	}
	dinfourl := ddata.Endpoint.String() + "/info"
	if err = cl.Perform(http.MethodGet, dinfourl, nil, &dinfo); err != nil {
		err = fmt.Errorf("could not get directory info from %s: %s", dinfourl, err)
	}
	return
}

// returns relays of this sc's directory
func ContractRelays(cl *client.Client, sc *texturl.URL) (rl relaylist.T, err error) {
	ddata, err := DirectoryData(cl, sc)
	if err != nil {
		return
	}
	dinfo, err := DirectoryInfo(cl, sc)
	if err != nil {
		return
	}
	if dinfo.PublicKey.String() != ddata.PublicKey.String() {
		return nil, fmt.Errorf(
			"directory public key does not match: contract says %s, directory says %s",
			ddata.PublicKey.String(), dinfo.PublicKey.String(),
		)
	}
	dirurl := dinfo.Endpoint.String() + "/relays"
	if err = cl.Perform(http.MethodGet, dirurl, nil, &rl, auth.Directory); err != nil {
		err = fmt.Errorf("could not perform request towards %s: %w", dirurl, err)
	}
	return
}
