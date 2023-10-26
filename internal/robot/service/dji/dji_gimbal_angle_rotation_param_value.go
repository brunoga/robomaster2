package dji

type DJIGimbalAngleRotationParamValue struct {
	Pitch int16 `json:"pitch"`
	Yaw   int16 `json:"yaw"`
	Time  int16 `json:"time"`
}

func NewDJIGimbalAngleRotationParamValue(pitch, yaw int16, time int16) *DJIGimbalAngleRotationParamValue {
	return &DJIGimbalAngleRotationParamValue{
		Pitch: pitch,
		Yaw:   yaw,
		Time:  time,
	}
}
