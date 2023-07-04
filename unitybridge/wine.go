//go:build linux && (amd64 || arm64)

package unitybridge

import (
	"bytes"
	"debug/pe"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

const (
	dllHostExe = "dllhost.exe"
)

var (
	unityBridge unityBridgeImpl

	localReadPipe  *os.File
	localWritePipe *os.File
)

func init() {
	// Check if wine is available.
	winePath, err := getWinePath()
	if err != nil {
		panic(err)
	}

	// Check if dllhost.exe is available and is a Windows executable.
	dllHostPath, err := getDLLHostPath()
	if err != nil {
		panic(err)
	}

	var remoteWritePipe *os.File
	localReadPipe, remoteWritePipe, err = os.Pipe()
	if err != nil {
		panic(err)
	}

	var remoteReadPipe *os.File
	remoteReadPipe, localWritePipe, err = os.Pipe()
	if err != nil {
		panic(err)
	}

	err = startDllHost(winePath, dllHostPath, remoteReadPipe, remoteWritePipe)
	if err != nil {
		panic(err)
	}

	go loop()
}

type unityBridgeImpl struct{}

func (ub unityBridgeImpl) Create(name string, debuggable bool,
	logPath string) {
	var b bytes.Buffer

	b.WriteByte(0x00)
	binary.Write(&b, binary.LittleEndian, int64(1+2+len(name)+2+len(logPath)))
	if debuggable {
		b.WriteByte(1)
	} else {
		b.WriteByte(0)
	}

	binary.Write(&b, binary.LittleEndian, uint16(len(name)))
	b.WriteString(name)

	binary.Write(&b, binary.LittleEndian, uint16(len(logPath)))
	b.WriteString(logPath)

	_, err := b.WriteTo(localWritePipe)
	if err != nil {
		panic(err)
	}

	response := make([]byte, 1)
	_, err = localReadPipe.Read(response)
	if err != nil {
		panic(err)
	}

	if response[0] != 0x00 {
		panic(fmt.Sprintf("Unexpected function byte: %d", response))
	}
}

func (ub unityBridgeImpl) Destroy() {
	var b bytes.Buffer

	b.WriteByte(0x01)

	_, err := b.WriteTo(localWritePipe)
	if err != nil {
		panic(err)
	}

	response := make([]byte, 1)
	_, err = localReadPipe.Read(response)
	if err != nil {
		panic(err)
	}

	if response[0] != 0x01 {
		panic(fmt.Sprintf("Unexpected function byte: %d", response))
	}
}

func (ub unityBridgeImpl) Initialize() bool {
	var b bytes.Buffer

	b.WriteByte(0x02)

	_, err := b.WriteTo(localWritePipe)
	if err != nil {
		panic(err)
	}

	response := make([]byte, 2)
	_, err = localReadPipe.Read(response)
	if err != nil {
		panic(err)
	}

	if response[0] != 0x02 {
		panic(fmt.Sprintf("Unexpected function byte: %d", response))
	}

	var result bool
	binary.Read(localReadPipe, binary.LittleEndian, &result)

	return result
}

func (ub unityBridgeImpl) Uninitialize() {
	var b bytes.Buffer

	b.WriteByte(0x03)

	_, err := b.WriteTo(localWritePipe)
	if err != nil {
		panic(err)
	}

	response := make([]byte, 1)
	_, err = localReadPipe.Read(response)
	if err != nil {
		panic(err)
	}

	if response[0] != 0x03 {
		panic(fmt.Sprintf("Unexpected function byte: %d", response))
	}
}

func (ub unityBridgeImpl) SendEvent(eventCode int64, data []byte, tag int64) {
	var b bytes.Buffer

	b.WriteByte(0x04)
	binary.Write(&b, binary.LittleEndian, eventCode)
	binary.Write(&b, binary.LittleEndian, tag)
	binary.Write(&b, binary.LittleEndian, uint16(len(data)))
	b.Write(data)

	_, err := b.WriteTo(localWritePipe)
	if err != nil {
		panic(err)
	}

	response := make([]byte, 1)
	_, err = localReadPipe.Read(response)
	if err != nil {
		panic(err)
	}

	if response[0] != 0x04 {
		panic(fmt.Sprintf("Unexpected function byte: %d", response))
	}
}

func (ub unityBridgeImpl) SendEventWithString(eventCode int64, data string,
	tag int64) {
	var b bytes.Buffer

	b.WriteByte(0x05)
	binary.Write(&b, binary.LittleEndian, eventCode)
	binary.Write(&b, binary.LittleEndian, tag)
	binary.Write(&b, binary.LittleEndian, uint16(len(data)))
	b.WriteString(data)

	_, err := b.WriteTo(localWritePipe)
	if err != nil {
		panic(err)
	}

	response := make([]byte, 1)
	_, err = localReadPipe.Read(response)
	if err != nil {
		panic(err)
	}

	if response[0] != 0x05 {
		panic(fmt.Sprintf("Unexpected function byte: %d", response))
	}
}

func (ub unityBridgeImpl) SendEventWithNumber(eventCode int64, data int64,
	tag int64) {
	var b bytes.Buffer

	b.WriteByte(0x06)
	binary.Write(&b, binary.LittleEndian, eventCode)
	binary.Write(&b, binary.LittleEndian, tag)
	binary.Write(&b, binary.LittleEndian, data)

	_, err := b.WriteTo(localWritePipe)
	if err != nil {
		panic(err)
	}

	response := make([]byte, 1)
	_, err = localReadPipe.Read(response)
	if err != nil {
		panic(err)
	}
}

func (ub unityBridgeImpl) SetEventCallback(eventCode int64,
	handler EventCallbackHandler) {
	var b bytes.Buffer

	b.WriteByte(0x07)
	binary.Write(&b, binary.LittleEndian, eventCode)
	binary.Write(&b, binary.LittleEndian, handler != nil)

	_, err := b.WriteTo(localWritePipe)
	if err != nil {
		panic(err)
	}

	response := make([]byte, 1)
	_, err = localReadPipe.Read(response)
	if err != nil {
		panic(err)
	}

	if response[0] != 0x07 {
		panic(fmt.Sprintf("Unexpected function byte: %d", response))
	}
}

func (ub unityBridgeImpl) GetSecurityKeyByKeyChainIndex(index int64) string {
	var b bytes.Buffer

	b.WriteByte(0x08)
	binary.Write(&b, binary.LittleEndian, index)

	_, err := b.WriteTo(localWritePipe)
	if err != nil {
		panic(err)
	}

	response := make([]byte, 3)
	_, err = localReadPipe.Read(response)
	if err != nil {
		panic(err)
	}

	if response[0] != 0x08 {
		panic(fmt.Sprintf("Unexpected function byte: %d", response))
	}

	var length uint16
	binary.Read(localReadPipe, binary.LittleEndian, &length)

	data := make([]byte, length)
	localReadPipe.Read(data)

	return string(data)
}

func getWinePath() (string, error) {
	return exec.LookPath("wine")
}

func getDLLHostPath() (string, error) {
	dllHostPath, err := exec.LookPath(dllHostExe)
	if err != nil {
		// Try current directory.
		dllHostPath, err = exec.LookPath("./" + dllHostExe)
		if err != nil {
			return "", err
		}
	}

	peFile, err := pe.Open(dllHostPath)
	if err != nil {
		return "", fmt.Errorf("%q does not look like a Windows executable: %w",
			dllHostPath, err)
	}
	peFile.Close()

	return dllHostPath, nil
}

func startDllHost(winePath, dllHostPath string, remoteReadPipe,
	remoteWritePipe *os.File) error {
	argv := []string{
		winePath,
		dllHostPath,
		fmt.Sprintf("%d", getFd(remoteReadPipe)),
		fmt.Sprintf("%d", getFd(remoteWritePipe)),
	}

	// Disable close on exec for the pipes.
	disableCloseOnExec(remoteReadPipe)
	disableCloseOnExec(remoteWritePipe)

	_, err := syscall.ForkExec(winePath, argv,
		&syscall.ProcAttr{
			Files: []uintptr{
				getFd(os.Stdin),
				getFd(os.Stdout),
				getFd(os.Stderr),
			},
			Env: []string{"WINEDEBUG=-all"},
			Sys: &syscall.SysProcAttr{
				Foreground: false,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("error executing windows program: %w", err)
	}

	remoteReadPipe.Close()
	remoteWritePipe.Close()

	return nil
}

func disableCloseOnExec(file *os.File) {
	_, _, err := syscall.Syscall(syscall.SYS_FCNTL, getFd(file),
		syscall.F_SETFD, 0)
	if err != syscall.Errno(0) {
		panic(fmt.Sprintf("Error disabling close on exec: %s", err))
	}
}

func getFd(file *os.File) uintptr {
	rawConn, err := file.SyscallConn()
	if err != nil {
		panic(fmt.Sprintf("Error getting raw connection "+
			"for file: %s", err))
	}

	var fileFd uintptr
	err = rawConn.Control(func(fd uintptr) {
		fileFd = fd
	})
	if err != nil {
		panic(fmt.Sprintf("Error controlling raw "+
			"connection: %s", err))
	}

	return fileFd
}

func loop() {
	var b bytes.Buffer

	for {
		_, err := b.ReadFrom(localReadPipe)
		if err != nil {
			panic(err)
		}

		var eventCode int64
		binary.Read(&b, binary.LittleEndian, &eventCode)

		var tag int64
		binary.Read(&b, binary.LittleEndian, &tag)

		var length uint16
		binary.Read(&b, binary.LittleEndian, &length)

		data := make([]byte, length)
		b.Read(data)

		runEventCallback(eventCode, data, tag)
	}
}
