package process

type Info struct {
	PID       uint32 `json:"pid"`
	PPID      uint32 `json:"ppid"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	CmdLine   string `json:"cmdLine"`
	User      string `json:"user"`
	StartTime int64  `json:"startTime"`
}
