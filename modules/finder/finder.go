package finder

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/brunoga/robomaster2/support"
)

const (
	ipBroadcastAddrPort = ":45678"
	listenerRemotePort  = ":56789"
)

// Finder provides an interface for finding a robot broadcasting its ip in
// the network.
type Finder struct {
	l           *support.Logger
	m           sync.Mutex
	ip          net.IP
	remoteAppID uint64
}

// New returns a Finder instance with no associated ip.
func New(l *support.Logger) *Finder {
	return &Finder{
		l,
		sync.Mutex{},
		nil,
		0,
	}
}

// GetOrFindIP returns the ip of a robot if it is already know or tries to
// detect a robot broadcasting its ip in the network. The search will go on
// until a robot is detected or a timeout happens. Returns the robot ip and a
// nil error on success and a non-nil error on failure.
func (f *Finder) GetOrFindIP(timeout time.Duration) (net.IP, error) {
	f.m.Lock()
	defer f.m.Unlock()

	if f.ip == nil {
		ip, remoteAppID, err := findRobotIP(timeout)
		if err != nil {
			return nil, fmt.Errorf("error finding robot ip: %w", err)
		}

		f.l.INFO("Detected robot with ip %s", ip.String())

		f.ip = ip
		f.remoteAppID = remoteAppID
	}

	return f.ip, nil
}

// SetIP forces the associated ip to be the given one. Useful for when
// connecting to a robot with a known ip.
func (f *Finder) SetIP(ip net.IP) {
	f.m.Lock()
	defer f.m.Unlock()

	f.ip = ip
}

func (f *Finder) SendACK() {
	f.l.INFO("Sending ACK to %s%s.\n", f.ip.String(), listenerRemotePort)

	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffer, f.remoteAppID)

	udpAddr, err := net.ResolveUDPAddr("udp4", f.ip.String()+listenerRemotePort)
	if err != nil {
		f.l.ERROR("Error resolving UDP address: %s", err.Error())
		return
	}

	conn, err := net.DialUDP("udp4", nil, udpAddr)
	if err != nil {
		f.l.ERROR("Error dialing UDP: %s", err.Error())
		return
	}
	defer conn.Close()

	_, err = conn.Write(buffer)
	if err != nil {
		f.l.ERROR("Error writing to UDP: %s", err.Error())
		return
	}
}

func findRobotIP(timeout time.Duration) (net.IP, uint64, error) {
	packetConn, err := net.ListenPacket("udp4", ipBroadcastAddrPort)
	if err != nil {
		return nil, 0, fmt.Errorf("error starting packet listner: %w", err)
	}
	defer packetConn.Close()

	buf := make([]byte, 1024)

	err = packetConn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return nil, 0, fmt.Errorf("error setting deadline: %w", err)
	}

	n, addr, err := packetConn.ReadFrom(buf)
	if err != nil {
		return nil, 0, fmt.Errorf("error reading packet: %w", err)
	}

	ip, remoteAppID, err := parseAndValidateMessage(buf[:n], addr)
	if err != nil {
		return nil, 0, fmt.Errorf("error validating message: %w", err)
	}

	return ip, remoteAppID, nil
}

func parseAndValidateMessage(buf []byte, addr net.Addr) (net.IP, uint64, error) {
	broadcastMessage, err := ParseBroadcastMessageData(buf)
	if err != nil {
		return nil, 0, fmt.Errorf("error parsing broadcast message: %w", err)
	}

	if !broadcastMessage.IsPairing() {
		return nil, 0, nil
	}

	// Get IP and make sure it is IPv4
	ip := net.IP(broadcastMessage.SourceIp()).To4()
	if ip == nil {
		return nil, 0, fmt.Errorf("not an IPv4 address")
	}

	if !ip.Equal(addr.(*net.UDPAddr).IP) {
		return nil, 0, fmt.Errorf("broadcast message source does not match reported IP")
	}

	return ip, broadcastMessage.AppId(), nil
}
