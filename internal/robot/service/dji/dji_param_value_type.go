package dji

type DJIParamValueType int

const (
	DJIParamValueTypeString DJIParamValueType = iota
	DJIParamValueTypeNumber
	DJIParamValueTypeStruct
	DJIParamValueTypeOther
)
