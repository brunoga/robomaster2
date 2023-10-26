package dji

type DJIBoolParamValue struct {
	Value bool `json:"value"`
}

func NewDJIBoolParamValue(value bool) *DJIBoolParamValue {
	return &DJIBoolParamValue{
		Value: value,
	}
}
