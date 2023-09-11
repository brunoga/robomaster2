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
		panic(err)
	}

	time.Sleep(10 * time.Second)

	err = c.Stop()
	if err != nil {
		panic(err)
	}
}
