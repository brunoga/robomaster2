//go:build ios && arm64

package unitybridge

/*
#include <stdbool.h>
#include <stdlib.h>

extern void CreateUnityBridge(const char* name, bool debuggable, const char* logPath);
extern void DestroyUnityBridge();
extern bool UnityBridgeInitialize();
extern void UnityBridgeUninitialze();
extern void UnityBridgeSendEvent(uint64_t event_code, intptr_t data, uint64_t tag);
extern void UnityBridgeSendEventWithString(uint64_t event_code, const char* data, uint64_t tag);
extern void UnityBridgeSendEventWithNumber(uint64_t event_code, uint64_t data, uint64_t tag);
extern void UnitySetEventCallback(uint64_t event_code, EventCallback event_callback);
extern uintptr_t UnityGetSecurityKeyByKeyChainIndex(uint64_t index);
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

func (ub unityBridgeImpl) SendEvent(eventCode int64, data []byte, tag int64) {
	if len(data) == 0 {
		dataUintptr = uintptr(unsafe.Pointer(nil))
	} else {
		dataUintptr = uintptr(unsafe.Pointer(&data[0]))
	}

	C.UnityBridgeSendEvent(eventCode, dataUintptr, tag)
}

func (ub unityBridgeImpl) SendEventWithString(eventCode int64, data string,
	tag int64) {
	C.UnityBridgeSendEventWithString(eventCode, C.CString(data), tag)
}

func (ub unityBridgeImpl) SendEventWithNumber(eventCode int64, data int64,
	tag int64) {
	C.UnityBridgeSendEventWithNumber(eventCode, data, tag)
}

func (ub unityBridgeImpl) SetEventCallback(eventCode int64,
	handler EventCallbackHandler) {
	setEventCallbackHandler(eventCode, handler)

	var eventCallbackC C.EventCallback
	if handler != nil {
		eventCallbackC = C.eventCallbackC
	}

	C.UnitySetEventCallback(eventCode, eventCallbackC)
}

func (ub unityBridgeImpl) GetSecurityKeyByKeyChainIndex(index int64) string {
	cKeyUintptr := C.UnityBridgeGetSecurityKeyByKeyChainIndex(C.uint64_t(index))

	defer C.free(unsafe.Pointer(cKeyUintptr))

	return C.GoString((*C.char)(unsafe.Pointer(cKeyUintptr)))
}
