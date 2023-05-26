//go:build (darwin && amd64) || (android && arm) || (android && arm64)

package gounitybridge

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
		"ios/arm64":     "libunitybridge.dylib",
		"windows/amd64": "./lib/windows/amd64/unitybridge.dll",
	}
)

func Init() error {
	log.Println("Loading Unity Bridge library")
	libPath, ok := libPaths[fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)]
	if !ok {
		// Should never happen.
		return fmt.Errorf("platform \"%s/%s\" not supported by Unity Bridge",
			runtime.GOOS, runtime.GOARCH)
	}

	log.Printf("Using path \"%s\"\n", libPath)

	cLibPath := C.CString(libPath)
	defer C.free(unsafe.Pointer(cLibPath))

	unityBridge.unityBridgeHandle = C.dlopen(cLibPath, C.RTLD_NOW)
	if unityBridge.unityBridgeHandle == nil {
		cError := C.dlerror()

		return fmt.Errorf("could not load Unity Bridge library at \"%s\": %s",
			libPath, C.GoString(cError))
	}

	log.Println("Unity Bridge library loaded. Handle obtained.")

	var err error

	unityBridge.createUnityBridge, err =
		unityBridge.getSymbol("CreateUnityBridge")
	if err != nil {
		return err
	}
	unityBridge.destroyUnityBridge, err =
		unityBridge.getSymbol("DestroyUnityBridge")
	if err != nil {
		return err
	}
	unityBridge.unityBridgeInitialize, err =
		unityBridge.getSymbol("UnityBridgeInitialize")
	if err != nil {
		return err
	}
	unityBridge.unityBridgeUninitialize, err =
		unityBridge.getSymbol("UnityBridgeUninitialze") // Typo in C code.
	if err != nil {
		return err
	}
	unityBridge.unitySendEvent, err =
		unityBridge.getSymbol("UnitySendEvent")
	if err != nil {
		return err
	}
	unityBridge.unitySendEventWithString, err =
		unityBridge.getSymbol("UnitySendEventWithString")
	if err != nil {
		return err
	}
	unityBridge.unitySendEventWithNumber, err =
		unityBridge.getSymbol("UnitySendEventWithNumber")
	if err != nil {
		return err
	}
	unityBridge.unitySetEventCallback, err =
		unityBridge.getSymbol("UnitySetEventCallback")
	if err != nil {
		return err
	}
	unityBridge.UnityGetSecurityKeyByKeyChainIndex, err =
		unityBridge.getSymbol("UnityGetSecurityKeyByKeyChainIndex")
	if err != nil {
		return err
	}

	log.Println("Unity Bridge library symbols loaded.")

	return nil
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

func (h unityBridgeImpl) getSymbol(name string) (unsafe.Pointer, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	symbol := C.dlsym(h.unityBridgeHandle, cName)
	if symbol == nil {
		cError := C.dlerror()

		return nil, fmt.Errorf("could not load symbol \"%s\": %s",
			name, C.GoString(cError))
	}

	return symbol, nil
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
	setEventCallbackHandler(eventCode, handler)

	C.UnitySetEventCallbackCaller(unsafe.Pointer(u.unitySetEventCallback),
		C.uint64_t(eventCode), C.EventCallback(C.eventCallbackC))
}

func (u unityBridgeImpl) GetSecurityKeyByKeyChainIndex(index int64) string {
	cKey := C.UnityGetSecurityKeyByKeyChainIndexCaller(
		unsafe.Pointer(u.UnityGetSecurityKeyByKeyChainIndex),
		C.int(index))

	return C.GoString(cKey)
}
