package main

import (
	"fmt"
	"os"
	"time"

	"github.com/brunoga/robomaster2"
	"github.com/brunoga/robomaster2/modules/video"
	"github.com/brunoga/robomaster2/support"
)

type ExampleVideoHandler struct{}

func (h *ExampleVideoHandler) HandleVideo(img *video.RGB) {
	fmt.Println("Got frame!")
}

func main() {
	l := support.NewLogger(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	c, err := robomaster2.NewClient(l)
	if err != nil {
		panic(err)
	}

	err = c.Start()

	c.Video().AddVideoHandler(&ExampleVideoHandler{})

	if err != nil {
		l.ERROR("Error starting client: %s", err.Error())
		// TODO(bga): Exit.
	}
	defer c.Stop()

	time.Sleep(10 * time.Second)
}
