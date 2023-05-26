//go:build linux && (amd64 || arm64)

package gounitybridge

import (
	"debug/pe"
	"fmt"
	"os/exec"
)

const (
	dllHostExe = "dllhost.exe"
)

var (
	unityBridge unityBridgeImpl
)

func init() {
	// Check if wine is available.
	_, err := getWinePath()
	if err != nil {
		panic(err)
	}

	// Check if unitybridge_dll_host.exe is available and is a Windows
	// executable.
	err = checkDLLHostExe()
	if err != nil {
		panic(err)
	}
}

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

func getWinePath() (string, error) {
	return exec.LookPath("wine")
}

func checkDLLHostExe() error {
	peFile, err := pe.Open(dllHostExe)
	if err != nil {
		return fmt.Errorf("%q does not look like a Windows executable: %w",
			dllHostExe, err)
	}
	peFile.Close()

	return nil
}
