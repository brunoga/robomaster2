package main

import "github.com/brunoga/robomaster2/unitybridge"

func main() {
	ub := unitybridge.Get()

	ub.Create("Robomaster", true, "./log")
	defer ub.Destroy()

	if !ub.Initialize() {
		panic("UnityBridge failed to initialize")
	}
	defer ub.Uninitialize()
}
