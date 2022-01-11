package plist

// #cgo pkg-config: libimobiledevice-1.0
// #include <stdlib.h>
// #include <libimobiledevice/lockdown.h>
import "C"
import (
	"unsafe"
)

const (
	// PListTypeBoolean .
	PListTypeBoolean = 0
	// PListTypeUint .
	PListTypeUint = 1
	// PListTypeReal .
	PListTypeReal = 2
	// PListTypeString .
	PListTypeString = 3
	// PListTypeArray .
	PListTypeArray = 4
	// PListTypeDict .
	PListTypeDict = 5
	// PListTypeDate .
	PListTypeDate = 6
	// PListTypeData .
	PListTypeData = 7
	// PListTypeKey .
	PListTypeKey = 8
	// PListTypeUID .
	PListTypeUID = 9
	// PListTypeNone .
	PListTypeNone = 10
)

// PList object
type PList interface {
	Type() int
	ArraySize() int
	ArrayItem(index int) PList
	SetItem(key string, value interface{})
	GetItem(key string) PList
	Append(value interface{})
	String() string
	Free()
}

type plist struct {
	p C.plist_t
}

// Create a plist
func Create() PList {
	p := C.plist_new_dict()
	if p == nil {
		return nil
	}
	return &plist{p}
}

// CreateString .
func CreateString(value string) PList {
	valueC := C.CString(value)
	defer C.free(unsafe.Pointer(valueC))

	p := C.plist_new_string(valueC)
	if p == nil {
		return nil
	}
	return &plist{p}
}

// CreateArray .
func CreateArray() PList {
	p := C.plist_new_array()
	if p == nil {
		return nil
	}
	return &plist{p}
}

// SetItem .
func (s *plist) Type() int {
	return (int)(C.plist_get_node_type(s.p))
}

// ArraySize .
func (s *plist) ArraySize() int {
	return (int)(C.plist_array_get_size(s.p))
}

// ArrayItem .
func (s *plist) ArrayItem(index int) PList {
	p := C.plist_array_get_item(s.p, C.uint(index))
	if p == nil {
		return nil
	}
	return &plist{p}
}

// SetItem .
func (s *plist) SetItem(key string, value interface{}) {
	keyC := C.CString(key)
	defer C.free(unsafe.Pointer(keyC))

	C.plist_dict_set_item(s.p, keyC, (C.plist_t)(GetPointer(convertToPList(value))))
}

// GetItem .
func (s *plist) GetItem(key string) PList {
	keyC := C.CString(key)
	defer C.free(unsafe.Pointer(keyC))

	p := C.plist_dict_get_item(s.p, keyC)
	if p == nil {
		return nil
	}
	return &plist{p}
}

func (s *plist) Append(value interface{}) {
	C.plist_array_append_item(s.p, (C.plist_t)(GetPointer(convertToPList(value))))
}

func (s *plist) String() string {
	var p *C.char
	C.plist_get_string_val(s.p, &p)
	var result string
	if p != nil {
		result = C.GoString(p)
		C.free(unsafe.Pointer(p))
	}
	return result
}

func (s *plist) Free() {
	if s.p != nil {
		C.plist_free(s.p)
		s.p = nil
	}
}

// GetPointer returns the raw internal pointer to the device.
func GetPointer(p PList) unsafe.Pointer {
	internal, _ := p.(*plist)
	return unsafe.Pointer(internal.p)
}

// FromPointer returns a PList from a raw plist_t pointer
func FromPointer(p unsafe.Pointer) PList {
	if p == nil {
		return nil
	}
	return &plist{(C.plist_t)(p)}
}

func convertToPList(val interface{}) PList {
	if str, ok := val.(string); ok {
		return CreateString(str)
	} else if p, ok := val.(PList); ok {
		return p
	} else {
		panic("invalid type")
	}
}
