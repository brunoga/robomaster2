package gounitybridge

type EventCallbackHandler interface {
	HandleEventCallback(eventCode int64, data []byte, tag int64)
}

type UnityBridge interface {
	Create(name string, debuggable bool, logPath string)
	Destroy()
	Initialize() bool
	Uninitialize()
	SendEvent(eventCode int64, data []byte, tag int64)
	SendEventWithString(eventCode int64, data string, tag int64)
	SendEventWithNumber(eventCode int64, data int64, tag int64)
	SetEventCallback(eventCode int64, handler EventCallbackHandler)
	GetSecurityKeyByKeyChainIndex(index int64) string
}

func Get() UnityBridge {
	return unityBridge
}
