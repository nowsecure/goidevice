package installation

// #cgo pkg-config: libimobiledevice-1.0
// #include <stdlib.h>
// #include <libimobiledevice/installation_proxy.h>
import "C"
import (
	"unsafe"

	"github.com/pauldotknopf/goidevice/idevice"
	"github.com/pauldotknopf/goidevice/plist"
)

// Proxy is an installation proxy
type Proxy interface {
	Browse(clientOptions plist.PList) (plist.PList, error)
	Close() error
}

type proxy struct {
	p C.instproxy_client_t
}

// NewClientStartService creates a installation proxy
func NewClientStartService(device idevice.Device, label string) (Proxy, error) {
	labelC := C.CString(label)
	defer C.free(unsafe.Pointer(labelC))

	var p C.instproxy_client_t
	err := resultToError(C.instproxy_client_start_service((C.idevice_t)(idevice.GetPointer(device)), &p, labelC))
	if err != nil {
		return nil, err
	}
	return &proxy{p}, nil
}

func (s *proxy) Browse(clientOptions plist.PList) (plist.PList, error) {
	var p C.plist_t
	err := resultToError(C.instproxy_browse(s.p, (C.plist_t)(plist.GetPointer(clientOptions)), &p))
	return plist.FromPointer(unsafe.Pointer(p)), err
}

func (s *proxy) Close() error {
	err := resultToError(C.instproxy_client_free(s.p))
	if err == nil {
		s.p = nil
	}
	return err
}
