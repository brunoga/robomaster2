//go:build (darwin && amd64) || (android && arm) || (android && arm64)

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
	"log"
	"runtime"
	"unsafe"
)

var (
	unityBridge unityBridgeImpl

	libPaths = map[string]string{
		"android/arm":   "./lib/android/arm/libunitybridge.so",
		"android/arm64": "./lib/android/arm64/libunitybridge.so",
		"darwin/amd64":  "./lib/darwin/amd64/unitybridge.bundle/Contents/MacOS/unitybridge",
	}
)

func init() {
	log.Println("Loading Unity Bridge library")
	libPath, ok := libPaths[fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)]
	if !ok {
		// Should never happen.
		panic(fmt.Sprintf("platform \"%s/%s\" not supported by Unity Bridge",
			runtime.GOOS, runtime.GOARCH))
	}

	log.Printf("Using path \"%s\"\n", libPath)

	cLibPath := C.CString(libPath)
	defer C.free(unsafe.Pointer(cLibPath))

	unityBridge.unityBridgeHandle = C.dlopen(cLibPath, C.RTLD_NOW)
	if unityBridge.unityBridgeHandle == nil {
		cError := C.dlerror()

		panic(fmt.Sprintf("could not load Unity Bridge library at \"%s\": %s",
			libPath, C.GoString(cError)))
	}

	log.Println("Unity Bridge library loaded. Handle obtained.")

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

	log.Println("Unity Bridge library symbols loaded.")
}

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

func (h unityBridgeImpl) getSymbol(name string) unsafe.Pointer {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	symbol := C.dlsym(h.unityBridgeHandle, cName)
	if symbol == nil {
		cError := C.dlerror()

		panic(fmt.Sprintf("could not load symbol \"%s\": %s",
			name, C.GoString(cError)))
	}

	return symbol
}

func (u unityBridgeImpl) Create(name string, debuggable bool, logPath string) {
	log.Println("Creating Unity Bridge")
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cLogPath := C.CString(logPath)
	defer C.free(unsafe.Pointer(cLogPath))

	C.CreateUnityBridgeCaller(unsafe.Pointer(u.createUnityBridge), cName,
		C.bool(debuggable), cLogPath)
}

func (u unityBridgeImpl) Destroy() {
	log.Println("Destroying Unity Bridge")
	C.DestroyUnityBridgeCaller(unsafe.Pointer(u.destroyUnityBridge))
}

func (u unityBridgeImpl) Initialize() bool {
	log.Println("Initializing Unity Bridge")
	return bool(C.UnityBridgeInitializeCaller(unsafe.Pointer(u.unityBridgeInitialize)))
}

func (u unityBridgeImpl) Uninitialize() {
	log.Println("Uninitializing Unity Bridge")
	C.UnityBridgeUninitializeCaller(unsafe.Pointer(u.unityBridgeUninitialize))
}

func (u unityBridgeImpl) SendEvent(eventCode uint64, data []byte,
	tag uint64) {
	var dataUintptr uintptr
	if len(data) > 0 {
		dataUintptr = uintptr(unsafe.Pointer(&data[0]))
	}

	C.UnitySendEventCaller(unsafe.Pointer(u.unitySendEvent),
		C.uint64_t(eventCode), C.uintptr_t(dataUintptr),
		C.int(len(data)), C.uint64_t(tag))
}

func (u unityBridgeImpl) SendEventWithString(eventCode uint64, data string,
	tag uint64) {
	cData := C.CString(data)
	defer C.free(unsafe.Pointer(cData))

	C.UnitySendEventWithStringCaller(unsafe.Pointer(u.unitySendEventWithString),
		C.uint64_t(eventCode), cData, C.uint64_t(tag))
}

func (u unityBridgeImpl) SendEventWithNumber(eventCode uint64, data uint64,
	tag uint64) {
	C.UnitySendEventWithNumberCaller(unsafe.Pointer(u.unitySendEventWithNumber),
		C.uint64_t(eventCode), C.uint64_t(data), C.uint64_t(tag))
}

func (u unityBridgeImpl) SetEventCallback(eventCode uint64,
	add bool) {
	var eventCallback C.EventCallback
	if add {
		eventCallback = C.EventCallback(C.eventCallbackC)
	}

	C.UnitySetEventCallbackCaller(unsafe.Pointer(u.unitySetEventCallback),
		C.uint64_t(eventCode), eventCallback)
}

func (u unityBridgeImpl) GetSecurityKeyByKeyChainIndex(index uint64) string {
	cKey := C.UnityGetSecurityKeyByKeyChainIndexCaller(
		unsafe.Pointer(u.UnityGetSecurityKeyByKeyChainIndex),
		C.int(index))

	return C.GoString(cKey)
}
