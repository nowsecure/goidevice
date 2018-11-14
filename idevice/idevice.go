package idevice

/*
#cgo pkg-config: libimobiledevice-1.0
#include <stdlib.h>
#include <stdio.h>
#include <libimobiledevice/libimobiledevice.h>
void event_proxy();
static void device_callback(const idevice_event_t *event, void *user_data) {
	printf("test    %d\n", event->event);
	event_proxy(event);
}
static void _device_callback(void *user_data) {
	printf("registering callback\n");
	idevice_event_subscribe(device_callback, user_data);
}
*/
import "C"

import (
	"sync"
	"unsafe"

	"github.com/mattn/go-pointer"
)

var (
	mutex           sync.Mutex
	events          []*internalEvent
	isSubscribed    bool
	callbackPointer unsafe.Pointer
)

const (
	// DeviceAdd a device was added
	DeviceAdd = 1
	// DeviceRemove a device was removed
	DeviceRemove = 2
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
func Subscribe() {
	mutex.Lock()
	defer mutex.Unlock()

	if isSubscribed {
		return
	}

	callbackPointer = pointer.Save(&callback{})
	C._device_callback(callbackPointer)
}

// Unsubscribe from events.
func Unsubscribe() {
	mutex.Lock()
	defer mutex.Unlock()

	if !isSubscribed {
		return
	}

	C.idevice_event_unsubscribe()

	pointer.Unref(callbackPointer)
}

// AddEvent Adds an event to be raised raise a device event happens.
func AddEvent(callback func(deviceEvent DeviceEvent, userData interface{}), userData interface{}) (func(), error) {
	mutex.Lock()
	defer mutex.Unlock()
	newEvent := &internalEvent{callback, userData}
	events = append(events, newEvent)
	return func() {
		mutex.Lock()
		defer mutex.Unlock()
		for eventIndex, event := range events {
			if event == newEvent {
				events = append(events[:eventIndex], events[eventIndex+1:]...)
				return
			}
		}
	}, nil
}

//export event_proxy
func event_proxy(deviceEvent unsafe.Pointer) {
	dInternal := (*(*deviceEventInternal)(deviceEvent))
	d := DeviceEvent{}
	d.Event = dInternal.Event
	d.ConnectionType = dInternal.ConnectionType
	d.UUID = C.GoString(dInternal.UUID)
	for _, event := range events {
		event.event(d, event.userData)
	}
}
