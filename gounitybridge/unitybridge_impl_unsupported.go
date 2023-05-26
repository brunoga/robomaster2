//go:build !((darwin && amd64) || (android && arm) || (android && arm64) || (ios && arm64) || (windows && amd64) || (linux && (amd64 || arm64)))

package gounitybridge

import (
	"fmt"
	"runtime"
)

var (
	unityBridge unityBridgeImpl
)

func init() {
	panic(fmt.Sprintf("Platform \"%s/%s\" not supported by Unity Bridge",
		runtime.GOOS, runtime.GOARCH))
}

type unityBridgeImpl struct{}

// Empty placeholders for unsupported platforms.

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
