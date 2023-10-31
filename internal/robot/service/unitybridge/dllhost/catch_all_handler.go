package main

import (
	"bytes"
	"encoding/binary"
	"os"
)

type CatchAllHandler struct {
	eventFile *os.File
}

func (c *CatchAllHandler) HandleEventCallback(eventCode int64, data []byte, tag int64) {
	var b bytes.Buffer

	// TOOD(bga): Add error checking.
	binary.Write(&b, binary.LittleEndian, eventCode)
	binary.Write(&b, binary.LittleEndian, tag)
	binary.Write(&b, binary.LittleEndian, uint16(len(data)))
	b.Write(data)

	b.WriteTo(c.eventFile)
}
