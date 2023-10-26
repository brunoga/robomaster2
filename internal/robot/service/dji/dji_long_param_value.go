package dji

type DJILongParamValue struct {
	Value int64 `json:"value"`
}

func NewDJILongParamValue(value int64) *DJILongParamValue {
	return &DJILongParamValue{
		Value: value,
	}
}
