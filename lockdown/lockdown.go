package lockdown

// #cgo pkg-config: libimobiledevice-1.0
// #include <stdlib.h>
// #include <libimobiledevice/libimobiledevice.h>
// #include <libimobiledevice/lockdown.h>
// #include <libimobiledevice/service.h>
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
	PList(domain string) (*plist.PList, error)
	Free() error
	GetClient() unsafe.Pointer
	StartService(d idevice.Device, serviceName string) (*Service, error)
	StartServiceClient(d idevice.Device, serviceName string) (*Service, error)
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

func (s *client) PList(domain string) (*plist.PList, error) {
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

	list, err := plist.FromPointer(unsafe.Pointer(node))
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s *client) Free() error {
	return resultToError(C.lockdownd_client_free(s.p))
}

func (s *client) GetClient() unsafe.Pointer {
	return unsafe.Pointer(s.p)
}

// GetDescriptor gets the lockdown descriptor for the service
func (s *Service) GetDescriptor() unsafe.Pointer {
	return unsafe.Pointer(s.descriptor)
}

func (s *Service) FreeDescriptor() {
	C.lockdownd_service_descriptor_free(s.descriptor)
}

func (s *Service) Free() {
	C.lockdownd_service_descriptor_free(s.descriptor)
	C.service_client_free(s.s)
}

func (s *client) StartService(d idevice.Device, serviceName string) (*Service, error) {
	var p C.lockdownd_service_descriptor_t
	svc := C.CString(serviceName)
	defer C.free(unsafe.Pointer(svc))
	err := resultToError(C.lockdownd_start_service(s.p, svc, &p))
	if err != nil {
		return nil, err
	}

	return &Service{descriptor: p}, nil
}

type Service struct {
	s          C.service_client_t
	descriptor C.lockdownd_service_descriptor_t
}

const (
	CRASH_REPORT_MOVER_SERVICE       = "com.apple.crashreportmover"
	CRASH_REPORT_COPY_MOBILE_SERVICE = "com.apple.crashreportcopymobile"
)

func (s *client) StartServiceClient(d idevice.Device, serviceName string) (*Service, error) {
	svc, err := s.StartService(d, serviceName)
	if err != nil {
		return nil, err
	}

	err = serviceResultToError(
		C.service_client_new((C.idevice_t)(idevice.GetPointer(d)),
			svc.descriptor,
			&svc.s,
		),
	)
	return svc, err
}

func (s *Service) ReadPing() error {
	var msg [4]int8
	var n C.uint32_t

	var attempts = 0
	for {
		res := C.service_receive_with_timeout(s.s, (*C.char)(&msg[0]), 4, &n, 2000)
		switch res {
		case 0:
			return nil
		case -7:
			attempts++
			if attempts == 10 {
				return errors.New("failed 10 attempts to ping")
			}
			continue
		default:
			return errors.New(":(((")
		}
	}
}
