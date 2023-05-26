//go:build none

package unitybridge

/*
#cgo LDFLAGS: -L. -lunitybridge.framework/unitybridge

#include <stdbool.h>
#include <stdlib.h>

void CreateUnityBridge(const char* name, bool debuggable, const char* logPath);
void DestroyUnityBridge();
bool UnityBridgeInitialize();
void UnityBridgeUninitialze();
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

func (ub unityBridgeImpl) SendEvent(eventCode int64, data []byte, tag int64) {}

func (ub unityBridgeImpl) SendEventWithString(eventCode int64, data string,
	tag int64) {
}

func (ub unityBridgeImpl) SendEventWithNumber(eventCode int64, data int64,
	tag int64) {
}

func (ub unityBridgeImpl) SetEventCallback(eventCode int64,
	handler EventCallbackHandler) {
}

func (ub unityBridgeImpl) GetSecurityKeyByKeyChainIndex(index int64) string {
	return ""
}