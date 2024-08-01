package afc

// #cgo pkg-config: libimobiledevice-1.0
// #include <stdlib.h>
// #include <libimobiledevice/afc.h>
import "C"
import (
	"errors"
	"fmt"
)

func resultToError(result C.afc_error_t) error {
	switch result {
	case 0:
		return nil
	default:
		fmt.Println(result)
		return errors.New("unknown")
	}
}
