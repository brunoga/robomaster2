package dji

type DJIFloatParamValue struct {
	Value float32 `json:"value"`
}

func NewDJIFloatParamValue(value float32) *DJIFloatParamValue {
	return &DJIFloatParamValue{
		Value: value,
	}
}
