package netstat

import (
	"context"

	"github.com/shirou/gopsutil/v3/net"
)

const (
	sockStream = 1 // SOCK_STREAM
	sockDgram  = 2 // SOCK_DGRAM
)

// afINET6 returns the AF_INET6 family constant for the current platform.
// Windows reports IPv6 as 23; Unix-like systems report 10 (or 30 on some BSDs).
// We accept both here to stay robust against future OS changes.
func afINET6(family uint32) bool {
	return family == 23 || family == 10 || family == 30
}

// gopsutilSnapshot reads TCP+UDP, IPv4+IPv6 connections via gopsutil and maps
// them to our ConnInfo schema. Shared by all platform implementations.
func gopsutilSnapshot(ctx context.Context) ([]ConnInfo, error) {
	rows, err := net.ConnectionsWithContext(ctx, "all")
	if err != nil {
		return nil, err
	}
	out := make([]ConnInfo, 0, len(rows))
	for _, r := range rows {
		ci := ConnInfo{
			Protocol:   mapProto(r.Type, r.Family),
			LocalAddr:  ipString(r.Laddr.IP),
			LocalPort:  portU16(r.Laddr.Port),
			RemoteAddr: ipString(r.Raddr.IP),
			RemotePort: portU16(r.Raddr.Port),
			State:      mapState(r.Status),
			PID:        uint32(r.Pid),
		}
		ci.Key = Key(ci.Protocol, ci.LocalAddr, ci.LocalPort, ci.RemoteAddr, ci.RemotePort, ci.PID)
		out = append(out, ci)
	}
	return out, nil
}

// mapProto combines gopsutil's (Type, Family) into our string protocol tag.
func mapProto(sockType, family uint32) Protocol {
	v6 := afINET6(family)
	switch sockType {
	case sockStream:
		if v6 {
			return ProtocolTCP6
		}
		return ProtocolTCP4
	case sockDgram:
		if v6 {
			return ProtocolUDP6
		}
		return ProtocolUDP4
	}
	return ProtocolTCP4
}

// mapState maps gopsutil's string state into our State enum.
// UDP entries commonly come back as empty/NONE — we render them as LISTEN to
// match the legacy behavior of the go-netstat implementation.
func mapState(s string) State {
	switch s {
	case "ESTABLISHED":
		return StateEstablished
	case "SYN_SENT":
		return StateSynSent
	case "SYN_RECEIVED":
		return StateSynReceived
	case "FIN_WAIT1":
		return StateFinWait1
	case "FIN_WAIT2":
		return StateFinWait2
	case "TIME_WAIT":
		return StateTimeWait
	case "CLOSE":
		return StateClosed
	case "CLOSE_WAIT":
		return StateCloseWait
	case "LISTEN":
		return StateListen
	case "CLOSING":
		return StateClosing
	case "LAST_ACK":
		return StateDeleteTCB
	case "NONE", "":
		return StateListen
	}
	return StateUnknown
}

func ipString(s string) string {
	if s == "" {
		return "0.0.0.0"
	}
	return s
}

func portU16(p uint32) uint16 {
	if p > 65535 {
		return 0
	}
	return uint16(p)
}
