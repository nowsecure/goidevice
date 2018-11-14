package idevice

/*
#cgo pkg-config: libimobiledevice-1.0
#include <libimobiledevice/libimobiledevice.h>
void event_proxy();
static void device_callback(const idevice_event_t *event, void *user_data) {
	event_proxy(event);
}
static idevice_error_t _device_callback(void *user_data) {
	return idevice_event_subscribe(device_callback, user_data);
}
*/
import "C"

import (
	"sync"
	"unsafe"

	"github.com/mattn/go-pointer"
	"github.com/olebedev/emitter"
)

var (
	mutex           sync.Mutex
	emit            *emitter.Emitter
	isSubscribed    bool
	callbackPointer unsafe.Pointer
)

func init() {
	emit = &emitter.Emitter{}
	emit.Use("*", emitter.Void)
}

const (
	// DeviceAdded a device was added
	DeviceAdded = 1
	// DeviceRemoved a device was removed
	DeviceRemoved = 2
	// DevicePaired a device was paired
	DevicePaired = 3
)

type deviceEventInternal struct {
	Event          int32
	UUID           *C.char
	ConnectionType int32
}

// DeviceEvent struct
type DeviceEvent struct {
	Event          int32
	UUID           string
	ConnectionType int32
}

type internalEvent struct {
	event    func(deviceEvent DeviceEvent, userData interface{})
	userData interface{}
}

type callback struct {
}

// Subscribe to event.
func Subscribe() error {
	mutex.Lock()
	defer mutex.Unlock()

	if isSubscribed {
		return nil
	}

	callbackPointer = pointer.Save(&callback{})
	err := resultToError(C._device_callback(callbackPointer))

	if err == nil {
		isSubscribed = true
	} else {
		pointer.Unref(callbackPointer)
		callbackPointer = nil
	}

	return err
}

// Unsubscribe from events.
func Unsubscribe() error {
	mutex.Lock()
	defer mutex.Unlock()

	if !isSubscribed {
		return nil
	}

	err := resultToError(C.idevice_event_unsubscribe())

	if err == nil {
		isSubscribed = false
		pointer.Unref(callbackPointer)
		callbackPointer = nil
	}

	return err
}

// AddEvent Adds an event to be raised raise a device event happens.
func AddEvent() (<-chan DeviceEvent, func()) {
	mutex.Lock()
	defer mutex.Unlock()
	out := make(chan DeviceEvent)
	in := emit.On("event", func(event *emitter.Event) {
		out <- event.Args[0].(DeviceEvent)
	})
	cancel := func() {
		emit.Off("mediaAdded", in)
		close(out)
	}
	return out, cancel
}

//export event_proxy
func event_proxy(deviceEvent unsafe.Pointer) {
	dInternal := (*(*deviceEventInternal)(deviceEvent))
	d := DeviceEvent{}
	d.Event = dInternal.Event
	d.ConnectionType = dInternal.ConnectionType
	d.UUID = C.GoString(dInternal.UUID)
	emit.Emit("event", d)
}
