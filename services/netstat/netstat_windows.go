//go:build windows

package netstat

import (
	"context"
	"fmt"
	"net"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	afinet  = 2
	afinet6 = 23

	tcpTableOwnerPidAll = 4
	udpTableOwnerPid    = 1

	errorInsufficientBuffer = 122
	noError                 = 0
)

var (
	modIphlpapi                = windows.NewLazySystemDLL("iphlpapi.dll")
	procGetExtendedTcpTable    = modIphlpapi.NewProc("GetExtendedTcpTable")
	procGetExtendedUdpTable    = modIphlpapi.NewProc("GetExtendedUdpTable")
)

type mibTcpRowOwnerPid struct {
	State      uint32
	LocalAddr  uint32
	LocalPort  uint32
	RemoteAddr uint32
	RemotePort uint32
	OwningPid  uint32
}

type mibTcp6RowOwnerPid struct {
	LocalAddr       [16]byte
	LocalScopeId    uint32
	LocalPort       uint32
	RemoteAddr      [16]byte
	RemoteScopeId   uint32
	RemotePort      uint32
	State           uint32
	OwningPid       uint32
}

type mibUdpRowOwnerPid struct {
	LocalAddr uint32
	LocalPort uint32
	OwningPid uint32
}

type mibUdp6RowOwnerPid struct {
	LocalAddr    [16]byte
	LocalScopeId uint32
	LocalPort    uint32
	OwningPid    uint32
}

type windowsProvider struct{}

func NewWindowsProvider() Provider { return &windowsProvider{} }

func (w *windowsProvider) Snapshot(ctx context.Context) ([]ConnInfo, error) {
	var out []ConnInfo
	if v, err := w.tcp4(); err != nil {
		return nil, fmt.Errorf("tcp4: %w", err)
	} else {
		out = append(out, v...)
	}
	if v, err := w.tcp6(); err != nil {
		return nil, fmt.Errorf("tcp6: %w", err)
	} else {
		out = append(out, v...)
	}
	if v, err := w.udp4(); err != nil {
		return nil, fmt.Errorf("udp4: %w", err)
	} else {
		out = append(out, v...)
	}
	if v, err := w.udp6(); err != nil {
		return nil, fmt.Errorf("udp6: %w", err)
	} else {
		out = append(out, v...)
	}
	return out, nil
}

func (w *windowsProvider) tcp4() ([]ConnInfo, error) {
	buf, err := getExtendedTable(procGetExtendedTcpTable, afinet, tcpTableOwnerPidAll)
	if err != nil {
		return nil, err
	}
	if len(buf) < 4 {
		return nil, nil
	}
	num := *(*uint32)(unsafe.Pointer(&buf[0]))
	rowSize := int(unsafe.Sizeof(mibTcpRowOwnerPid{}))
	rows := int(num)
	if 4+rows*rowSize > len(buf) {
		rows = (len(buf) - 4) / rowSize
	}
	out := make([]ConnInfo, 0, rows)
	for i := 0; i < rows; i++ {
		base := 4 + i*rowSize
		row := (*mibTcpRowOwnerPid)(unsafe.Pointer(&buf[base]))
		c := ConnInfo{
			Protocol:   ProtocolTCP4,
			LocalAddr:  ipv4String(row.LocalAddr),
			LocalPort:  ntohsLow(row.LocalPort),
			RemoteAddr: ipv4String(row.RemoteAddr),
			RemotePort: ntohsLow(row.RemotePort),
			State:      tcpState(row.State),
			PID:        row.OwningPid,
		}
		c.Key = Key(c.Protocol, c.LocalAddr, c.LocalPort, c.RemoteAddr, c.RemotePort, c.PID)
		out = append(out, c)
	}
	return out, nil
}

func (w *windowsProvider) tcp6() ([]ConnInfo, error) {
	buf, err := getExtendedTable(procGetExtendedTcpTable, afinet6, tcpTableOwnerPidAll)
	if err != nil {
		return nil, err
	}
	if len(buf) < 4 {
		return nil, nil
	}
	num := *(*uint32)(unsafe.Pointer(&buf[0]))
	rowSize := int(unsafe.Sizeof(mibTcp6RowOwnerPid{}))
	rows := int(num)
	if 4+rows*rowSize > len(buf) {
		rows = (len(buf) - 4) / rowSize
	}
	out := make([]ConnInfo, 0, rows)
	for i := 0; i < rows; i++ {
		base := 4 + i*rowSize
		row := (*mibTcp6RowOwnerPid)(unsafe.Pointer(&buf[base]))
		c := ConnInfo{
			Protocol:   ProtocolTCP6,
			LocalAddr:  ipv6String(row.LocalAddr),
			LocalPort:  ntohsLow(row.LocalPort),
			RemoteAddr: ipv6String(row.RemoteAddr),
			RemotePort: ntohsLow(row.RemotePort),
			State:      tcpState(row.State),
			PID:        row.OwningPid,
		}
		c.Key = Key(c.Protocol, c.LocalAddr, c.LocalPort, c.RemoteAddr, c.RemotePort, c.PID)
		out = append(out, c)
	}
	return out, nil
}

func (w *windowsProvider) udp4() ([]ConnInfo, error) {
	buf, err := getExtendedTable(procGetExtendedUdpTable, afinet, udpTableOwnerPid)
	if err != nil {
		return nil, err
	}
	if len(buf) < 4 {
		return nil, nil
	}
	num := *(*uint32)(unsafe.Pointer(&buf[0]))
	rowSize := int(unsafe.Sizeof(mibUdpRowOwnerPid{}))
	rows := int(num)
	if 4+rows*rowSize > len(buf) {
		rows = (len(buf) - 4) / rowSize
	}
	out := make([]ConnInfo, 0, rows)
	for i := 0; i < rows; i++ {
		base := 4 + i*rowSize
		row := (*mibUdpRowOwnerPid)(unsafe.Pointer(&buf[base]))
		c := ConnInfo{
			Protocol:  ProtocolUDP4,
			LocalAddr: ipv4String(row.LocalAddr),
			LocalPort: ntohsLow(row.LocalPort),
			State:     StateListen,
			PID:       row.OwningPid,
		}
		c.Key = Key(c.Protocol, c.LocalAddr, c.LocalPort, c.RemoteAddr, c.RemotePort, c.PID)
		out = append(out, c)
	}
	return out, nil
}

func (w *windowsProvider) udp6() ([]ConnInfo, error) {
	buf, err := getExtendedTable(procGetExtendedUdpTable, afinet6, udpTableOwnerPid)
	if err != nil {
		return nil, err
	}
	if len(buf) < 4 {
		return nil, nil
	}
	num := *(*uint32)(unsafe.Pointer(&buf[0]))
	rowSize := int(unsafe.Sizeof(mibUdp6RowOwnerPid{}))
	rows := int(num)
	if 4+rows*rowSize > len(buf) {
		rows = (len(buf) - 4) / rowSize
	}
	out := make([]ConnInfo, 0, rows)
	for i := 0; i < rows; i++ {
		base := 4 + i*rowSize
		row := (*mibUdp6RowOwnerPid)(unsafe.Pointer(&buf[base]))
		c := ConnInfo{
			Protocol:  ProtocolUDP6,
			LocalAddr: ipv6String(row.LocalAddr),
			LocalPort: ntohsLow(row.LocalPort),
			State:     StateListen,
			PID:       row.OwningPid,
		}
		c.Key = Key(c.Protocol, c.LocalAddr, c.LocalPort, c.RemoteAddr, c.RemotePort, c.PID)
		out = append(out, c)
	}
	return out, nil
}

func getExtendedTable(proc *windows.LazyProc, af uint32, tableClass uint32) ([]byte, error) {
	var size uint32 = 0
	r1, _, _ := proc.Call(0, uintptr(unsafe.Pointer(&size)), 1, uintptr(af), uintptr(tableClass), 0)
	if r1 != errorInsufficientBuffer && r1 != noError {
		if r1 == noError {
			return nil, nil
		}
		return nil, fmt.Errorf("query size failed: %d", r1)
	}
	if size == 0 {
		return nil, nil
	}
	buf := make([]byte, size)
	r1, _, _ = proc.Call(uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&size)), 1, uintptr(af), uintptr(tableClass), 0)
	if r1 != noError {
		return nil, fmt.Errorf("query table failed: %d", r1)
	}
	return buf, nil
}

func ipv4String(v uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(v&0xFF),
		byte((v>>8)&0xFF),
		byte((v>>16)&0xFF),
		byte((v>>24)&0xFF),
	)
}

func ipv6String(a [16]byte) string {
	return net.IP(a[:]).String()
}

func ntohsLow(v uint32) uint16 {
	return uint16(((v & 0xFF) << 8) | ((v >> 8) & 0xFF))
}

func tcpState(s uint32) State {
	switch s {
	case 1:
		return StateClosed
	case 2:
		return StateListen
	case 3:
		return StateSynSent
	case 4:
		return StateSynReceived
	case 5:
		return StateEstablished
	case 6:
		return StateFinWait1
	case 7:
		return StateFinWait2
	case 8:
		return StateCloseWait
	case 9:
		return StateClosing
	case 10:
		return StateTimeWait
	case 11:
		return StateDeleteTCB
	case 12:
		return StateBound
	default:
		return StateUnknown
	}
}
