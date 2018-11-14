package idevice

// #cgo pkg-config: libimobiledevice-1.0
// #include <stdlib.h>
// #include <libimobiledevice/libimobiledevice.h>
import "C"
import "unsafe"

// Device is an iOS device
type Device interface {
	UUID() (string, error)
	Close() error
}

type device struct {
	p C.idevice_t
}

// New up an iOS device object.
func New(uuid string) (Device, error) {
	uuidC := C.CString(uuid)
	defer C.free(unsafe.Pointer(uuidC))
	//C.myprint(cs)
	//var p unsafe.Pointer
	var p C.idevice_t
	err := resultToError(C.idevice_new(&p, uuidC))
	if err != nil {
		return nil, err
	}
	return &device{p}, nil
}

func (s *device) UUID() (string, error) {
	var p *C.char
	err := resultToError(C.idevice_get_udid(s.p, &p))
	var result string
	if p != nil {
		result = C.GoString(p)
		C.free(unsafe.Pointer(p))
	}
	return result, err
}

func (s *device) Close() error {
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
