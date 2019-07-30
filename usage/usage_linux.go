package usage

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"sync/atomic"
	"unsafe"
)

// Interface represents a network interface present on the system.
type Interface struct {
	name string
	path string

	prevIn, prevOut uint64
}

// OpenInterface opens the interface with the specified name.
func OpenInterface(name string) (*Interface, error) {
	i := &Interface{
		name: name,
		path: filepath.Join("/sys/class/net/", name),
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

func (i *Interface) read(name string, prev *uint64) (uint64, error) {
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

// In returns the number of bytes received by the interface since the
// last time that it was called, or since the interface was opened if
// it hasn't been called before.
func (i *Interface) In() (uint64, error) {
	return i.read("rx_bytes", &i.prevIn)
}

// Out returns the number of bytes sent by the interface since the
// last time that it was called, or since the interface was opened if
// it hasn't been called before.
func (i *Interface) Out() (uint64, error) {
	return i.read("tx_bytes", &i.prevOut)
}
