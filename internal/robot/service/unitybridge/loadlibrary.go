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

func (u unityBridgeImpl) SendEvent(eventCode uint64, data []byte, tag uint64) {
	var dataUintptr uintptr
	if len(data) == 0 {
		dataUintptr = uintptr(unsafe.Pointer(nil))
	} else {
		dataUintptr = uintptr(unsafe.Pointer(&data[0]))
	}
	_, _, _ = u.unitySendEvent.Call(
		uintptr(eventCode),
		dataUintptr,
		uintptr(tag),
	)
}

func (u unityBridgeImpl) SendEventWithString(eventCode uint64, data string,
	tag uint64) {

	cData := C.CString(data)
	defer C.free(unsafe.Pointer(cData))

	_, _, _ = u.unitySendEventWithString.Call(
		uintptr(eventCode),
		uintptr(unsafe.Pointer(cData)),
		uintptr(tag),
	)
}

func (u unityBridgeImpl) SendEventWithNumber(eventCode uint64, data uint64,
	tag uint64) {
	_, _, _ = u.unitySendEventWithNumber.Call(
		uintptr(eventCode),
		uintptr(data),
		uintptr(tag),
	)
}

func (u unityBridgeImpl) SetEventCallback(eventCode uint64, add bool) {
	var eventCallbackUintptr uintptr
	if add {
		eventCallbackUintptr = uintptr(C.eventCallbackC)
	}

	_, _, _ = u.unitySetEventCallback.Call(
		uintptr(eventCode),
		eventCallbackUintptr,
	)
}

func (u unityBridgeImpl) GetSecurityKeyByKeyChainIndex(index int) string {
	cKeyUintptr, _, _ := u.UnityGetSecurityKeyByKeyChainIndex.Call(
		uintptr(index),
	)

	defer C.free(unsafe.Pointer(cKeyUintptr))

	return C.GoString((*C.char)(unsafe.Pointer(cKeyUintptr)))
}
