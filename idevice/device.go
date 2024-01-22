package idevice

// #cgo pkg-config: libimobiledevice-1.0
// #include <stdlib.h>
// #include <libimobiledevice/libimobiledevice.h>
import "C"
import "unsafe"

// Device is an iOS device
type Device interface {
	UDID() (string, error)
	Free() error
}

type device struct {
	p C.idevice_t
}

// New up an iOS device object.
func New(udid string) (Device, error) {
	udidC := C.CString(udid)
	defer C.free(unsafe.Pointer(udidC))

	var p C.idevice_t
	err := resultToError(C.idevice_new(&p, udidC))
	if err != nil {
		return nil, err
	}
	return &device{p}, nil
}

func (s *device) UDID() (string, error) {
	var p *C.char = nil
	err := resultToError(C.idevice_get_udid(s.p, &p))
	defer C.free(unsafe.Pointer(p))
	return C.GoString(p), err
}

func (s *device) Free() error {
	err := resultToError(C.idevice_free(s.p))
	if err == nil {
		s.p = nil
	}
	return err
}

// GetPointer returns the raw internal pointer to the device.
func GetPointer(d Device) unsafe.Pointer {
	internal, _ := d.(*device)
	return unsafe.Pointer(internal.p)
}
