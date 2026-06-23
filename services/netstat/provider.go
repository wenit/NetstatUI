package netstat

import "context"

type Provider interface {
	// Snapshot returns all current connections across TCP4/TCP6/UDP4/UDP6.
	// PID resolution happens here; process name/path may be filled by caller via ProcessCache.
	Snapshot(ctx context.Context) ([]ConnInfo, error)
}

var defaultProvider Provider

func SetProvider(p Provider) { defaultProvider = p }
func Get() ([]ConnInfo, error) {
	if defaultProvider == nil {
		return nil, ErrNoProvider
	}
	return defaultProvider.Snapshot(context.Background())
}

// NewPlatformProvider returns the Provider implementation for the current
// operating system. Each platform file defines its own platformProvider()
// helper; this wrapper picks the active one.
func NewPlatformProvider() Provider {
	return platformProvider()
}
