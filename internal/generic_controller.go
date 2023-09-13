package internal

import (
	"fmt"

	"github.com/brunoga/robomaster2/internal/event"
	"github.com/brunoga/robomaster2/internal/unitybridge"
)

type GenericController struct {
	eventHandler unitybridge.EventCallbackHandler

	eventHandlerIndexMap map[event.Type]struct{}
}

func NewGenericController(eventHandler unitybridge.EventCallbackHandler) *GenericController {
	return &GenericController{
		eventHandler,
		make(map[event.Type]struct{}),
	}
}

func (g *GenericController) StartControllingEvent(typ event.Type) error {
	b := unitybridge.Get()

	_, ok := g.eventHandlerIndexMap[typ]
	if ok {
		return fmt.Errorf("event type %s is already being controlled",
			event.TypeName(typ))
	}

	b.SetEventCallback(uint64(typ), g.eventHandler)

	g.eventHandlerIndexMap[typ] = struct{}{}

	return nil
}

func (g *GenericController) StopControllingEvent(typ event.Type) error {
	b := unitybridge.Get()

	_, ok := g.eventHandlerIndexMap[typ]
	if !ok {
		return fmt.Errorf("event type %s is not being controlled",
			event.TypeName(typ))
	}

	b.SetEventCallback(uint64(typ), nil)

	delete(g.eventHandlerIndexMap, typ)

	return nil
}
