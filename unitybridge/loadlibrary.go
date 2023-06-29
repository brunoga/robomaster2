//go:build windows && amd64

package unitybridge

/*
#include <stdbool.h>
#include <stdlib.h>

#include "event_callback.h"
*/
import "C"
import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	unityBridge unityBridgeImpl

	libPath = "./lib/windows/amd64/unitybridge.dll"
)

func init() {
	var err error

	unityBridge.unityBridgeHandle, err = syscall.LoadDLL(libPath)
	if err != nil {
		panic(fmt.Sprintf("Could not load Unity Bridge library at \"%s\": %s", libPath, err))
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

type unityBridgeImpl struct {
	unityBridgeHandle *syscall.DLL

	createUnityBridge                  *syscall.Proc
	destroyUnityBridge                 *syscall.Proc
	unityBridgeInitialize              *syscall.Proc
	unityBridgeUninitialize            *syscall.Proc
	unitySendEvent                     *syscall.Proc
	unitySendEventWithString           *syscall.Proc
	unitySendEventWithNumber           *syscall.Proc
	unitySetEventCallback              *syscall.Proc
	UnityGetSecurityKeyByKeyChainIndex *syscall.Proc
}

func (h unityBridgeImpl) getSymbol(name string) *syscall.Proc {
	symbol, err := h.unityBridgeHandle.FindProc(name)
	if err != nil {
		panic(fmt.Sprintf("Could not load symbol \"%s\": %s", name, err))
	}

	return symbol
}

func (u unityBridgeImpl) Create(name string, debuggable bool, logPath string) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	intDebuggable := 0
	if debuggable {
		intDebuggable = 1
	}

	cLogPath := C.CString(logPath)
	defer C.free(unsafe.Pointer(cLogPath))

	_, _, _ = u.createUnityBridge.Call(
		uintptr(unsafe.Pointer(cName)),
		uintptr(intDebuggable),
		uintptr(unsafe.Pointer(cLogPath)),
	)
}

func (u unityBridgeImpl) Destroy() {
	_, _, _ = u.destroyUnityBridge.Call()
}

func (u unityBridgeImpl) Initialize() bool {
	ret, _, _ := u.unityBridgeInitialize.Call()
	return ret != 0
}

func (u unityBridgeImpl) Uninitialize() {
	_, _, _ = u.unityBridgeUninitialize.Call()
}

func (u unityBridgeImpl) SendEvent(eventCode int64, data []byte, tag int64) {
	_, _, _ = u.unitySendEvent.Call(
		uintptr(eventCode),
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)),
		uintptr(tag),
	)
}

func (u unityBridgeImpl) SendEventWithString(eventCode int64, data string,
	tag int64) {
	cData := C.CString(data)
	defer C.free(unsafe.Pointer(cData))

	_, _, _ = u.unitySendEventWithString.Call(
		uintptr(eventCode),
		uintptr(unsafe.Pointer(cData)),
		uintptr(tag),
	)
}

func (u unityBridgeImpl) SendEventWithNumber(eventCode int64, data int64,
	tag int64) {
	_, _, _ = u.unitySendEventWithNumber.Call(
		uintptr(eventCode),
		uintptr(data),
		uintptr(tag),
	)
}

func (u unityBridgeImpl) SetEventCallback(eventCode int64,
	handler EventCallbackHandler) {
	setEventCallbackHandler(eventCode, handler)

	_, _, _ = u.unitySetEventCallback.Call(
		uintptr(eventCode),
		uintptr(C.eventCallbackC),
	)
}

func (u unityBridgeImpl) GetSecurityKeyByKeyChainIndex(index int64) string {
	result, _, _ := u.UnityGetSecurityKeyByKeyChainIndex.Call(
		uintptr(index),
	)

	return C.GoString((*C.char)(unsafe.Pointer(result)))
}
