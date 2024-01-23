package plist

// #cgo pkg-config: libplist-2.0
// #include <stdlib.h>
// #include <plist/plist.h>
import "C"
import (
	"fmt"
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
	Size() int
	ArraySize() int
	ArrayItem(index int) PList
	SetItem(key string, value interface{})
	GetItem(key string) (PList, error)
	GetItemValue(key string) (PList, error)
	Append(value interface{})
	String() string
	Int() int
	Bool() bool
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

// Type return the plist type
func (s *plist) Type() int {
	return (int)(C.plist_get_node_type(s.p))
}

// ArraySize .
func (s *plist) ArraySize() int {
	return (int)(C.plist_array_get_size(s.p))
}

// Size size of the plist dict
func (s *plist) Size() int {
	return (int)(C.plist_dict_get_size(s.p))
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

// GetItemValue is similar to GetItem however it creates
// a copy of the current node only perserving it and its
// childern for traversal
func (s *plist) GetItemValue(key string) (PList, error) {
	p, err := s.getItem(key)
	if err != nil {
		return nil, err
	}
	pCopy := C.plist_copy(p.p)
	C.plist_free(p.p)
	return &plist{pCopy}, err
}

// GetItem returns the plist/node at the specified key
// maintains parent/child relationship of the plist
// if the key is return nothing nil and an error is returned
func (s *plist) GetItem(key string) (PList, error) {
	return s.getItem(key)
}

// getItem is a wrapper around plist_get_dict_item
func (s *plist) getItem(key string) (*plist, error) {
	keyC := C.CString(key)
	defer C.free(unsafe.Pointer(keyC))

	p := C.plist_dict_get_item(s.p, keyC)
	if p == nil {
		return nil, fmt.Errorf("no valid value in plist at key of %s", key)
	}
	return &plist{p}, nil
}

func (s *plist) Append(value interface{}) {
	C.plist_array_append_item(s.p, (C.plist_t)(GetPointer(convertToPList(value))))
}

// String returns the node as string,
// an empty string if node is not a string.
// (see Type() to check node type)
func (s *plist) String() string {
	var p *C.char = nil
	C.plist_get_string_val(s.p, &p)
	defer C.free(unsafe.Pointer(p))
	return C.GoString(p)
}

// Int returns the node as an int,
// 0 if node is not an int.
// (see Type() to check node type)
func (s *plist) Int() int {
	var p C.uint64_t
	C.plist_get_uint_val(s.p, &p)
	return int(p)
}

// Bool returns the node as a bool
// false if node not an bool
// (see Type() to check node type)
func (s *plist) Bool() bool {
	var p C.uint8_t
	C.plist_get_bool_val(s.p, &p)
	return int(p) == 1
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

// FromPointer returns PList from a raw plist_t pointer
// if the plist is invalid it will return nil
func FromPointer(p unsafe.Pointer) PList {
	if p == nil {
		return nil
	}

	plist := &plist{(C.plist_t)(p)}
	if plist.Size() == 0 {
		return nil
	}
	return plist
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
