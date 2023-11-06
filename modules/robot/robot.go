package robot

import (
	"github.com/brunoga/robomaster2/internal/robot/service"
	"github.com/brunoga/robomaster2/internal/robot/service/dji"
	"github.com/brunoga/robomaster2/support"
)

type Robot struct {
	logger *support.Logger
}

func NewRobot(logger *support.Logger) *Robot {
	return &Robot{
		logger,
	}
}

func (r *Robot) SendPose(leftVertical, leftHorizontal, rightVertical,
	rightHorizontal float32, isLeftCtl, isRightCtl bool, ctrlMode uint64) {
	leftVerticalVal := int64(lerp(-660, 660, (leftVertical+1)*0.5) + 1024)
	leftHorizontalVal := uint64(lerp(-660, 660, (leftHorizontal+1)*0.5) + 1024)
	leftCtlVal := uint64(0)
	if isLeftCtl {
		leftCtlVal = 1
	}

	rightVerticalVal := uint64(lerp(-660, 660, (rightVertical+1)*0.5) + 1024)
	rightHorizontalVal := uint64(lerp(-660, 660, (rightHorizontal+1)*0.5) + 1024)
	rightCtlVal := uint64(0)
	if isRightCtl {
		rightCtlVal = 1
	}

	commandValue :=
		uint64(leftVerticalVal) |
			leftHorizontalVal<<11 |
			rightVerticalVal<<22 |
			rightHorizontalVal<<33 |
			leftCtlVal<<44 |
			rightCtlVal<<45 |
			ctrlMode<<46

	cc := service.DJICommandControllerInstance()

	cc.DirectSendValue(dji.DJIMainControllerVirtualStick, int64(commandValue))
}

func lerp(a, b, t float32) float32 {
	return a*(1-t) + b*t
}
