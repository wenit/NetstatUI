//go:build darwin

package netstat

import "context"

type darwinProvider struct{}

func NewDarwinProvider() Provider { return &darwinProvider{} }

func (d *darwinProvider) Snapshot(ctx context.Context) ([]ConnInfo, error) {
	return gopsutilSnapshot(ctx)
}
