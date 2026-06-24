package geo

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/lionsoul2014/ip2region/binding/golang/service"
)

const (
	cacheCapacity = 4096
	searcherCount = 4
)

var (
	v4FileName = "ip2region_v4.xdb"
	v6FileName = "ip2region_v6.xdb"
)

type Resolver struct {
	dataDir string
	cache   *lruCache

	mu          sync.Mutex
	ip2r        *service.Ip2Region
	ready       bool
	initErr     error
	initStarted bool
	cb          func()
}

// New returns a Resolver configured to load xdb files from dataDir.
// It performs no I/O; call InitAsync to verify the files and build the
// searcher pool in a background goroutine. Lookup before init completes
// returns "".
func New(dataDir string) (*Resolver, error) {
	if dataDir == "" {
		return nil, fmt.Errorf("geo: empty data directory")
	}
	return &Resolver{dataDir: dataDir, cache: newLRU(cacheCapacity)}, nil
}

// InitAsync verifies the xdb files in dataDir and builds the searcher
// pool in a background goroutine. onReady is invoked on the goroutine
// when init succeeds. The call is idempotent: only the FIRST InitAsync
// call registers a callback and triggers init; later calls return
// immediately without firing onReady again. If init fails, the error is
// captured and surfaced via InitError; onReady is not invoked.
func (r *Resolver) InitAsync(onReady func()) {
	if r == nil {
		return
	}
	r.mu.Lock()
	if r.initStarted {
		r.mu.Unlock()
		return
	}
	r.initStarted = true
	r.cb = onReady
	r.mu.Unlock()
	go r.runInit()
}

func (r *Resolver) runInit() {
	r.mu.Lock()
	err := r.initLocked()
	if err != nil {
		r.initErr = err
		r.mu.Unlock()
		log.Printf("geo: async init failed: %v", err)
		return
	}
	r.ready = true
	cb := r.cb
	r.mu.Unlock()
	if cb != nil {
		cb()
	}
}

func (r *Resolver) initLocked() error {
	v4Path, v6Path, err := resolveDataFiles(r.dataDir)
	if err != nil {
		return err
	}
	v4Cfg, err := service.NewV4Config(service.VIndexCache, v4Path, searcherCount)
	if err != nil {
		return fmt.Errorf("v4 config: %w", err)
	}
	v6Cfg, err := service.NewV6Config(service.VIndexCache, v6Path, searcherCount)
	if err != nil {
		return fmt.Errorf("v6 config: %w", err)
	}
	ip2r, err := service.NewIp2Region(v4Cfg, v6Cfg)
	if err != nil {
		return fmt.Errorf("ip2region init: %w", err)
	}
	r.ip2r = ip2r
	return nil
}

// IsReady reports whether the searcher pool has finished initializing.
func (r *Resolver) IsReady() bool {
	if r == nil {
		return false
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.ready
}

// InitError returns the error from a failed init, or nil if init
// succeeded or has not yet completed.
func (r *Resolver) InitError() error {
	if r == nil {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.initErr
}

func (r *Resolver) Close() {
	if r == nil {
		return
	}
	r.mu.Lock()
	ip2r := r.ip2r
	r.mu.Unlock()
	if ip2r == nil {
		return
	}
	ip2r.CloseTimeout(2 * time.Second)
}

func (r *Resolver) Lookup(addr string) string {
	if r == nil {
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

	r.mu.Lock()
	ip2r := r.ip2r
	r.mu.Unlock()
	if ip2r == nil {
		return ""
	}

	if v, ok := r.cache.get(key); ok {
		return v
	}

	region, err := ip2r.Search(key)
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

// resolveDataFiles checks that the required xdb files exist in dataDir.
// It returns their absolute paths. The caller (InitAsync) runs this in a
// background goroutine so missing files do not block application startup.
func resolveDataFiles(dataDir string) (string, string, error) {
	v4 := filepath.Join(dataDir, v4FileName)
	v6 := filepath.Join(dataDir, v6FileName)
	if !fileExists(v4) {
		return "", "", fmt.Errorf("missing %s", v4)
	}
	if !fileExists(v6) {
		return "", "", fmt.Errorf("missing %s", v6)
	}
	return v4, v6, nil
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
