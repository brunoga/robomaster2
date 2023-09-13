package unitybridge

type EventCallbackHandler interface {
	HandleEventCallback(eventCode uint64, data []byte, tag uint64)
}

type UnityBridge interface {
	Create(name string, debuggable bool, logPath string)
	Destroy()
	Initialize() bool
	Uninitialize()
	SendEvent(eventCode uint64, data []byte, tag uint64)
	SendEventWithString(eventCode uint64, data string, tag uint64)
	SendEventWithNumber(eventCode uint64, data uint64, tag uint64)
	SetEventCallback(eventCode uint64, handler EventCallbackHandler)
	GetSecurityKeyByKeyChainIndex(index uint64) string
}

func Get() UnityBridge {
	return unityBridge
}
