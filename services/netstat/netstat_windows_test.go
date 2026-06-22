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
	found9245 := false
	for _, c := range conns {
		t.Logf("%s %s:%d -> %s:%d state=%s pid=%d proc=%s",
			c.Protocol, c.LocalAddr, c.LocalPort, c.RemoteAddr, c.RemotePort, c.State, c.PID, c.ProcessName)
		if c.State == "LISTEN" {
			listen++
		}
		if c.State == "ESTABLISHED" {
			estab++
		}
		if c.LocalPort == 9245 && c.LocalAddr == "127.0.0.1" {
			found9245 = true
		}
	}
	t.Logf("listen=%d established=%d found9245=%v", listen, estab, found9245)
	if !found9245 {
		t.Log("WARN: 127.0.0.1:9245 (PID 36020) not found — go-netstat GetTcpTable2 may also miss this entry")
	}
	_ = context.Background()
}
