package service

import (
	"sync"

	"github.com/brunoga/robomaster2/internal/robot/service/dji"
)

type DJIListener struct {
	Name     string
	Callback func(result *dji.DJIResult)
}

func NewDJIListener(name string, callback func(result *dji.DJIResult)) *DJIListener {
	return &DJIListener{
		Name:     name,
		Callback: callback,
	}
}

type CallbackType int

const (
	Getter CallbackType = iota
	Setter
)

type CallbackDictionary struct {
	mList    [3]map[uint32]func(result *dji.DJIResult)
	dicCount [3]uint32
	mu       sync.Mutex
}

func NewCallbackDictionary() *CallbackDictionary {
	cd := &CallbackDictionary{}
	for i := 0; i < 3; i++ {
		cd.mList[i] = make(map[uint32]func(result *dji.DJIResult))
		cd.dicCount[i] = 0
	}
	return cd
}

func (cd *CallbackDictionary) AddAction(callbackType CallbackType, action func(result *dji.DJIResult)) uint32 {
	cd.mu.Lock()
	defer cd.mu.Unlock()

	cd.dicCount[callbackType]++
	cd.mList[callbackType][cd.dicCount[callbackType]] = action

	return cd.dicCount[callbackType]
}

func (cd *CallbackDictionary) SetAction(callbackType CallbackType, seq uint32, action func(result *dji.DJIResult)) {
	cd.mu.Lock()
	defer cd.mu.Unlock()

	cd.mList[callbackType][seq] = action
}

func (cd *CallbackDictionary) RemoveAction(callbackType CallbackType, seq uint32) {
	cd.mu.Lock()
	defer cd.mu.Unlock()

	delete(cd.mList[callbackType], seq)
}

func (cd *CallbackDictionary) Invoke(callbackType CallbackType, param *dji.DJIResult) {
	cd.mu.Lock()
	action, exists := cd.mList[callbackType][param.SequenceNumber()]
	cd.mu.Unlock()

	if exists && action != nil {
		action(param)
	}
}

func (cd *CallbackDictionary) Clear() {
	cd.mu.Lock()
	defer cd.mu.Unlock()

	for i := 0; i < 3; i++ {
		cd.mList[i] = make(map[uint32]func(result *dji.DJIResult))
		cd.dicCount[i] = 0
	}
}
