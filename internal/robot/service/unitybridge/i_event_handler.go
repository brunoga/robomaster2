package unitybridge

type IEventHandler interface {
	OnEventCallback(e *DJIUnityEvent, data []byte, tag uint64)
}
