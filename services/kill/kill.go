package kill

type Result struct {
	PID    uint32 `json:"pid"`
	OK     bool   `json:"ok"`
	Reason string `json:"reason"`
}
