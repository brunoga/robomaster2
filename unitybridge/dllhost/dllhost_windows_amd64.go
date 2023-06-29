package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

const (
	numFds = 2
)

func main() {
	ntDll, err := syscall.LoadDLL("ntdll.dll")
	if err != nil {
		panic(fmt.Sprintf("error loading ntdll.dll: %s", err))
	}
	defer ntDll.Release()

	wineServerFdToHandleProc, err := ntDll.FindProc(
		"wine_server_fd_to_handle")
	if err != nil {
		panic(fmt.Sprintf(
			"error finding wine_server_fd_to_handle: %w", err))
	}

	if len(os.Args) != numFds {
		panic(fmt.Sprintf("wrong number of file descriptors: "+
			"expected %d, got %d", numFds, len(os.Args)))
	}

	files := make([]*os.File, len(os.Args))
	for i, fdArg := range fdArgs {
		fd, err := strconv.Atoi(fdArg)
		if err != nil {
			return nil, fmt.Errorf("error converting file "+
				"descriptor to integer: %s", err)
		}

		file := fdToFile(wineServerFdToHandleProc, uintptr(fd),
			syscall.GENERIC_READ|syscall.GENERIC_WRITE,
			fmt.Sprintf("linux%d", i))
		files[i] = file
	}
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
