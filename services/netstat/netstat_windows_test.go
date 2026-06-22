package netstat

import (
	"context"
	"testing"
)

func TestWindowsSnapshot(t *testing.T) {
	SetProvider(NewWindowsProvider())
	conns, err := Get()
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if len(conns) == 0 {
		t.Fatal("no connections returned")
	}
	p := &windowsProvider{}
	_ = p
	t.Logf("got %d connections", len(conns))
	listen, estab := 0, 0
	for _, c := range conns {
		t.Logf("%s %s:%d -> %s:%d state=%s pid=%d proc=%s",
			c.Protocol, c.LocalAddr, c.LocalPort, c.RemoteAddr, c.RemotePort, c.State, c.PID, c.ProcessName)
		if c.State == "LISTEN" {
			listen++
		}
		if c.State == "ESTABLISHED" {
			estab++
		}
	}
	t.Logf("listen=%d established=%d", listen, estab)
	_ = context.Background()
}
