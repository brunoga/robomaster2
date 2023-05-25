//go:build linux && (amd64 || arm64)

package unitybridge

var (
	unityBridge unityBridgeImpl
)

type unityBridgeImpl struct{}

func (ub unityBridgeImpl) Create(name string, debuggable bool,
	logPath string) {
}

func (ub unityBridgeImpl) Destroy() {}

func (ub unityBridgeImpl) Initialize() bool {
	return false
}

func (ub unityBridgeImpl) Uninitialize() {}

func (ub unityBridgeImpl) SendEvent(eventCode int64, data []byte, tag int64) {}

func (ub unityBridgeImpl) SendEventWithString(eventCode int64, data string,
	tag int64) {
}

func (ub unityBridgeImpl) SendEventWithNumber(eventCode int64, data int64,
	tag int64) {
}

func (ub unityBridgeImpl) SetEventCallback(eventCode int64,
	handler EventCallbackHandler) {
}

func (ub unityBridgeImpl) GetSecurityKeyByKeyChainIndex(index int64) string {
	return ""
}
