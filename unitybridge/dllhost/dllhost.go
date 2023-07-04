//go:build windows && amd64

package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/brunoga/robomaster2/unitybridge"
)

var (
	// Command line flags.
	readFd  = flag.Int("read-fd", -1, "file descriptor to read from")
	writeFd = flag.Int("write-fd", -1, "file descriptor to write to")

	catchAll *CatchAllHandler
)

func main() {
	flag.Parse()

	if *readFd < 0 || *writeFd < 0 {
		fmt.Fprintln(os.Stderr, "Flags -read-fd and -write-fd must be "+
			"provided and non-negative")
		os.Exit(1)
	}

	files, err := fdsToFiles([]int{*readFd, *writeFd})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting file descriptors to "+
			"files: %s\n", err)
		os.Exit(1)
	}

	catchAll = &CatchAllHandler{writeFile: files[1]}

	err = loop(files[0], files[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in loop: %s\n", err)
		os.Exit(1)
	}
}

func fdsToFiles(fds []int) ([]*os.File, error) {
	ntDll, err := syscall.LoadDLL("ntdll.dll")
	if err != nil {
		return nil, fmt.Errorf("error loading ntdll.dll: %s", err)
	}
	defer ntDll.Release()

	wineServerFdToHandleProc, err := ntDll.FindProc(
		"wine_server_fd_to_handle")
	if err != nil {
		return nil, fmt.Errorf(
			"error finding wine_server_fd_to_handle: %w", err)
	}

	files := make([]*os.File, len(fds))
	for i, fd := range fds {
		file := fdToFile(wineServerFdToHandleProc, uintptr(fd),
			syscall.GENERIC_READ|syscall.GENERIC_WRITE,
			fmt.Sprintf("linux%d", i))
		files[i] = file
	}

	return files, nil
}

func fdToFile(proc *syscall.Proc, fd uintptr, flags uintptr,
	name string) *os.File {
	var fdHandle uintptr
	ntStatus, _, _ := proc.Call(fd, flags|syscall.SYNCHRONIZE,
		2 /*OBJ_INHERIT*/, uintptr(unsafe.Pointer(&fdHandle)))
	if ntStatus != 0 {
		panic(ntStatus)
	}

	return os.NewFile(fdHandle, name)
}

func loop(readFile, writeFile *os.File) error {
	var b bytes.Buffer
	for {
		_, err := b.ReadFrom(readFile)
		if err != nil {
			return err
		}

		function, err := b.ReadByte()
		if err != nil {
			return err
		}

		var length uint16
		if err := binary.Read(&b, binary.LittleEndian, &length); err != nil {
			return err
		}

		data := make([]byte, length)
		n, err := b.Read(data)
		if err != nil {
			return err
		}
		if n != int(length) {
			return err
		}

		process(writeFile, function, data)
	}

	return nil
}

func process(writeFile *os.File, function byte, data []byte) {
	var b bytes.Buffer

	b.WriteByte(function)

	switch function {
	case 0x00:
		runCreateUnityBridge(data, &b)
	case 0x01:
		runDestroyUnityBridge(data, &b)
	case 0x02:
		runInitializeUnityBridge(data, &b)
	case 0x03:
		runUnitializeUnityBridge(data, &b)
	case 0x04:
		runUnitySendEvent(data, &b)
	case 0x05:
		runUnitySendEventWithString(data, &b)
	case 0x06:
		runUnitySendEventWithNumber(data, &b)
	case 0x07:
		runUnitySetEventCallback(data, &b)
	case 0x08:
		runGetSecurityKeyByKeyChainIndex(data, &b)
	}

	_, err := b.WriteTo(writeFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file: %s\n", err)
		return
	}
}

func runCreateUnityBridge(data []byte, b *bytes.Buffer) {
	debuggable := data[0] != 0
	length := binary.LittleEndian.Uint16(data[1:3])
	name := string(data[3 : 3+length])
	length = binary.LittleEndian.Uint16(data[3+length : 3+length+2])
	logPath := string(data[3+length+2:])

	unitybridge.Get().Create(name, debuggable, logPath)
}

func runDestroyUnityBridge(data []byte, b *bytes.Buffer) {
	unitybridge.Get().Destroy()
}

func runInitializeUnityBridge(data []byte, b *bytes.Buffer) {
	binary.Write(b, binary.LittleEndian, uint16(1))

	if unitybridge.Get().Initialize() {
		b.WriteByte(0x01)
	} else {
		b.WriteByte(0x00)
	}
}

func runUnitializeUnityBridge(data []byte, b *bytes.Buffer) {
	unitybridge.Get().Uninitialize()
}

func runUnitySendEvent(data []byte, b *bytes.Buffer) {
	eventCode := binary.LittleEndian.Uint64(data[0:8])
	tag := binary.LittleEndian.Uint64(data[8:16])
	length := binary.LittleEndian.Uint16(data[16:18])
	data2 := data[18 : 18+length]

	unitybridge.Get().SendEvent(int64(eventCode), data2, int64(tag))
}

func runUnitySendEventWithString(data []byte, b *bytes.Buffer) {
	eventCode := binary.LittleEndian.Uint64(data[0:8])
	tag := binary.LittleEndian.Uint64(data[8:16])
	length := binary.LittleEndian.Uint16(data[16:18])
	data2 := string(data[18 : 18+length])

	unitybridge.Get().SendEventWithString(int64(eventCode), data2, int64(tag))
}

func runUnitySendEventWithNumber(data []byte, b *bytes.Buffer) {
	eventCode := binary.LittleEndian.Uint64(data[0:8])
	tag := binary.LittleEndian.Uint64(data[8:16])
	data2 := binary.LittleEndian.Uint64(data[16:24])

	unitybridge.Get().SendEventWithNumber(int64(eventCode), int64(data2), int64(tag))
}

func runUnitySetEventCallback(data []byte, b *bytes.Buffer) {
	eventCode := binary.LittleEndian.Uint64(data[0:8])
	add := data[8] != 0

	if add {
		unitybridge.Get().SetEventCallback(int64(eventCode), catchAll)
	} else {
		unitybridge.Get().SetEventCallback(int64(eventCode), nil)
	}
}

func runGetSecurityKeyByKeyChainIndex(data []byte, b *bytes.Buffer) {
	index := binary.LittleEndian.Uint64(data[0:8])

	key := unitybridge.Get().GetSecurityKeyByKeyChainIndex(int64(index))

	binary.Write(b, binary.LittleEndian, uint16(len(key)))
	b.WriteString(key)
}
