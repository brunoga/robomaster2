package gimbal

import (
	"fmt"
	"sync"

	"github.com/brunoga/robomaster2/internal/robot/service"
	"github.com/brunoga/robomaster2/internal/robot/service/dji"
	"github.com/brunoga/robomaster2/support"
)

type Gimbal struct {
	logger *support.Logger
}

func New(logger *support.Logger) *Gimbal {
	return &Gimbal{
		logger,
	}
}

func (g *Gimbal) Start() {
	cc := service.DJICommandControllerInstance()

	connectionWg := sync.WaitGroup{}

	connectionWg.Add(1)
	cc.StartListeningOnKey(dji.DJIGimbalConnection, g,
		func(result *dji.DJIResult) {
			if result.Value().(bool) {
				// Enable gimbal updates.
				fmt.Println("Gimbal connection established.")
				cc.PerformAction(dji.DJIGimbalOpenAttitudeUpdates, nil)
			} else {
				fmt.Println("Gimbal connection failed.")
			}

			connectionWg.Done()
		}, false)

	connectionWg.Wait()
}

func (g *Gimbal) MoveToAbsoluteAngle(angle, axis int16, duration float32) {
	if axis == 1 {
		service.DJICommandControllerInstance().PerformActionWithParam(
			dji.DJIGimbalAngleFrontYawRotation,
			dji.NewDJIGimbalAngleRotationParamValue(
				0, angle*10, int16(duration*1000)), nil)
		return
	}

	service.DJICommandControllerInstance().PerformActionWithParam(
		dji.DJIGimbalAngleFrontPitchRotation,
		dji.NewDJIGimbalAngleRotationParamValue(
			angle*10, 0, int16(duration*1000)), nil)
}

func (g *Gimbal) MoveToRelativeAngle(angle, axis int16, duration float32) {
	time := int16(duration * 1000)
	yaw := int16(0)
	pitch := int16(angle * 10)
	if axis == 1 {
		yaw = pitch
		pitch = 0
	}

	service.DJICommandControllerInstance().PerformActionWithParam(
		dji.DJIGimbalAngleIncrementRotation,
		dji.NewDJIGimbalAngleRotationParamValue(pitch, yaw, time), nil)
}

func (g *Gimbal) Reset() {
	service.DJICommandControllerInstance().PerformAction(
		dji.DJIGimbalResetPosition, nil)
}
