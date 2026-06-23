//go:build windows

package netstat

import "context"

type windowsProvider struct{}

func NewWindowsProvider() Provider { return &windowsProvider{} }

func (w *windowsProvider) Snapshot(ctx context.Context) ([]ConnInfo, error) {
	return gopsutilSnapshot(ctx)
}
