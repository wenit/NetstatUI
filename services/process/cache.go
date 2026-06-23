package process

import (
	"errors"
	"sync"
	"time"
)

var ErrNotFound = errors.New("process: not found")

type Cache struct {
	mu       sync.RWMutex
	entries  map[uint32]*entry
	pathTTL  time.Duration
	provider provider
}

type entry struct {
	info     Info
	pathAt   time.Time
	hasPath  bool
}

type provider interface {
	SnapshotAll() (map[uint32]Info, error)
	QueryPath(pid uint32) (string, error)
}

func NewCache() *Cache {
	return &Cache{
		entries:  make(map[uint32]*entry),
		pathTTL:  30 * time.Second,
		provider: platformProvider(),
	}
}

func (c *Cache) Refresh() error {
	all, err := c.provider.SnapshotAll()
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for pid, info := range all {
		e, ok := c.entries[pid]
		if !ok {
			c.entries[pid] = &entry{info: info}
			continue
		}
		e.info.Name = info.Name
		e.info.PPID = info.PPID
		if e.info.PID == 0 {
			e.info.PID = info.PID
		}
	}
	for pid := range c.entries {
		if _, ok := all[pid]; !ok {
			delete(c.entries, pid)
		}
	}
	return nil
}

func (c *Cache) Get(pid uint32) (Info, error) {
	c.mu.RLock()
	e, ok := c.entries[pid]
	c.mu.RUnlock()
	if !ok {
		return Info{}, ErrNotFound
	}
	if !e.hasPath || time.Since(e.pathAt) > c.pathTTL {
		if p, err := c.provider.QueryPath(pid); err == nil {
			c.mu.Lock()
			e.info.Path = p
			e.hasPath = true
			e.pathAt = time.Now()
			c.mu.Unlock()
		}
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return e.info, nil
}

func (c *Cache) Enrich(pids []uint32) map[uint32]Info {
	out := make(map[uint32]Info, len(pids))
	for _, pid := range pids {
		if info, err := c.Get(pid); err == nil {
			out[pid] = info
		}
	}
	return out
}
