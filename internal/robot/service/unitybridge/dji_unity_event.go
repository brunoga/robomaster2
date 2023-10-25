package unitybridge

type DJIUnityEvent struct {
	typ    DJIUnityEventType
	subTyp uint32
}

func NewDJIUnityEvent(code uint64) *DJIUnityEvent {
	return &DJIUnityEvent{
		typ:    DJIUnityEventType(code >> 32),
		subTyp: uint32(code & uint64(^uint(0))),
	}
}

func NewDJIUnityEventWithType(typ DJIUnityEventType) *DJIUnityEvent {
	return &DJIUnityEvent{
		typ: typ,
	}
}

func NewDJIUnityEventWithTypeAndSubType(typ DJIUnityEventType, subTyp uint32) *DJIUnityEvent {
	return &DJIUnityEvent{
		typ:    typ,
		subTyp: subTyp,
	}
}

func NewDJIUnityEventZero() *DJIUnityEvent {
	return &DJIUnityEvent{}
}

func (e *DJIUnityEvent) Type() DJIUnityEventType {
	return e.typ
}

func (e *DJIUnityEvent) SubType() uint32 {
	return e.subTyp
}

func (e *DJIUnityEvent) GetCode() uint64 {
	return uint64(e.typ)<<32 | uint64(e.subTyp)
}

func (e *DJIUnityEvent) Reset(typ DJIUnityEventType, subTyp uint32) {
	e.typ = typ
	e.subTyp = subTyp
}

func (e *DJIUnityEvent) ResetSubType(subTyp uint32) {
	e.subTyp = subTyp
}
