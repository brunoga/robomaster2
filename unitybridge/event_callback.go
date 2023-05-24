package unitybridge

import "C"

//export goEventCallback
func goEventCallback(eventCode uint64, data []byte, tag uint64) {}
