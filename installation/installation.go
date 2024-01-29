package installation

// #cgo pkg-config: libimobiledevice-1.0
// #include <stdlib.h>
// #include <libimobiledevice/installation_proxy.h>
// #include <plist/plist.h>
//
// void instproxy_client_options_add_pair(plist_t client_opts, char * key, char * value)
// {
//		instproxy_client_options_add(client_opts, key, value, NULL);
// }
import "C"
import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/nowsecure/goidevice/idevice"
	"github.com/nowsecure/goidevice/plist"
)

type Proxy struct {
	instproxy_client C.instproxy_client_t
}

// NewClientStartService creates a installation proxy
func NewClientStartService(device idevice.Device, label string) (*Proxy, error) {
	labelC := C.CString(label)
	defer C.free(unsafe.Pointer(labelC))

	var p C.instproxy_client_t
	err := resultToError(C.instproxy_client_start_service((C.idevice_t)(idevice.GetPointer(device)), &p, labelC))
	if err != nil {
		return nil, err
	}
	return &Proxy{p}, nil
}

func (s *Proxy) Browse(opts *ClientOptions) (*plist.PList, error) {
	var apps C.plist_t
	err := resultToError(C.instproxy_browse(s.instproxy_client, (C.plist_t)(&opts.client_opts.P), &apps))
	if err != nil {
		return nil, err
	}

	list, err := plist.FromPointer(unsafe.Pointer(apps))
	if err != nil {
		return nil, err
	}
	if list.Type() != plist.PListTypeArray {
		return nil, errors.New("instproxy_browse returned an invalid plist, must be an array")
	}

	return list, err
}

func (s *Proxy) Free() error {
	return resultToError(C.instproxy_client_free(s.instproxy_client))
}

type AppScope string

const (
	User   AppScope = "User"
	System AppScope = "System"
	All    AppScope = ""
)

type App struct {
	ID      string
	Name    string
	Version string
}

func (a *App) String() string {
	return fmt.Sprintf("Identifier: %s, Name: %s, Version: %s", a.ID, a.Name, a.Version)
}

// Apps returns an array of App
func (s *Proxy) Apps(scope AppScope) ([]App, error) {
	opts := NewClientOptions()
	defer opts.Free()
	if scope != All {
		opts.AddOption("ApplicationType", string(scope))
	}

	appsPList, err := s.Browse(opts)
	if err != nil {
		return nil, err
	}
	defer appsPList.Free()

	apps := []App{}
	for i := 0; i < appsPList.ArraySize(); i++ {
		item := appsPList.ArrayItem(i)
		id, _ := item.GetItem("CFBundleIdentifier")
		name, _ := item.GetItem("CFBundleDisplayName")
		version, _ := item.GetItem("CFBundleShortVersionString")
		apps = append(apps, App{
			ID:      id.String(),
			Name:    name.String(),
			Version: version.String(),
		})
		item.Free()
	}

	return apps, nil
}

type ClientOptions struct {
	client_opts *plist.PList
}

func NewClientOptions() *ClientOptions {
	opts := C.instproxy_client_options_new()
	return &ClientOptions{(*plist.PList)(opts)}
}

func (c *ClientOptions) AddOption(key, value string) {
	keyC := C.CString(key)
	valueC := C.CString(value)
	defer C.free(unsafe.Pointer(keyC))
	defer C.free(unsafe.Pointer(valueC))

	C.instproxy_client_options_add_pair((C.plist_t)(c.client_opts), keyC, valueC)
}

func (c *ClientOptions) RemoveOption(key string) {
	keyC := C.CString(key)
	defer C.free(unsafe.Pointer(keyC))
	c.client_opts.RemoveItem(key)
}

func (c *ClientOptions) Free() {
	c.client_opts.Free()
}
