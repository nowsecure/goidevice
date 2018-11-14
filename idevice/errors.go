package idevice

import (
	"errors"
)

// #cgo pkg-config: libimobiledevice-1.0
// #include <libimobiledevice/libimobiledevice.h>
import "C"

var (
	// ErrInvalidArgs .
	ErrInvalidArgs = errors.New("invalid args")
	// ErrUnknown .
	ErrUnknown = errors.New("unknown")
	// ErrNoDevice .
	ErrNoDevice = errors.New("no device")
	// ErrNotEnoughData .
	ErrNotEnoughData = errors.New("not enough data")
	// ErrBadHeader .
	ErrBadHeader = errors.New("bad header")
	// ErrSslError .
	ErrSslError = errors.New("ssl error")
)

func resultToError(result C.idevice_error_t) error {
	switch result {
	case 0:
		return nil
	case -1:
		return ErrInvalidArgs
	case -2:
		return ErrUnknown
	case -3:
		return ErrNoDevice
	case -4:
		return ErrNotEnoughData
	case -5:
		return ErrBadHeader
	case -6:
		return ErrSslError
	default:
		return ErrUnknown
	}
}
