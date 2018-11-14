package lockdown

import (
	"errors"
)

// #cgo pkg-config: libimobiledevice-1.0
// #include <libimobiledevice/lockdown.h>
import "C"

var (
	// ErrUnknown .
	ErrUnknown = errors.New("unknown")
	// TODO: Add all the other error types
)

func resultToError(result C.lockdownd_error_t) error {
	switch result {
	case 0:
		return nil
	default:
		return ErrUnknown
	}
}
