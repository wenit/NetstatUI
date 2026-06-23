//go:build linux

package netstat

import "context"

type linuxProvider struct{}

func NewLinuxProvider() Provider { return &linuxProvider{} }

func (l *linuxProvider) Snapshot(ctx context.Context) ([]ConnInfo, error) {
	return gopsutilSnapshot(ctx)
}
