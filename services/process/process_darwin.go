//go:build darwin

package process

import (
	"github.com/shirou/gopsutil/v3/process"
)

type darwinProvider struct{}

func platformProvider() provider { return &darwinProvider{} }

func (d *darwinProvider) SnapshotAll() (map[uint32]Info, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}
	out := make(map[uint32]Info, len(procs))
	for _, p := range procs {
		pid := uint32(p.Pid)
		info := Info{PID: pid}
		if name, err := p.Name(); err == nil {
			info.Name = name
		}
		if exe, err := p.Exe(); err == nil {
			info.Path = exe
		}
		if ppid, err := p.Ppid(); err == nil {
			info.PPID = uint32(ppid)
		}
		out[pid] = info
	}
	return out, nil
}

func (d *darwinProvider) QueryPath(pid uint32) (string, error) {
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return "", err
	}
	return p.Exe()
}
