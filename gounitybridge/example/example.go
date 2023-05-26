package main

import (
	"github.com/brunoga/robomaster2/gounitybridge"
)

func main() {
	ub := gounitybridge.Get()

	ub.Create("Robomaster", true, "./log")
	defer ub.Destroy()

	if !ub.Initialize() {
		panic("UnityBridge failed to initialize")
	}
	defer ub.Uninitialize()
}
