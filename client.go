package robomaster2

import (
	"fmt"
	"sync"
	"time"

	"github.com/brunoga/robomaster2/internal"
	"github.com/brunoga/robomaster2/internal/dji"
	"github.com/brunoga/robomaster2/internal/event"
	"github.com/brunoga/robomaster2/internal/unitybridge"
	"github.com/brunoga/robomaster2/modules/finder"
	"github.com/brunoga/robomaster2/support"
)

type Client struct {
	logger *support.Logger

	finder *finder.Finder
	cc     *internal.CommandController
}

func NewClient(logger *support.Logger) (*Client, error) {
	cc, err := internal.NewCommandController()
	if err != nil {
		return nil, err
	}

	return &Client{
		logger,
		finder.New(logger),
		cc,
	}, nil
}

func (c *Client) Start() error {
	ub := unitybridge.Get()

	ub.Create("Robomaster", true, "./log")
	if !ub.Initialize() {
		return fmt.Errorf("unable to initialize UnityBridge")
	}

	c.cc.StartListening(dji.KeyAirLinkConnection,
		func(result *dji.Result, wg *sync.WaitGroup) {
			if result.Value().(bool) {
				c.finder.SendACK()
			}

			wg.Done()
		})

	// Reset connection to defaults.
	ev := event.NewWithSubType(
		event.TypeConnection, 2)
	ub.SendEventWithString(int64(ev.Code()), "192.168.2.1", 0)

	ev = event.NewWithSubType(event.TypeConnection, 3)
	ub.SendEventWithNumber(int64(ev.Code()), int64(10607), 0)

	ev = event.New(event.TypeConnection)
	ub.SendEvent(int64(ev.Code()), nil, 0)

	ip, err := c.finder.GetOrFindIP(5 * time.Second)
	if err != nil {
		return err
	}

	ev = event.NewWithSubType(event.TypeConnection, 1)
	ub.SendEvent(int64(ev.Code()), nil, 0)

	ev = event.NewWithSubType(event.TypeConnection, 2)
	ub.SendEventWithString(int64(ev.Code()), ip.String(), 0)

	ev = event.NewWithSubType(event.TypeConnection, 3)
	ub.SendEventWithNumber(int64(ev.Code()), int64(10607), 0)

	ev = event.New(event.TypeConnection)
	ub.SendEvent(int64(ev.Code()), nil, 0)

	return nil
}

func (c *Client) Stop() error {
	return fmt.Errorf("not implemented")
}
