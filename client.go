package robomaster2

import (
	"time"

	"github.com/brunoga/robomaster2/internal/robot/service"
	"github.com/brunoga/robomaster2/internal/robot/service/dji"
	"github.com/brunoga/robomaster2/internal/robot/service/unitybridge"
	"github.com/brunoga/robomaster2/modules/chassis"
	"github.com/brunoga/robomaster2/modules/finder"
	"github.com/brunoga/robomaster2/modules/gimbal"
	"github.com/brunoga/robomaster2/modules/video"
	"github.com/brunoga/robomaster2/support"
)

type Client struct {
	logger *support.Logger

	finder *finder.Finder
	cc     service.DJICommandController

	chassis *chassis.Chassis
	gimbal  *gimbal.Gimbal
	video   *video.Video
}

func NewClient(logger *support.Logger) (*Client, error) {
	cc := service.DJICommandControllerInstance()

	return &Client{
		logger,
		finder.New(logger),
		cc,
		chassis.New(logger),
		gimbal.New(logger),
		video.New(logger),
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

	c.Gimbal().Start()
	//c.Chassis().Start()

	return nil
}

func (c *Client) Stop() {
	c.video.Stop()
	c.cc.UnInit()
	unitybridge.DJIUnityBridgeInstance().UnInit()
}

func (c *Client) Chassis() *chassis.Chassis {
	return c.chassis
}

func (c *Client) Gimbal() *gimbal.Gimbal {
	return c.gimbal
}

func (c *Client) Video() *video.Video {
	return c.video
}
