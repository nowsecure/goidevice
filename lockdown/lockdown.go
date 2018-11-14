package lockdown

// #cgo pkg-config: libimobiledevice-1.0
// #include <stdlib.h>
// #include <libimobiledevice/lockdown.h>
import "C"

import (
	"unsafe"

	"github.com/pauldotknopf/goidevice/idevice"
)

// Client is a lockdown client
type Client interface {
	Type() (string, error)
	Close() error
}

type client struct {
	p C.lockdownd_client_t
}

// NewClient creates a new lockdown client
func NewClient(device idevice.Device, label string) (Client, error) {
	labelC := C.CString(label)
	defer C.free(unsafe.Pointer(labelC))

	var p C.lockdownd_client_t
	err := resultToError(C.lockdownd_client_new((C.idevice_t)(idevice.GetPointer(device)), &p, labelC))
	if err != nil {
		return nil, err
	}
	return &client{p}, nil
}

func (s *client) Type() (string, error) {
	var p *C.char
	err := resultToError(C.lockdownd_query_type(s.p, &p))
	var result string
	if p != nil {
		result = C.GoString(p)
		C.free(unsafe.Pointer(p))
	}
	return result, err
}

func (s *client) Close() error {
	err := resultToError(C.lockdownd_client_free(s.p))
	if err == nil {
		s.p = nil
	}
	return err
}
