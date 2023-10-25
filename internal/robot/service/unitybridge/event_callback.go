package unitybridge

import "C"
import (
	"fmt"
	"log"
	"reflect"
	"sync"
)

var (
	m                          sync.Mutex
	eventCodeIEventHandlersMap = make(map[DJIUnityEventType]map[uintptr]IEventHandler)
)

func registerIEventHandler(eventType DJIUnityEventType, handler IEventHandler) {
	m.Lock()
	defer m.Unlock()
	if handler == nil {
		return
	} else {
		handlers, ok := eventCodeIEventHandlersMap[eventType]
		if !ok {
			handlers = make(map[uintptr]IEventHandler)
			eventCodeIEventHandlersMap[eventType] = handlers
		}
		handlers[getInterfaceValuePointer(handler)] = handler
	}
}

func unregisterEventHandler(eventHandler IEventHandler) {
	m.Lock()
	defer m.Unlock()
	for eventCode, handlers := range eventCodeIEventHandlersMap {
		delete(handlers, getInterfaceValuePointer(eventHandler))
		if len(handlers) == 0 {
			delete(eventCodeIEventHandlersMap, eventCode)
		}
	}
}

func runEventCallback(eventCode uint64, data []byte, tag uint64) error {
	event := NewDJIUnityEvent(eventCode)

	m.Lock()
	defer m.Unlock()
	handlers, ok := eventCodeIEventHandlersMap[event.Type()]
	if !ok {
		return fmt.Errorf("no handlers for event type %d", event.Type())
	}
	for _, handler := range handlers {
		go handler.OnEventCallback(event, data, tag)
	}
	return nil
}

//export eventCallbackGo
func eventCallbackGo(eventCode uint64, data []byte, tag uint64) {
	err := runEventCallback(eventCode, data, tag)
	if err != nil {
		log.Printf("error running event callback: %s\n", err)
	}
}

func pointToSameAddress(a, b interface{}) bool {
	return getInterfaceValuePointer(a) == getInterfaceValuePointer(b)
}

func getInterfaceValuePointer(i interface{}) uintptr {
	return reflect.ValueOf(i).Pointer()
}
