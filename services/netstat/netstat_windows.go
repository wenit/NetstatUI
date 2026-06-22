//go:build windows

package netstat

import (
	"context"
	"fmt"

	gns "github.com/cakturk/go-netstat/netstat"
)

type windowsProvider struct{}

func NewWindowsProvider() Provider { return &windowsProvider{} }

func (w *windowsProvider) Snapshot(ctx context.Context) ([]ConnInfo, error) {
	var out []ConnInfo

	// TCP4 — uses GetTcpTable2 internally (newer than the manual GetExtendedTcpTable path)
	socks, err := gns.TCPSocks(gns.NoopFilter)
	if err != nil {
		return nil, fmt.Errorf("tcp4: %w", err)
	}
	for _, s := range socks {
		out = append(out, sockToConnInfo(ProtocolTCP4, &s))
	}

	// TCP6 — uses GetTcp6Table2
	socks, err = gns.TCP6Socks(gns.NoopFilter)
	if err != nil {
		return nil, fmt.Errorf("tcp6: %w", err)
	}
	for _, s := range socks {
		out = append(out, sockToConnInfo(ProtocolTCP6, &s))
	}

	// UDP4 — uses GetExtendedUdpTable (owner PID)
	socks, err = gns.UDPSocks(gns.NoopFilter)
	if err != nil {
		return nil, fmt.Errorf("udp4: %w", err)
	}
	for _, s := range socks {
		out = append(out, sockToConnInfo(ProtocolUDP4, &s))
	}

	// UDP6
	socks, err = gns.UDP6Socks(gns.NoopFilter)
	if err != nil {
		return nil, fmt.Errorf("udp6: %w", err)
	}
	for _, s := range socks {
		out = append(out, sockToConnInfo(ProtocolUDP6, &s))
	}

	return out, nil
}

func sockToConnInfo(p Protocol, s *gns.SockTabEntry) ConnInfo {
	c := ConnInfo{
		Protocol:   p,
		LocalAddr:  sockAddrToString(s.LocalAddr),
		LocalPort:  s.LocalAddr.Port,
		RemoteAddr: sockAddrToString(s.RemoteAddr),
		RemotePort: s.RemoteAddr.Port,
		State:      stateFromSock(p, s),
		PID:        pidPtr(s.Process),
	}
	if s.Process != nil {
		c.ProcessName = s.Process.Name
	}
	c.Key = Key(c.Protocol, c.LocalAddr, c.LocalPort, c.RemoteAddr, c.RemotePort, c.PID)
	return c
}

func sockAddrToString(sa *gns.SockAddr) string {
	if sa == nil || sa.IP == nil {
		return "0.0.0.0"
	}
	return sa.IP.String()
}

func pidPtr(p *gns.Process) uint32 {
	if p == nil {
		return 0
	}
	return uint32(p.Pid)
}

// stateFromSock maps the library's SkState to our string State.
// For UDP the library always returns "CLOSE" (empty string), we override to LISTEN.
func stateFromSock(p Protocol, s *gns.SockTabEntry) State {
	if p == ProtocolUDP4 || p == ProtocolUDP6 {
		return StateListen
	}
	return State(s.State.String())
}
