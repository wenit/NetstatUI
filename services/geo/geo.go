package geo

import (
	"fmt"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lionsoul2014/ip2region/binding/golang/service"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

const (
	cacheCapacity = 4096
)

var (
	v4FileName = "ip2region_v4.xdb"
	v6FileName = "ip2region_v6.xdb"
)

type Resolver struct {
	ip2r  *service.Ip2Region
	cache *lruCache
}

func New(src fs.FS) (*Resolver, error) {
	v4Path, v6Path, err := materialize(src)
	if err != nil {
		return nil, err
	}

	ip2r, err := service.NewIp2RegionWithPath(v4Path, v6Path)
	if err != nil {
		return nil, fmt.Errorf("ip2region init: %w", err)
	}

	return &Resolver{
		ip2r:  ip2r,
		cache: newLRU(cacheCapacity),
	}, nil
}

func (r *Resolver) Close() {
	if r == nil || r.ip2r == nil {
		return
	}
	r.ip2r.CloseTimeout(2 * time.Second)
}

func (r *Resolver) Lookup(addr string) string {
	if r == nil || r.ip2r == nil {
		return ""
	}
	addr = strings.TrimSpace(addr)
	if addr == "" || addr == "*" || addr == "0.0.0.0" || addr == "::" {
		return ""
	}

	ip := net.ParseIP(addr)
	if ip == nil {
		return ""
	}
	if isNonPublic(ip) {
		return ""
	}
	key := ip.String()

	if v, ok := r.cache.get(key); ok {
		return v
	}

	region, err := r.ip2r.Search(key)
	if err != nil || region == "" {
		r.cache.put(key, "")
		return ""
	}
	short := formatRegion(region)
	r.cache.put(key, short)
	return short
}

func (r *Resolver) CacheSize() int {
	if r == nil {
		return 0
	}
	return r.cache.len()
}

// formatRegion converts "中国|江苏省|南京市|0|CN" → "中国-南京".
// Empty/zero fields are dropped.
//
// ip2region v4 layout: country|state|city|0|countryCode
// ip2region v6 layout: country|province|city|isp|countryCode
// (slot 3 is "0" when the data only resolves to a higher level; ISP lives there
// in v6 records that include one).
func formatRegion(region string) string {
	parts := strings.Split(region, "|")
	if len(parts) < 3 {
		return cleanField(parts[0])
	}
	country := cleanField(parts[0])
	city := cleanField(parts[2])
	if country == "" && city == "" {
		return ""
	}
	if country == "" {
		return city
	}
	if city == "" {
		return country
	}
	return country + "-" + city
}

func cleanField(s string) string {
	s = strings.TrimSpace(s)
	if s == "" || s == "0" {
		return ""
	}
	return s
}

// materialize places the xdb files in a stable on-disk location and
// returns their absolute paths. Resolution order:
//   1. <exe-dir>/data/                (side-load: power users can swap xdb files)
//   2. embedded FS  →  user cache dir (one-time extraction; ~36MB total)
func materialize(src fs.FS) (string, string, error) {
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		if v4, v6, ok := sideLoad(exeDir); ok {
			return v4, v6, nil
		}
	}

	cacheDir, err := cacheRoot()
	if err != nil {
		return "", "", err
	}
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return "", "", err
	}

	v4Dest := filepath.Join(cacheDir, v4FileName)
	v6Dest := filepath.Join(cacheDir, v6FileName)
	if err := extractIfMissing(src, v4FileName, v4Dest); err != nil {
		return "", "", err
	}
	if err := extractIfMissing(src, v6FileName, v6Dest); err != nil {
		return "", "", err
	}
	if err := xdb.VerifyFromFile(v4Dest); err != nil {
		return "", "", fmt.Errorf("verify %s: %w", v4FileName, err)
	}
	if err := xdb.VerifyFromFile(v6Dest); err != nil {
		return "", "", fmt.Errorf("verify %s: %w", v6FileName, err)
	}
	return v4Dest, v6Dest, nil
}

func sideLoad(exeDir string) (string, string, bool) {
	v4 := filepath.Join(exeDir, "data", v4FileName)
	v6 := filepath.Join(exeDir, "data", v6FileName)
	if fileExists(v4) && fileExists(v6) {
		return v4, v6, true
	}
	return "", "", false
}

func cacheRoot() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "NetstatUI", "data"), nil
}

func extractIfMissing(src fs.FS, name, dest string) error {
	if fileExists(dest) {
		return nil
	}
	data, err := fs.ReadFile(src, name)
	if err != nil {
		return fmt.Errorf("read embedded %s: %w", name, err)
	}
	tmp := dest + ".tmp-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", tmp, err)
	}
	if err := os.Rename(tmp, dest); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("rename %s: %w", dest, err)
	}
	return nil
}

func fileExists(p string) bool {
	st, err := os.Stat(p)
	return err == nil && !st.IsDir()
}

func isNonPublic(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsUnspecified() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsMulticast() {
		return true
	}
	if ip.IsPrivate() {
		return true
	}
	if v4 := ip.To4(); v4 != nil {
		if v4[0] == 100 && v4[1] >= 64 && v4[1] <= 127 {
			return true
		}
		if v4[0] >= 240 {
			return true
		}
	}
	return false
}

// lruCache is a tiny generic LRU used to memoize ip2region lookups.
type lruCache struct {
	mu       sync.Mutex
	capacity int
	items    map[string]*lruNode
	head     *lruNode
	tail     *lruNode
}

type lruNode struct {
	key, val string
	prev     *lruNode
	next     *lruNode
}

func newLRU(capacity int) *lruCache {
	if capacity <= 0 {
		capacity = 1024
	}
	return &lruCache{
		capacity: capacity,
		items:    make(map[string]*lruNode, capacity),
	}
}

func (c *lruCache) get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	n, ok := c.items[key]
	if !ok {
		return "", false
	}
	c.touch(n)
	return n.val, true
}

func (c *lruCache) put(key, val string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if n, ok := c.items[key]; ok {
		n.val = val
		c.touch(n)
		return
	}
	n := &lruNode{key: key, val: val}
	c.items[key] = n
	c.pushFront(n)
	if len(c.items) > c.capacity {
		old := c.tail
		if old != nil {
			c.remove(old)
			delete(c.items, old.key)
		}
	}
}

func (c *lruCache) len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.items)
}

func (c *lruCache) touch(n *lruNode) {
	if c.head == n {
		return
	}
	c.remove(n)
	c.pushFront(n)
}

func (c *lruCache) pushFront(n *lruNode) {
	n.prev = nil
	n.next = c.head
	if c.head != nil {
		c.head.prev = n
	}
	c.head = n
	if c.tail == nil {
		c.tail = n
	}
}

func (c *lruCache) remove(n *lruNode) {
	if n.prev != nil {
		n.prev.next = n.next
	} else {
		c.head = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	} else {
		c.tail = n.prev
	}
	n.prev, n.next = nil, nil
}
