package dji

type DJIRealControlParamValue struct {
	LeftVerticalInput    float32
	LeftHorizontalInput  float32
	RightVerticalInput   float32
	RightHorizontalInput float32
	IsLeftInputTouched   bool
	IsRightInputTouched  bool
	YawFollowMode        bool
}
