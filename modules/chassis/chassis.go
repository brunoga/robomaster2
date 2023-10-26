package chassis

import (
	"fmt"
	"sync"

	"github.com/brunoga/robomaster2/internal/robot/service"
	"github.com/brunoga/robomaster2/internal/robot/service/dji"
	"github.com/brunoga/robomaster2/support"
)

type Chassis struct {
	logger *support.Logger
}

func New(logger *support.Logger) *Chassis {
	return &Chassis{
		logger,
	}
}

func (c *Chassis) Start() {
	cc := service.DJICommandControllerInstance()

	connectionWg := sync.WaitGroup{}

	connectionWg.Add(1)
	cc.StartListeningOnKey(dji.DJIRobomasterSystemConnection, c,
		func(result *dji.DJIResult) {
			if result.Value().(bool) {
				fmt.Println("Chassis connection established.")
				cc.PerformAction(dji.DJIRobomasterOpenChassisSpeedUpdates, nil)
			} else {
				fmt.Println("Chassis connection failed.")
			}

			connectionWg.Done()
		}, false)

	connectionWg.Wait()
}
