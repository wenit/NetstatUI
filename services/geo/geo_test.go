package geo

import (
	"io/fs"
	"net"
	"os"
	"testing"
)

func TestFormatRegion(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		// v4 layout: country|state|city|0|countryCode
		{"US|California|Mountain View|0|US", "US-Mountain View"},
		{"AU|Queensland|Brisbane|0|AU", "AU-Brisbane"},
		{"US|California|0|Google LLC|US", "US"},
		// v6 layout: country|province|city|isp|countryCode
		{"US|California|San Jose|xTom|US", "US-San Jose"},
		// edge cases
		{"US|0|0|0|US", "US"},
		{"0|0|0|0|US", ""},
		{"", ""},
		{"|0|0|0|0", ""},
		{"0|state|city|0|CC", "city"},
		// short
		{"US", "US"},
		{"US|California", "US"},
		{"US|California|San Jose", "US-San Jose"},
		// trimming
		{"US|0|San Jose|0|US", "US-San Jose"},
	}
	for _, c := range cases {
		got := formatRegion(c.in)
		if got != c.want {
			t.Errorf("formatRegion(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestIsNonPublic(t *testing.T) {
	cases := []struct {
		addr string
		want bool
	}{
		{"127.0.0.1", true},
		{"10.0.0.1", true},
		{"192.168.1.1", true},
		{"172.16.0.1", true},
		{"172.31.255.255", true},
		{"172.32.0.1", false},
		{"169.254.1.1", true},
		{"100.64.0.1", true},
		{"100.127.255.255", true},
		{"100.128.0.1", false},
		{"224.0.0.1", true},
		{"255.255.255.255", true},
		{"8.8.8.8", false},
		{"1.1.1.1", false},
		{"::1", true},
		{"fe80::1", true},
		{"fc00::1", true},
		{"fd00::1", true},
		{"ff02::1", true},
		{"2001:4860:4860::8888", false},
	}
	for _, c := range cases {
		ip := net.ParseIP(c.addr)
		if ip == nil {
			t.Errorf("parse %s failed", c.addr)
			continue
		}
		if got := isNonPublic(ip); got != c.want {
			t.Errorf("isNonPublic(%s) = %v, want %v", c.addr, got, c.want)
		}
	}
}

func TestLookupSkipsNonPublic(t *testing.T) {
	var r *Resolver
	if got := r.Lookup("8.8.8.8"); got != "" {
		t.Errorf("nil resolver Lookup = %q, want empty", got)
	}
	empty := &Resolver{}
	if got := empty.Lookup("8.8.8.8"); got != "" {
		t.Errorf("empty resolver Lookup = %q, want empty", got)
	}
	if got := empty.Lookup(""); got != "" {
		t.Errorf("empty addr = %q", got)
	}
	if got := empty.Lookup("*"); got != "" {
		t.Errorf("wildcard = %q", got)
	}
	if got := empty.Lookup("0.0.0.0"); got != "" {
		t.Errorf("0.0.0.0 = %q", got)
	}
	if got := empty.Lookup("::"); got != "" {
		t.Errorf(":: = %q", got)
	}
	if got := empty.Lookup("not-an-ip"); got != "" {
		t.Errorf("invalid = %q", got)
	}
	if got := empty.Lookup("127.0.0.1"); got != "" {
		t.Errorf("loopback = %q", got)
	}
	if got := empty.Lookup("192.168.1.1"); got != "" {
		t.Errorf("private = %q", got)
	}
}

func TestLRUCache(t *testing.T) {
	c := newLRU(3)
	c.put("a", "1")
	c.put("b", "2")
	c.put("c", "3")
	if v, ok := c.get("a"); !ok || v != "1" {
		t.Errorf("get a = (%q, %v), want (1, true)", v, ok)
	}
	c.put("d", "4")
	if _, ok := c.get("b"); ok {
		t.Errorf("b should be evicted")
	}
	if v, ok := c.get("c"); !ok || v != "3" {
		t.Errorf("get c = (%q, %v), want (3, true)", v, ok)
	}
	if v, ok := c.get("d"); !ok || v != "4" {
		t.Errorf("get d = (%q, %v), want (4, true)", v, ok)
	}
	if c.len() != 3 {
		t.Errorf("len = %d, want 3", c.len())
	}
	c.put("a", "1-updated")
	if v, _ := c.get("a"); v != "1-updated" {
		t.Errorf("updated a = %q, want 1-updated", v)
	}
}

func TestCleanField(t *testing.T) {
	cases := map[string]string{
		"":        "",
		"0":       "",
		"  ":      "",
		" 0 ":     "",
		" China ": "China",
		"Tokyo":   "Tokyo",
	}
	for in, want := range cases {
		if got := cleanField(in); got != want {
			t.Errorf("cleanField(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestLookupPublicIP(t *testing.T) {
	disk := os.DirFS("../../data")
	if _, err := fs.Stat(disk, "ip2region_v4.xdb"); err != nil {
		t.Skip("no ../../data/ip2region_v4.xdb; skipping live lookup test")
	}
	r, err := New(disk)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()
	for _, ip := range []string{"8.8.8.8", "1.1.1.1", "223.5.5.5"} {
		got := r.Lookup(ip)
		if got == "" {
			t.Errorf("%s lookup returned empty", ip)
			continue
		}
		t.Logf("%s -> %s", ip, got)
	}
}
