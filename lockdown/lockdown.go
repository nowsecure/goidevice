package lockdown

// #cgo pkg-config: libimobiledevice-1.0
// #include <stdlib.h>
// #include <libimobiledevice/lockdown.h>
import "C"
import (
	"errors"
	"unsafe"

	"github.com/nowsecure/goidevice/idevice"
	"github.com/nowsecure/goidevice/plist"
)

// Client is a lockdown client
type Client interface {
	Type() (string, error)
	Pair() error
	ValidatePair() error
	DeviceName() (string, error)
	PList(domain string) (plist.PList, error)
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

// NewClientWithHandshake creates a new lockdown with a handshake
func NewClientWithHandshake(device idevice.Device, label string) (Client, error) {
	labelC := C.CString(label)
	defer C.free(unsafe.Pointer(labelC))

	var p C.lockdownd_client_t
	err := resultToError(C.lockdownd_client_new_with_handshake((C.idevice_t)(idevice.GetPointer(device)), &p, labelC))
	if err != nil {
		return nil, err
	}
	return &client{p}, nil
}

func (s *client) Type() (string, error) {
	var p *C.char = nil
	err := resultToError(C.lockdownd_query_type(s.p, &p))
	defer C.free(unsafe.Pointer(p))
	return C.GoString(p), err
}

func (s *client) Pair() error {
	return resultToError(C.lockdownd_pair(s.p, nil))
}

func (s *client) ValidatePair() error {
	return resultToError(C.lockdownd_validate_pair(s.p, nil))
}

func (s *client) DeviceName() (string, error) {
	var p *C.char
	err := resultToError(C.lockdownd_get_device_name(s.p, &p))
	defer C.free(unsafe.Pointer(p))
	return C.GoString(p), err
}

func (s *client) PList(domain string) (plist.PList, error) {
	var domainC *C.char = nil

	if domain != "" {
		domainC = C.CString(domain)
		defer C.free(unsafe.Pointer(domainC))
	}

	var node C.plist_t
	err := resultToError(C.lockdownd_get_value(s.p, domainC, nil, &node))
	if err != nil {
		return nil, err
	}

	list := plist.FromPointer(unsafe.Pointer(node))
	if list == nil {
		return nil, errors.New("no plist was found for the query")
	}

	return list, nil
}

func (s *client) Close() error {
	err := resultToError(C.lockdownd_client_free(s.p))
	if err == nil {
		s.p = nil
	}
	return err
}
