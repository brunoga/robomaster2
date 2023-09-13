//go:build ios && arm64

package unitybridge

/*
#include <stdbool.h>
#include <stdlib.h>

#include "event_callback.h"

extern void CreateUnityBridge(const char* name, bool debuggable, const char* logPath);
extern void DestroyUnityBridge();
extern bool UnityBridgeInitialize();
extern void UnityBridgeUninitialze();
extern void UnitySendEvent(uint64_t event_code, intptr_t data, uint64_t tag);
extern void UnitySendEventWithString(uint64_t event_code, const char* data, uint64_t tag);
extern void UnitySendEventWithNumber(uint64_t event_code, uint64_t data, uint64_t tag);
extern void UnitySetEventCallback(uint64_t event_code, EventCallback event_callback);
extern intptr_t UnityGetSecurityKeyByKeyChainIndex(uint64_t index);
*/
import "C"
import (
	"log"
	"unsafe"
)

var (
	unityBridge unityBridgeImpl
)

type unityBridgeImpl struct{}

func (ub unityBridgeImpl) Create(name string, debuggable bool,
	logPath string) {
	log.Println("Creating Unity Bridge")
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cLogPath := C.CString(logPath)
	defer C.free(unsafe.Pointer(cLogPath))

	C.CreateUnityBridge(cName, C.bool(debuggable), cLogPath)
}

func (ub unityBridgeImpl) Destroy() {
	C.DestroyUnityBridge()
}

func (ub unityBridgeImpl) Initialize() bool {
	return bool(C.UnityBridgeInitialize())
}

func (ub unityBridgeImpl) Uninitialize() {
	C.UnityBridgeUninitialze()
}

func (ub unityBridgeImpl) SendEvent(eventCode uint64, data []byte, tag uint64) {
	var dataUintptr uintptr
	if len(data) == 0 {
		dataUintptr = uintptr(unsafe.Pointer(nil))
	} else {
		dataUintptr = uintptr(unsafe.Pointer(&data[0]))
	}

	C.UnitySendEvent(C.uint64_t(eventCode), C.intptr_t(dataUintptr), C.uint64_t(tag))
}

func (ub unityBridgeImpl) SendEventWithString(eventCode uint64, data string,
	tag uint64) {
	C.UnitySendEventWithString(C.uint64_t(eventCode), C.CString(data), C.uint64_t(tag))
}

func (ub unityBridgeImpl) SendEventWithNumber(eventCode uint64, data uint64,
	tag uint64) {
	C.UnitySendEventWithNumber(C.uint64_t(eventCode), C.uint64_t(data), C.uint64_t(tag))
}

func (ub unityBridgeImpl) SetEventCallback(eventCode uint64,
	handler EventCallbackHandler) {
	setEventCallbackHandler(eventCode, handler)

	var eventCallbackC C.EventCallback
	if handler != nil {
		eventCallbackC = C.EventCallback(C.eventCallbackC)
	}

	C.UnitySetEventCallback(C.uint64_t(eventCode), eventCallbackC)
}

func (ub unityBridgeImpl) GetSecurityKeyByKeyChainIndex(index uint64) string {
	cKeyUintptr := C.UnityGetSecurityKeyByKeyChainIndex(C.uint64_t(index))

	defer C.free(unsafe.Pointer(uintptr(cKeyUintptr)))

	return C.GoString((*C.char)(unsafe.Pointer(uintptr(cKeyUintptr))))
}
