package unitybridge

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

var (
	mInstance             DJIUnityBridge
	mGetAvailableValuePtr []byte = make([]byte, 2048)
)

func init() {
	mInstance = &djiUnityBridge{}
}

func DJIUnityBridgeInstance() DJIUnityBridge {
	return mInstance
}

type DJIUnityBridge interface {
	Init()
	UnInit()
	RegisterEventHandler(handler IEventHandler, typ DJIUnityEventType)
	UnregisterEventHandler(handler IEventHandler)
	SendEvent(e *DJIUnityEvent, data []byte, tag uint64)
	SendEventWithoutTag(e *DJIUnityEvent, data []byte)
	SendEventWithoutDataOrTag(e *DJIUnityEvent)
	SendEventWithNumber(e *DJIUnityEvent, data uint64, tag uint64)
	SendEventWithString(e *DJIUnityEvent, data string, tag uint64)
	GetStringValueWithEvent(e *DJIUnityEvent) string
	GetInt32ValueWithEvent(e *DJIUnityEvent) int32
	GetSecurityKeyByKeyChainIndex(index int) string
}

type djiUnityBridge struct{}

func (d *djiUnityBridge) Init() {
	unityBridge.Create("Robomaster", true, "./log")
	d.registerCallbacks()
	if !unityBridge.Initialize() {
		panic("Failed to initialize Unity Bridge")
	}
}
func (d *djiUnityBridge) UnInit() {
	d.unregisterCallbacks()
	unityBridge.Uninitialize()
	unityBridge.Destroy()
}
func (d *djiUnityBridge) RegisterEventHandler(handler IEventHandler, typ DJIUnityEventType) {
	registerIEventHandler(typ, handler)
}
func (d *djiUnityBridge) UnregisterEventHandler(handler IEventHandler) {
	unregisterEventHandler(handler)
}
func (d *djiUnityBridge) SendEvent(e *DJIUnityEvent, data []byte, tag uint64) {
	unityBridge.SendEvent(e.GetCode(), data, tag)
}
func (d *djiUnityBridge) SendEventWithoutTag(e *DJIUnityEvent, data []byte) {
	unityBridge.SendEvent(e.GetCode(), data, 0)
}
func (d *djiUnityBridge) SendEventWithoutDataOrTag(e *DJIUnityEvent) {
	fmt.Println(e.Type(), e.SubType(), e.GetCode())
	unityBridge.SendEvent(e.GetCode(), nil, 0)
}
func (d *djiUnityBridge) SendEventWithNumber(e *DJIUnityEvent, data uint64, tag uint64) {
	unityBridge.SendEventWithNumber(e.GetCode(), data, tag)
}
func (d *djiUnityBridge) SendEventWithString(e *DJIUnityEvent, info string, tag uint64) {
	unityBridge.SendEventWithString(e.GetCode(), info, tag)
}
func (d *djiUnityBridge) GetStringValueWithEvent(e *DJIUnityEvent) string {
	d.SendEventWithoutTag(e, mGetAvailableValuePtr)
	n := bytes.IndexByte(mGetAvailableValuePtr, 0)
	if n == -1 {
		return string(mGetAvailableValuePtr)
	}
	return string(mGetAvailableValuePtr[:n])
}
func (d *djiUnityBridge) GetInt32ValueWithEvent(e *DJIUnityEvent) int32 {
	d.SendEventWithoutTag(e, mGetAvailableValuePtr)
	return int32(binary.NativeEndian.Uint32(mGetAvailableValuePtr))
}
func (d *djiUnityBridge) GetSecurityKeyByKeyChainIndex(index int) string {
	return unityBridge.GetSecurityKeyByKeyChainIndex(index)
}

func (d *djiUnityBridge) registerCallbacks() {
	event := NewDJIUnityEventZero()
	for _, eventType := range DJIUnityEventTypes() {
		event.Reset(eventType, 0)
		unityBridge.SetEventCallback(event.GetCode(), true)
	}

}
func (d *djiUnityBridge) unregisterCallbacks() {
	event := NewDJIUnityEventZero()
	for _, eventType := range DJIUnityEventTypes() {
		event.Reset(eventType, 0)
		unityBridge.SetEventCallback(event.GetCode(), false)
	}
}
