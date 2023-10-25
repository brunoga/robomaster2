package chassis

import "github.com/brunoga/robomaster2/support"

type Chassis struct {
	logger *support.Logger
}

func New(logger *support.Logger) *Chassis {
	return &Chassis{
		logger,
	}
}
