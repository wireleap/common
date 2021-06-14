// Copyright (c) 2021 Wireleap

package upgrade

import (
	"bufio"
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/blang/semver"
	"github.com/wireleap/common/cli/fsdir"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/clearsign"
)

type Config struct {
	fm fsdir.T

	q     string
	qpath string

	binfile string
	binpath string

	user, pass  string
	interactive bool

	skipversion *semver.Version
}

func NewConfig(fm fsdir.T, arg0 string, interactive bool) *Config {
	u := &Config{
		fm:          fm,
		q:           "quarantine",
		binfile:     arg0,
		interactive: interactive,
	}
	u.qpath = fm.Path(u.q)
	u.binpath = fm.Path(u.binfile)
	return u
}

const SKIP_FILENAME = ".skip-upgrade-version"

func (u *Config) SkippedVersion() (v *semver.Version) {
	// ignore error -- maybe it's not there
	u.fm.Get(&v, SKIP_FILENAME)
	return
}

func (u *Config) SkipVersion(v semver.Version) error {
	log.Printf("Skipping version %s", v)
	return u.fm.Set(v, SKIP_FILENAME)
}

func (u *Config) GetChangelog(ver semver.Version) (_ string, err error) {
	var (
		chgfile  = u.binfile + ".md"
		qchgpath = u.fm.Path(u.q, chgfile)
		chgurl   = u.ChangelogURL(ver)
	)
	if err = u.Download(qchgpath, chgurl, u.user, u.pass); err != nil {
		return
	}
	b := []byte{}
	if b, err = ioutil.ReadFile(qchgpath); err != nil {
		return
	}
	return string(b), nil
}

func (u *Config) Download(out, url, distuser, distpass string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(path.Dir(out), 0700); err != nil {
		return err
	}
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	req.SetBasicAuth(distuser, distpass)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"%s download request returned code %d: %s",
			url, res.StatusCode, res.Status,
		)
	}
	var r io.Reader
	if u.interactive {
		r = io.TeeReader(res.Body, &WriteCounter{out, 0, res.ContentLength})
	} else {
		fmt.Printf("Downloading %s...\n", out)
		r = res.Body
	}
	if _, err = io.Copy(f, r); err != nil {
		return err
	}
	fmt.Println()
	return nil
}

func (u *Config) GetHash(ver semver.Version) (sha512 []byte, err error) {
	var (
		hshfile  = u.binfile + ".hash"
		hshurl   = u.HashURL(ver)
		qhshpath = u.fm.Path(u.q, hshfile)
	)
	keyring, err := openpgp.ReadArmoredKeyRing(strings.NewReader(BuildKey))
	if err != nil {
		return nil, err
	}
	if err = u.Download(qhshpath, hshurl, u.user, u.pass); err != nil {
		return nil, err
	}
	hf, err := os.Open(qhshpath)
	if err != nil {
		return nil, err
	}
	defer hf.Close()
	msg, err := ioutil.ReadAll(hf)
	if err != nil {
		return nil, err
	}
	block, _ := clearsign.Decode(msg)
	e, err := openpgp.CheckDetachedSignature(keyring, bytes.NewReader(block.Bytes), block.ArmoredSignature.Body)
	if err != nil {
		return nil, err
	}
	ids := []string{}
	for k := range e.Identities {
		ids = append(ids, k)
	}
	fmt.Printf(
		"Good signature from %s %s (%s)\n",
		ids, e.PrimaryKey.KeyIdString(), e.PrimaryKey.CreationTime,
	)
	// update this if hash file format changes
	sha512rex, err := regexp.Compile(fmt.Sprintf(` *([a-f0-9]{128})  %s`, u.binfile))
	if err != nil {
		return nil, err
	}
	_, err = hf.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	matches := sha512rex.FindReaderSubmatchIndex(bufio.NewReader(hf))
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("sha512 sum not found in hash file %s", qhshpath)
	}
	off, end := matches[2], matches[3]
	hex512sum := make([]byte, end-off)
	_, err = hf.ReadAt(hex512sum, int64(off))
	if err != nil {
		return nil, err
	}
	sha512 = make([]byte, hex.DecodedLen(len(hex512sum)))
	_, err = hex.Decode(sha512, hex512sum)
	if err != nil {
		return nil, err
	}
	return sha512, nil
}

func (u *Config) GetBinary(ver semver.Version, binpath string, s512 []byte) (err error) {
	var (
		binurl = u.BinaryURL(ver)
	)
	if err = u.Download(binpath, binurl, u.user, u.pass); err != nil {
		return err
	}
	bf, err := os.Open(binpath)
	if err != nil {
		return err
	}
	defer bf.Close()
	h := sha512.New()
	if _, err = io.Copy(h, bf); err != nil {
		return err
	}
	sha512sum := h.Sum(nil)
	if !bytes.Equal(sha512sum, s512) {
		return fmt.Errorf(
			"hash mismatch: expected %s, got %s",
			hex.EncodeToString(s512),
			hex.EncodeToString(sha512sum),
		)
	}
	return os.Chmod(binpath, 0755)
}

// full upgrade procedure
func (u *Config) Upgrade(ex Executor, v0, v1 semver.Version) (err error) {
	var (
		newbinfile  = u.binfile + ".next"
		newbinpath  = u.fm.Path(newbinfile)
		qnewbinpath = u.fm.Path(u.q, u.binfile)

		s512 []byte
	)
	// create quarantine
	if err = os.MkdirAll(u.qpath, 0700); err != nil {
		return fmt.Errorf("creating quarantine failed: %w", err)
	}
	// get .hash
	if s512, err = u.GetHash(v1); err != nil {
		return fmt.Errorf(
			"downloading or verifying new %s %s .hash file failed: %w",
			u.binfile, v1, err,
		)
	}
	// get & verify binary
	if err = u.GetBinary(v1, qnewbinpath, s512); err != nil {
		return fmt.Errorf(
			"downloading or verifying new %s %s binary failed: %w",
			u.binfile, v1, err,
		)
	}
	// move out of quarantine
	if err = os.Rename(qnewbinpath, newbinpath); err != nil {
		return fmt.Errorf(
			"moving new %s %s binary out of quarantine (%s -> %s) failed: %w",
			u.binfile, v1, qnewbinpath, newbinpath, err,
		)
	}
	return ex(ExecutorArgs{
		Root:   u.fm,
		SrcBin: newbinpath, DstBin: u.binpath,
		SrcVer: v0, DstVer: v1,
	})
}

func (u *Config) Cleanup() { os.RemoveAll(u.qpath) }
