//go:build (darwin && amd64) || (android && arm) || (android && arm64) || (ios && arm64)

package unitybridge

/*
#include <dlfcn.h>
#include <stdint.h>
#include <stdlib.h>

#include "event_callback.h"
#include "function_caller.h"
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

var (
	unityBridge unityBridgeImpl
)

type unityBridgeImpl struct {
	unityBridgeHandle unsafe.Pointer

	createUnityBridge                  unsafe.Pointer
	destroyUnityBridge                 unsafe.Pointer
	unityBridgeInitialize              unsafe.Pointer
	unityBridgeUninitialize            unsafe.Pointer
	unitySendEvent                     unsafe.Pointer
	unitySendEventWithString           unsafe.Pointer
	unitySendEventWithNumber           unsafe.Pointer
	unitySetEventCallback              unsafe.Pointer
	UnityGetSecurityKeyByKeyChainIndex unsafe.Pointer
}

var libPaths = map[string]string{
	"darwin/amd64":  "./lib/darwin/amd64/unitybridge.bundle/Contents/MacOS/unitybridge",
	"android/arm":   "./lib/android/arm/libunitybridge.so",
	"android/arm64": "./lib/android/arm64/libunitybridge.so",
	"ios/arm64":     "./lib/ios/arm64/unitybridge.framework/unitybridge",
	"windows/amd64": "./lib/windows/amd64/unitybridge.dll",
}

func init() {
	libPath, ok := libPaths[fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)]
	if !ok {
		// Should never happen.
		panic(fmt.Sprintf("Platform \"%s/%s\" not supported by Unity Bridge",
			runtime.GOOS, runtime.GOARCH))
	}

	cLibPath := C.CString(libPath)
	defer C.free(unsafe.Pointer(cLibPath))

	unityBridge.unityBridgeHandle = C.dlopen(cLibPath, C.RTLD_NOW)
	if unityBridge.unityBridgeHandle == nil {
		cError := C.dlerror()

		panic(fmt.Sprintf("Could not load Unity Bridge library at \"%s\": %s",
			libPath, C.GoString(cError)))
	}

	unityBridge.createUnityBridge =
		unityBridge.getSymbol("CreateUnityBridge")
	unityBridge.destroyUnityBridge =
		unityBridge.getSymbol("DestroyUnityBridge")
	unityBridge.unityBridgeInitialize =
		unityBridge.getSymbol("UnityBridgeInitialize")
	unityBridge.unityBridgeUninitialize =
		unityBridge.getSymbol("UnityBridgeUninitialze") // Typo in C code.
	unityBridge.unitySendEvent =
		unityBridge.getSymbol("UnitySendEvent")
	unityBridge.unitySendEventWithString =
		unityBridge.getSymbol("UnitySendEventWithString")
	unityBridge.unitySendEventWithNumber =
		unityBridge.getSymbol("UnitySendEventWithNumber")
	unityBridge.unitySetEventCallback =
		unityBridge.getSymbol("UnitySetEventCallback")
	unityBridge.UnityGetSecurityKeyByKeyChainIndex =
		unityBridge.getSymbol("UnityGetSecurityKeyByKeyChainIndex")
}

func (h unityBridgeImpl) getSymbol(name string) unsafe.Pointer {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	symbol := C.dlsym(h.unityBridgeHandle, cName)
	if symbol == nil {
		cError := C.dlerror()

		panic(fmt.Sprintf("Could not load symbol \"%s\": %s",
			name, C.GoString(cError)))
	}

	return symbol
}

func (u unityBridgeImpl) Create(name string, debuggable bool, logPath string) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cLogPath := C.CString(logPath)
	defer C.free(unsafe.Pointer(cLogPath))

	C.CreateUnityBridgeCaller(unsafe.Pointer(u.createUnityBridge), cName,
		C.bool(debuggable), cLogPath)
}

func (u unityBridgeImpl) Destroy() {
	C.DestroyUnityBridgeCaller(unsafe.Pointer(u.destroyUnityBridge))
}

func (u unityBridgeImpl) Initialize() bool {
	return bool(C.UnityBridgeInitializeCaller(unsafe.Pointer(u.unityBridgeInitialize)))
}

func (u unityBridgeImpl) Uninitialize() {
	C.UnityBridgeUninitializeCaller(unsafe.Pointer(u.unityBridgeUninitialize))
}

func (u unityBridgeImpl) SendEvent(eventCode int64, data []byte,
	tag int64) {
	C.UnitySendEventCaller(unsafe.Pointer(u.unitySendEvent),
		C.uint64_t(eventCode), C.uintptr_t(uintptr(unsafe.Pointer(&data[0]))),
		C.int(len(data)), C.uint64_t(tag))
}

func (u unityBridgeImpl) SendEventWithString(eventCode int64, data string,
	tag int64) {
	cData := C.CString(data)
	defer C.free(unsafe.Pointer(cData))

	C.UnitySendEventWithStringCaller(unsafe.Pointer(u.unitySendEventWithString),
		C.uint64_t(eventCode), cData, C.uint64_t(tag))
}

func (u unityBridgeImpl) SendEventWithNumber(eventCode int64, data int64,
	tag int64) {
	C.UnitySendEventWithNumberCaller(unsafe.Pointer(u.unitySendEventWithNumber),
		C.uint64_t(eventCode), C.uint64_t(data), C.uint64_t(tag))
}

func (u unityBridgeImpl) SetEventCallback(eventCode int64,
	handler EventCallbackHandler) {
	// TODO(bga): Keep track of event callback on the Go side.

	C.UnitySetEventCallbackCaller(unsafe.Pointer(u.unitySetEventCallback),
		C.uint64_t(eventCode), C.EventCallback(C.cEventCallback))
}

func (u unityBridgeImpl) GetSecurityKeyByKeyChainIndex(index int64) string {
	cKey := C.UnityGetSecurityKeyByKeyChainIndexCaller(
		unsafe.Pointer(u.UnityGetSecurityKeyByKeyChainIndex),
		C.int(index))

	return C.GoString(cKey)
}
