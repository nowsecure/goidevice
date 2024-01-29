package plist

// #cgo pkg-config: libplist-2.0
// #include <stdlib.h>
// #include <plist/plist.h>
import "C"
import (
	"errors"
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

type PList struct {
	P C.plist_t
}

// Create a plist
func Create() *PList {
	return &PList{C.plist_new_dict()}
}

// CreateString .
func CreateString(value string) *PList {
	valueC := C.CString(value)
	defer C.free(unsafe.Pointer(valueC))

	return &PList{C.plist_new_string(valueC)}
}

// CreateArray .
func CreateArray() *PList {
	return &PList{C.plist_new_array()}
}

// Type return the plist type
func (s *PList) Type() int {
	return (int)(C.plist_get_node_type(s.P))
}

// ArraySize .
func (s *PList) ArraySize() int {
	return (int)(C.plist_array_get_size(s.P))
}

// Size size of the plist dict
func (s *PList) Size() int {
	return (int)(C.plist_dict_get_size(s.P))
}

// ArrayItem .
func (s *PList) ArrayItem(index int) *PList {
	return &PList{C.plist_array_get_item(s.P, C.uint(index))}
}

// SetItem .
func (s *PList) SetItem(key string, value PList) {
	keyC := C.CString(key)
	defer C.free(unsafe.Pointer(keyC))

	C.plist_dict_set_item(s.P, keyC, value.P)
}

// GetItemValue is similar to GetItem however it creates
// a copy of the current node only perserving it and its
// childern for traversal
func (s *PList) GetItemValue(key string) (*PList, error) {
	p, err := s.getItem(key)
	if err != nil {
		return nil, err
	}
	pCopy := C.plist_copy(p.P)
	C.plist_free(p.P)
	return &PList{pCopy}, err
}

// GetItem returns the plist/node at the specified key
// maintains parent/child relationship of the plist
// if the key is return nothing nil and an error is returned
func (s *PList) GetItem(key string) (*PList, error) {
	return s.getItem(key)
}

// getItem is a wrapper around plist_get_dict_item
func (s *PList) getItem(key string) (*PList, error) {
	keyC := C.CString(key)
	defer C.free(unsafe.Pointer(keyC))

	p := C.plist_dict_get_item(s.P, keyC)
	if p == nil {
		return &PList{p}, fmt.Errorf("no valid value in plist at key of %s", key)
	}
	return &PList{p}, nil
}

// RemoveItem remove item from plist
func (s *PList) RemoveItem(key string) {
	keyC := C.CString(key)
	defer C.free(unsafe.Pointer(keyC))
	C.plist_dict_remove_item(s.P, keyC)
}

// ArrayAppendItem appent item to array
func (s *PList) ArrayAppendItem(value PList) {
	C.plist_array_append_item(s.P, value.P)
}

// String returns the node as string,
// an empty string if node is not a string.
// (see Type() to check node type)
func (s *PList) String() string {
	var p *C.char = nil
	C.plist_get_string_val(s.P, &p)
	defer C.free(unsafe.Pointer(p))
	return C.GoString(p)
}

// Int returns the node as an int,
// 0 if node is not an int.
// (see Type() to check node type)
func (s *PList) Int() int {
	var p C.uint64_t
	C.plist_get_uint_val(s.P, &p)
	return int(p)
}

// Bool returns the node as a bool
// false if node not an bool
// (see Type() to check node type)
func (s *PList) Bool() bool {
	var p C.uint8_t
	C.plist_get_bool_val(s.P, &p)
	return int(p) == 1
}

func (s *PList) Free() {
	C.plist_free(s.P)
}

func (s *PList) XML() string {
	var buf *C.char = nil
	var num C.uint
	C.plist_to_xml(s.P, &buf, &num)
	return C.GoString(buf)
}

// GetPointer returns the raw internal pointer to the device.
func GetPointer(p PList) unsafe.Pointer {
	return unsafe.Pointer(p.P)
}

// FromPointer returns PList from a raw plist_t pointer
// if the plist is invalid it will return nil
func FromPointer(p unsafe.Pointer) (*PList, error) {
	if p == nil {
		return nil, nil
	}

	plist := &PList{(C.plist_t)(p)}
	if plist.Size() == 0 && plist.ArraySize() == 0 {
		return nil, errors.New("no plist was found for the query")
	}
	return plist, nil
}
