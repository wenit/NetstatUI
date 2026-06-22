package monitor

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/zwb/network-ports/services/netstat"
	"github.com/zwb/network-ports/services/process"
)

type Stats struct {
	Total        int `json:"total"`
	Listen       int `json:"listen"`
	Established  int `json:"established"`
	UDP          int `json:"udp"`
	FilteredHits int `json:"filteredHits"`
}

type Diff struct {
	Added   []netstat.ConnInfo `json:"added"`
	Removed []string           `json:"removed"`
	Updated []netstat.ConnInfo `json:"updated"`
}

type Monitor struct {
	cache     *process.Cache
	interval  time.Duration
	mu        sync.Mutex
	last      map[string]netstat.ConnInfo
	ticker    *time.Ticker
	stopCh    chan struct{}
	running   bool
	lastStats Stats
}

func New(cache *process.Cache) *Monitor {
	return &Monitor{
		cache:    cache,
		interval: 2 * time.Second,
		last:     make(map[string]netstat.ConnInfo),
		stopCh:   make(chan struct{}),
	}
}

func (m *Monitor) app() *application.App {
	return application.Get()
}

func (m *Monitor) ServiceName() string { return "monitor" }

func (m *Monitor) ServiceStartup(_ context.Context, _ application.ServiceOptions) error {
	if err := m.cache.Refresh(); err != nil {
		log.Printf("monitor: initial process refresh: %v", err)
	}
	conns, stats, err := m.collect()
	if err != nil {
		return err
	}
	m.mu.Lock()
	m.last = indexBy(conns)
	m.lastStats = stats
	m.mu.Unlock()
	m.app().Event.Emit("conn:full", conns)
	m.app().Event.Emit("conn:stats", stats)
	return nil
}

func (m *Monitor) ServiceShutdown() error {
	m.Stop()
	return nil
}

type SnapshotResult struct {
	Conns []netstat.ConnInfo `json:"conns"`
	Stats Stats              `json:"stats"`
}

func (m *Monitor) GetSnapshot() (SnapshotResult, error) {
	if err := m.cache.Refresh(); err != nil {
		log.Printf("monitor: get snapshot process refresh: %v", err)
	}
	conns, stats, err := m.collect()
	if err != nil {
		return SnapshotResult{}, err
	}
	m.mu.Lock()
	m.last = indexBy(conns)
	m.lastStats = stats
	m.mu.Unlock()
	return SnapshotResult{Conns: conns, Stats: stats}, nil
}

func (m *Monitor) Start(intervalMs int) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.running {
		return true
	}
	m.interval = time.Duration(intervalMs) * time.Millisecond
	if m.interval < 200*time.Millisecond {
		m.interval = 200 * time.Millisecond
	}
	m.ticker = time.NewTicker(m.interval)
	m.stopCh = make(chan struct{})
	m.running = true
	go m.loop()
	return true
}

func (m *Monitor) Stop() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.running {
		return true
	}
	m.running = false
	m.ticker.Stop()
	close(m.stopCh)
	return true
}

func (m *Monitor) SetInterval(intervalMs int) bool {
	m.mu.Lock()
	m.interval = time.Duration(intervalMs) * time.Millisecond
	if m.interval < 200*time.Millisecond {
		m.interval = 200 * time.Millisecond
	}
	running := m.running
	m.mu.Unlock()
	if running {
		m.Stop()
		m.Start(intervalMs)
	}
	return true
}

func (m *Monitor) IsRunning() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.running
}

func (m *Monitor) GetInterval() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return int(m.interval / time.Millisecond)
}

func (m *Monitor) RefreshNow() bool {
	go func() {
		if err := m.cache.Refresh(); err != nil {
			log.Printf("monitor: refresh process: %v", err)
		}
		conns, stats, err := m.collect()
		if err != nil {
			log.Printf("monitor: refresh now: %v", err)
			return
		}
		m.mu.Lock()
		m.last = indexBy(conns)
		m.lastStats = stats
		m.mu.Unlock()
		m.app().Event.Emit("conn:full", conns)
		m.app().Event.Emit("conn:stats", stats)
	}()
	return true
}

func (m *Monitor) loop() {
	tick := m.ticker
	stop := m.stopCh
	for {
		select {
		case <-stop:
			return
		case <-tick.C:
			m.tick()
		}
	}
}

func (m *Monitor) tick() {
	if err := m.cache.Refresh(); err != nil {
		log.Printf("monitor: process refresh: %v", err)
	}
	conns, stats, err := m.collect()
	if err != nil {
		m.app().Event.Emit("conn:error", err.Error())
		return
	}
	diff := m.diff(conns)
	m.mu.Lock()
	m.last = indexBy(conns)
	m.lastStats = stats
	m.mu.Unlock()
	m.app().Event.Emit("conn:diff", diff)
	m.app().Event.Emit("conn:stats", stats)
}

func (m *Monitor) collect() ([]netstat.ConnInfo, Stats, error) {
	start := time.Now()
	conns, err := netstat.Get()
	if err != nil {
		return nil, Stats{}, err
	}
	pids := uniquePIDs(conns)
	procs := m.cache.Enrich(pids)
	for i := range conns {
		if info, ok := procs[conns[i].PID]; ok {
			conns[i].ProcessName = info.Name
			conns[i].ProcessPath = info.Path
		}
	}
	stats := computeStats(conns)
	_ = start
	return conns, stats, nil
}

func (m *Monitor) diff(cur []netstat.ConnInfo) Diff {
	m.mu.Lock()
	prev := m.last
	m.mu.Unlock()
	curMap := make(map[string]netstat.ConnInfo, len(cur))
	for _, c := range cur {
		curMap[c.Key] = c
	}
	var d Diff
	for k, c := range curMap {
		old, ok := prev[k]
		if !ok {
			d.Added = append(d.Added, c)
		} else if old.State != c.State || old.ProcessName != c.ProcessName {
			d.Updated = append(d.Updated, c)
		}
	}
	for k := range prev {
		if _, ok := curMap[k]; !ok {
			d.Removed = append(d.Removed, k)
		}
	}
	return d
}

func indexBy(conns []netstat.ConnInfo) map[string]netstat.ConnInfo {
	m := make(map[string]netstat.ConnInfo, len(conns))
	for _, c := range conns {
		m[c.Key] = c
	}
	return m
}

func uniquePIDs(conns []netstat.ConnInfo) []uint32 {
	seen := make(map[uint32]struct{}, len(conns))
	out := make([]uint32, 0, len(conns))
	for _, c := range conns {
		if _, ok := seen[c.PID]; !ok {
			seen[c.PID] = struct{}{}
			out = append(out, c.PID)
		}
	}
	return out
}

func computeStats(conns []netstat.ConnInfo) Stats {
	var s Stats
	s.Total = len(conns)
	for _, c := range conns {
		switch {
		case c.Protocol == netstat.ProtocolUDP4 || c.Protocol == netstat.ProtocolUDP6:
			s.UDP++
		case c.State == netstat.StateListen:
			s.Listen++
		case c.State == netstat.StateEstablished:
			s.Established++
		}
	}
	return s
}
