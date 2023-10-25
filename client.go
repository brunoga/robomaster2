package robomaster2

import (
	"time"

	"github.com/brunoga/robomaster2/internal/robot/service"
	"github.com/brunoga/robomaster2/internal/robot/service/dji"
	"github.com/brunoga/robomaster2/internal/robot/service/unitybridge"
	"github.com/brunoga/robomaster2/modules/finder"
	"github.com/brunoga/robomaster2/support"
)

type Client struct {
	logger *support.Logger

	finder *finder.Finder
	cc     service.DJICommandController
}

func NewClient(logger *support.Logger) (*Client, error) {
	cc := service.DJICommandControllerInstance()

	return &Client{
		logger,
		finder.New(logger),
		cc,
	}, nil
}

func (c *Client) Start() error {
	ub := unitybridge.DJIUnityBridgeInstance()
	ub.Init()

	c.cc.Init()

	c.cc.StartListeningOnKey(dji.DJIAirLinkConnection, c,
		func(result *dji.DJIResult) {
			if result.Value().(bool) {
				c.logger.INFO("Connected to Robot.")
			}
		}, false)

	ip, err := c.finder.GetOrFindIP(5 * time.Second)
	if err != nil {
		return err
	}

	// Send ack.
	c.finder.SendACK()

	ub.SendEventWithoutDataOrTag(
		unitybridge.NewDJIUnityEventWithTypeAndSubType(
			unitybridge.Connection, 1))
	ub.SendEventWithString(unitybridge.NewDJIUnityEventWithTypeAndSubType(
		unitybridge.Connection, 2), ip.String(), 0)
	ub.SendEventWithNumber(unitybridge.NewDJIUnityEventWithTypeAndSubType(
		unitybridge.Connection, 3), 10607, 0)
	ub.SendEventWithoutDataOrTag(unitybridge.NewDJIUnityEventWithType(
		unitybridge.Connection))

	return nil
}

func (c *Client) Stop() {
	c.cc.UnInit()
	unitybridge.DJIUnityBridgeInstance().UnInit()
}
