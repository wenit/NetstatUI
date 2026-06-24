package netstat

type Protocol string

const (
	ProtocolTCP4 Protocol = "tcp4"
	ProtocolTCP6 Protocol = "tcp6"
	ProtocolUDP4 Protocol = "udp4"
	ProtocolUDP6 Protocol = "udp6"
)

type State string

const (
	StateUnknown      State = "UNKNOWN"
	StateClosed       State = "CLOSED"
	StateListen       State = "LISTEN"
	StateSynSent      State = "SYN_SENT"
	StateSynReceived  State = "SYN_RECEIVED"
	StateEstablished  State = "ESTABLISHED"
	StateFinWait1     State = "FIN_WAIT_1"
	StateFinWait2     State = "FIN_WAIT_2"
	StateCloseWait    State = "CLOSE_WAIT"
	StateClosing      State = "CLOSING"
	StateTimeWait     State = "TIME_WAIT"
	StateDeleteTCB    State = "DELETE_TCB"
	StateBound        State = "BOUND"
)

type ConnInfo struct {
	Key         string   `json:"key"`
	Protocol    Protocol `json:"protocol"`
	LocalAddr   string   `json:"localAddr"`
	LocalPort   uint16   `json:"localPort"`
	RemoteAddr  string   `json:"remoteAddr"`
	RemotePort  uint16   `json:"remotePort"`
	State       State    `json:"state"`
	PID         uint32   `json:"pid"`
	ProcessName string   `json:"processName"`
	ProcessPath string   `json:"processPath"`
	Geo         string   `json:"geo"`
}

type Snapshot struct {
	Conns []ConnInfo `json:"conns"`
	Took  int64      `json:"took"`
}

func Key(p Protocol, local string, lport uint16, remote string, rport uint16, pid uint32) string {
	return string(p) + "|" + local + ":" + itoa(lport) + "|" + remote + ":" + itoa(rport) + "|" + itoa32(pid)
}

func itoa(n uint16) string {
	if n == 0 {
		return "0"
	}
	var buf [5]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}

func itoa32(n uint32) string {
	if n == 0 {
		return "0"
	}
	var buf [10]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
