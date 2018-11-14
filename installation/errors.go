package installation

import (
	"errors"
)

// #cgo pkg-config: libimobiledevice-1.0
// #include <stdlib.h>
// #include <libimobiledevice/installation_proxy.h>
import "C"

var (
	// ErrUnknown .
	ErrUnknown = errors.New("unknown")
	// TODO: Add the rest of the errors
)

func resultToError(result C.instproxy_error_t) error {
	switch result {
	case 0:
		return nil
	default:
		return ErrUnknown
	}
}
