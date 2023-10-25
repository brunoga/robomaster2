package service

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/brunoga/robomaster2/internal/robot/service/dji"
	"github.com/brunoga/robomaster2/internal/robot/service/unitybridge"
)

var (
	mInstance DJICommandController

	mGetEvent            = unitybridge.NewDJIUnityEvent(uint64(unitybridge.GetValue))
	mGetAvailableEvent   = unitybridge.NewDJIUnityEvent(uint64(unitybridge.GetAvailableValue))
	mSetEvent            = unitybridge.NewDJIUnityEvent(uint64(unitybridge.SetValue))
	mActionEvent         = unitybridge.NewDJIUnityEvent(uint64(unitybridge.PerformAction))
	mStartListeningEvent = unitybridge.NewDJIUnityEvent(uint64(unitybridge.StartListening))
	mStopListeningEvent  = unitybridge.NewDJIUnityEvent(uint64(unitybridge.StopListening))
)

func init() {
	mInstance = &djiCommandController{
		listenersByKeyAndName: make(map[dji.DJIKeys]map[string]DJIListener),
		callbackDelegate:      NewCallbackDictionary(),
	}
}

type DJICommandController interface {
	Init()
	UnInit()

	unitybridge.IEventHandler

	StartListeningOnKey(key dji.DJIKeys, listener interface{},
		callback func(*dji.DJIResult), fetchFromCache bool)
}

type djiCommandController struct {
	listenersByKeyAndName map[dji.DJIKeys]map[string]DJIListener
	callbackDelegate      *CallbackDictionary
}

func DJICommandControllerInstance() DJICommandController {
	return mInstance
}

func (d *djiCommandController) Init() {
	unitybridge.DJIUnityBridgeInstance().RegisterEventHandler(
		d, unitybridge.GetValue)
	unitybridge.DJIUnityBridgeInstance().RegisterEventHandler(
		d, unitybridge.SetValue)
	unitybridge.DJIUnityBridgeInstance().RegisterEventHandler(
		d, unitybridge.PerformAction)
	unitybridge.DJIUnityBridgeInstance().RegisterEventHandler(
		d, unitybridge.StartListening)
}

func (d *djiCommandController) UnInit() {
	unitybridge.DJIUnityBridgeInstance().UnregisterEventHandler(d)
}

func (d *djiCommandController) OnEventCallback(event *unitybridge.DJIUnityEvent,
	data []byte, tag uint64) {
	switch unitybridge.DJIUnityDataType(byte(tag>>56) & 0xff) {
	case unitybridge.String:
		n := bytes.IndexByte(data, 0)
		if n == -1 {
			d.onCommandEventCallback(event, string(data), tag)
		} else {
			d.onCommandEventCallback(event, string(data[:n]), tag)
		}
	case unitybridge.Number:
		d.onCommandEventCallback(event, fmt.Sprintf(
			"%d", binary.NativeEndian.Uint32(data)), tag)
	}
}

func (d *djiCommandController) StartListeningOnKey(key dji.DJIKeys,
	listener interface{}, callback func(*dji.DJIResult), fetchFromCache bool) {
	if callback == nil {
		return
	}

	if key.AccessType()&dji.AccessType_Read == 0 {
		panic(fmt.Sprintf("Key %d is not readable.", key))
	}

	flag := false

	listeners := d.listenersByKeyAndName[key]
	if listeners == nil {
		d.listenersByKeyAndName[key] = make(map[string]DJIListener)
		flag = true
	}

	djiListener := DJIListener{
		Name:     getFullName(listener),
		Callback: callback,
	}

	d.listenersByKeyAndName[key][djiListener.Name] = djiListener
	if flag {
		mStartListeningEvent.ResetSubType(key.Value())
		unitybridge.DJIUnityBridgeInstance().SendEventWithoutDataOrTag(mStartListeningEvent)
	}

	if !fetchFromCache {
		return
	}

	availableValueForKey := d.GetAvailableValueForKey(key)
	if !availableValueForKey.Succeeded() {
		return
	}

	callback(availableValueForKey)
}

func (d *djiCommandController) StopListening(listener interface{}) {
	name := getFullName(listener)
	for key, listenersByName := range d.listenersByKeyAndName {
		if _, exists := listenersByName[name]; exists {
			delete(listenersByName, name)
			if len(listenersByName) == 0 {
				mStopListeningEvent.ResetSubType(key.Value())
				unitybridge.DJIUnityBridgeInstance().SendEventWithoutDataOrTag(mStopListeningEvent)
			}
		}
	}
}

func (d *djiCommandController) StopListeningOnKey(key dji.DJIKeys,
	listener interface{}) {
	name := getFullName(listener)
	listenersByName, ok := d.listenersByKeyAndName[key]
	if !ok {
		return
	}

	delete(listenersByName, name)
	if len(listenersByName) == 0 {
		mStopListeningEvent.ResetSubType(key.Value())
		unitybridge.DJIUnityBridgeInstance().SendEventWithoutDataOrTag(mStopListeningEvent)
	}

	if len(d.listenersByKeyAndName) == 0 {
		delete(d.listenersByKeyAndName, key)
	}
}

func (d *djiCommandController) GetAvailableValueForKey(key dji.DJIKeys) *dji.DJIResult {
	if key.AccessType()&dji.AccessType_Read == 0 {
		panic(fmt.Sprintf("Key %d is not readable.", key))
	}

	mGetAvailableEvent.ResetSubType(key.Value())

	result := dji.NewDJIResultFromJSON(
		[]byte(unitybridge.DJIUnityBridgeInstance().GetStringValueWithEvent(
			mGetAvailableEvent)))

	return result
}

func (d *djiCommandController) GetValueForKey(key dji.DJIKeys, callback func(*dji.DJIResult)) {
	if key.AccessType()&dji.AccessType_Read == 0 {
		panic(fmt.Sprintf("Key %d is not readable.", key))
	}

	subType := key.Value()
	tag := d.callbackDelegate.AddAction(Getter, callback)
	mGetEvent.ResetSubType(subType)
	unitybridge.DJIUnityBridgeInstance().SendEvent(mGetEvent, nil, uint64(tag))
}

func (d *djiCommandController) SetValueForKey(key dji.DJIKeys,
	paramValue dji.DJIParamValue, callback func(*dji.DJIResult)) {
	if key.AccessType()&dji.AccessType_Write == 0 {
		panic(fmt.Sprintf("Key %d is not writable.", key))
	}

	subType := key.Value()
	tag := d.callbackDelegate.AddAction(Setter, callback)

	mSetEvent.ResetSubType(subType)

	data, err := json.Marshal(
		struct {
			Value any
		}{
			Value: paramValue,
		},
	)
	if err != nil {
		panic(err)
	}

	unitybridge.DJIUnityBridgeInstance().SendEvent(mSetEvent, data, uint64(tag))
}

func (d *djiCommandController) SetValueForKerWithNumber(key dji.DJIKeys,
	value int64, callback func(*dji.DJIResult)) {
	d.SetValueForKey(key, dji.DJILongParamValue(value), callback)
}

func (d *djiCommandController) PerformAction(key dji.DJIKeys, callback func(*dji.DJIResult)) {
	d.PerformActionWithParam(key, nil, callback)
}

func (d *djiCommandController) PerformActionWithParam(key dji.DJIKeys,
	value dji.DJIParamValue, callback func(*dji.DJIResult)) {
	if key.AccessType()&dji.AccessType_Action == 0 {
		panic(fmt.Sprintf("Key %d is not an action.", key))
	}

	subType := key.Value()
	tag := d.callbackDelegate.AddAction(Setter, callback)
	mActionEvent.ResetSubType(subType)

	var data []byte
	if value != nil {
		var err error
		data, err = json.Marshal(struct{ Value any }{Value: value})
		if err != nil {
			panic(err)
		}
	}

	unitybridge.DJIUnityBridgeInstance().SendEvent(mActionEvent, data, uint64(tag))
}

func (d *djiCommandController) PerformActionWithNumber(key dji.DJIKeys,
	value int64, callback func(*dji.DJIResult)) {
	d.PerformActionWithParam(key, dji.DJILongParamValue(value), callback)
}

func (d *djiCommandController) DirectSendValue(key dji.DJIKeys, value int64) {
	mActionEvent.ResetSubType(key.Value())

	data, err := json.Marshal(struct{ Value int64 }{Value: value})
	if err != nil {
		panic(err)
	}

	unitybridge.DJIUnityBridgeInstance().SendEvent(mActionEvent, data, 0)
}

func (d *djiCommandController) onCommandEventCallback(
	event *unitybridge.DJIUnityEvent, value string, tag uint64) {
	switch event.Type() {
	case unitybridge.SetValue:
		fallthrough
	case unitybridge.PerformAction:
		d.onSetValueCallback(value, tag)
	case unitybridge.GetValue:
		d.onGetValueCallback(value, tag)
	case unitybridge.StartListening:
		d.onListeningUpdates(value, tag)
	}
}

func (d *djiCommandController) onGetValueCallback(value string, tag uint64) {
	d.callbackDelegate.Invoke(Getter, dji.NewDJIResultFromJSON([]byte(value)))
	d.callbackDelegate.RemoveAction(Getter, uint32(tag))
}

func (d *djiCommandController) onSetValueCallback(value string, tag uint64) {
	d.callbackDelegate.Invoke(Setter, dji.NewDJIResultFromJSON([]byte(value)))
	d.callbackDelegate.RemoveAction(Setter, uint32(tag))
}

func (d *djiCommandController) onListeningUpdates(value string, tag uint64) {
	if len(value) == 0 {
		return
	}

	d.callbackListening(dji.NewDJIResultFromJSON([]byte(value)))
}

func (d *djiCommandController) callbackListening(param *dji.DJIResult) {
	listeners, ok := d.listenersByKeyAndName[param.Key()]
	if !ok {
		panic(fmt.Sprintf("No listeners for key %d.", param.Key()))
	}

	for _, listener := range listeners {
		listener.Callback(param)
	}
}

func getFullName(t interface{}) string {
	return reflect.TypeOf(t).PkgPath() + "." + reflect.TypeOf(t).Name()
}
