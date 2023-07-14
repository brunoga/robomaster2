package finder

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/brunoga/robomaster2/internal/support"
)

var (
	// Expected broadcast message length.
	broadcastMessageLen = 24

	// Expected decoded broadcast message header.
	broadcastMessageHeader = []byte{90, 91}
)

// BroadcastMessage represents a broadcast message sent by a Robomaster robot.
type BroadcastMessage struct {
	isPairing bool
	sourceIp  net.IP
	sourceMac net.HardwareAddr
	appId     uint64
}

// ParseBroadcastMessageData parses the given data as a BroadcastMessage. It
// returns the associated BroadcastMessage instance pointer and a nil error on
// success and a nil BroadcastMessage and a non-nil error on failure.
func ParseBroadcastMessageData(data []byte) (*BroadcastMessage, error) {
	if len(data) != broadcastMessageLen {
		return nil, fmt.Errorf("unexpected broadcast message length")
	}

	// Decode incoming data.
	support.SimpleEncryptDecrypt(data)

	if !bytes.HasPrefix(data, broadcastMessageHeader) {
		return nil, fmt.Errorf("invalid broadcast message header")
	}

	// First byte tells us if this is a pairing message.
	isPairing := (data[2] & 1) > 0

	// Then we get the rest of the data trivially.
	sourceIp := data[6:10]
	sourceMac := data[10:16]
	appId := binary.LittleEndian.Uint64(data[16:])

	return &BroadcastMessage{
		isPairing,
		sourceIp,
		sourceMac,
		appId,
	}, nil
}

func (b *BroadcastMessage) IsPairing() bool {
	return b.isPairing
}

func (b *BroadcastMessage) SourceIp() net.IP {
	return b.sourceIp
}

func (b *BroadcastMessage) SourceMac() net.HardwareAddr {
	return b.sourceMac
}

func (b *BroadcastMessage) AppId() uint64 {
	return b.appId
}

func (b *BroadcastMessage) String() string {
	return fmt.Sprintf("IsPairing:%t, SourceIp:%s, SourceMac:%s, AppId:%d",
		b.isPairing, b.sourceIp, b.sourceMac, b.appId)
}
