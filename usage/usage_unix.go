// +build linux freebsd

package usage

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"sync/atomic"
	"unsafe"
)

type inter struct {
	name string
	path string

	prevIn, prevOut uint64
}

func openInterface(name string) (*Interface, error) {
	i := &Interface{
		inter{
			name: name,
			path: filepath.Join("/sys/class/net/", name),
		},
	}

	_, err := i.In()
	if err != nil {
		return nil, err
	}

	_, err = i.Out()
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (i *inter) read(name string, prev *uint64) (uint64, error) {
	raw, err := ioutil.ReadFile(filepath.Join(i.path, "statistics", name))
	if err != nil {
		return 0, err
	}
	raw = bytes.TrimSpace(raw)

	n, err := strconv.ParseUint(*(*string)(unsafe.Pointer(&raw)), 10, 0)
	if err != nil {
		return 0, err
	}

	p := atomic.SwapUint64(prev, n)
	return n - p, nil
}

func (i *inter) In() (uint64, error) {
	return i.read("rx_bytes", &i.prevIn)
}

func (i *inter) Out() (uint64, error) {
	return i.read("tx_bytes", &i.prevOut)
}
