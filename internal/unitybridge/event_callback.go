package unitybridge

import "C"
import (
	"fmt"
	"log"
	"sync"
)

var (
	m                       sync.Mutex
	eventCallbackHandlerMap = make(map[uint64]EventCallbackHandler)
)

func setEventCallbackHandler(eventCode uint64, handler EventCallbackHandler) {
	m.Lock()
	defer m.Unlock()
	if handler == nil {
		delete(eventCallbackHandlerMap, eventCode)
	} else {
		eventCallbackHandlerMap[eventCode] = handler
	}
}

func runEventCallback(eventCode uint64, data []byte, tag uint64) error {
	m.Lock()
	defer m.Unlock()
	handler, ok := eventCallbackHandlerMap[eventCode]
	if !ok {
		return fmt.Errorf("event callback handler not found for event code %d",
			eventCode)
	}

	go handler.HandleEventCallback(eventCode, data, tag)

	return nil
}

//export eventCallbackGo
func eventCallbackGo(eventCode uint64, data []byte, tag uint64) {
	err := runEventCallback(eventCode, data, tag)
	if err != nil {
		log.Printf("error running event callback: %s\n", err)
	}
}
