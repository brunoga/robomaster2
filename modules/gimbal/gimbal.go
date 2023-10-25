package gimbal

import "github.com/brunoga/robomaster2/support"

type Gimbal struct {
	logger *support.Logger
}

func New(logger *support.Logger) *Gimbal {
	return &Gimbal{
		logger,
	}
}
