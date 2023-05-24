//go:build !((darwin && amd64) || (android && arm) || (android && arm64) || (ios && arm64) || (windows && amd64))

package unitybridge

import (
	"fmt"
	"runtime"
)

func init() {
	panic(fmt.Sprintf("Platform \"%s/%s\" not supported by Unity Bridge",
		runtime.GOOS, runtime.GOARCH))
}

type unityBridgeImpl struct{}

var unityBridge unityBridgeImpl

// Empty placeholders for unsupported platforms.

func (ub unityBridgeImpl) Create(name string, debuggable bool, logPath string) {}

func (ub unityBridgeImpl) Destroy() {}

func (ub unityBridgeImpl) Initialize() bool {
	return false
}

func (ub unityBridgeImpl) Uninitialize() {}

func (ub unityBridgeImpl) SendEvent(eventCode uint64, data []byte, tag uint64) {}

func (ub unityBridgeImpl) SendEventWithString(eventCode uint64, data string, tag uint64) {}

func (ub unityBridgeImpl) SendEventWithNumber(eventCode uint64, data uint64, tag uint64) {}

func (ub unityBridgeImpl) SetEventCallback(eventCode uint64, callback EventCallback) {}

func (ub unityBridgeImpl) GetSecurityKeyByKeyChainIndex(index uint64) string {
	return ""
}
