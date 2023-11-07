package mobile

import (
	"github.com/brunoga/robomaster2"
	"github.com/brunoga/robomaster2/support"
)

type Client struct {
	c *robomaster2.Client
}

func NewClient() (*Client, error) {
	l := support.NewLogger(nil, nil, nil, nil)

	c, err := robomaster2.NewClient(l)
	if err != nil {
		return nil, err
	}

	return &Client{
		c,
	}, nil
}

func (c *Client) Start() error {
	return c.c.Start()
}

func (c *Client) Stop() error {
	c.c.Stop()

	return nil
}
