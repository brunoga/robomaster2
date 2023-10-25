package main

import (
	"os"
	"time"

	"github.com/brunoga/robomaster2"
	"github.com/brunoga/robomaster2/support"
)

func main() {
	l := support.NewLogger(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	c, err := robomaster2.NewClient(l)
	if err != nil {
		panic(err)
	}

	err = c.Start()
	if err != nil {
		l.ERROR("Error starting client: %s", err.Error())
		// TODO(bga): Exit.
	}
	defer c.Stop()

	time.Sleep(10 * time.Second)
}
