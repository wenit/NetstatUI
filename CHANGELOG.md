# Changelog

All notable changes to NetstatUI are documented here. Versions follow [Semantic Versioning](https://semver.org/).

## [0.3.0] - 2026-06-24

### Added
- **IP geolocation column** for remote addresses, displayed right of the remote address column in the connection table. Format: `Country-City` (e.g. `United States-Mountain View`, `中国-杭州`).
- New `services/geo` package: offline IP-to-region lookup via [ip2region](https://github.com/lionsoul2014/ip2region) v3 (`VIndexCache` strategy, ~50µs/query, ~1MB index in RAM, searcher pool, built-in LRU 4096).
- Bundled `ip2region_v4.xdb` + `ip2region_v6.xdb` (embedded via `//go:embed all:data`; first launch extracts to user cache, `<exe-dir>/data/` side-load takes priority for power users).
- Settings → General → "IP 地理位置" toggle (localStorage `np.geo`, default **on**). Hides the UI column only; backend resolution always runs.
- Geolocation included in the connection detail panel and in the filter search haystack.
- `scripts/download-xdb.ps1` to refresh xdb files from the official ip2region repo.

### Performance
- Backend Geo lookup is cached in a bounded LRU; non-public IPs (loopback, RFC1918, link-local, ULA, CGNAT, multicast, reserved) are skipped before hitting the resolver.
- `monitor.diff()` is unchanged: still compares only `State` + `ProcessName`, so Geo updates never trigger spurious `conn:diff` events.

## [0.2.1] - 2026

### Changed
- Patch release.

## [0.2.0] - 2026

### Added
- macOS backend via `gopsutil/v3`.
- Linux backend via `gopsutil/v3`.
- Cross-platform CI matrix (Windows / macOS / Linux).

### Changed
- Windows backend migrated from `go-netstat` + `Toolhelp32` to `gopsutil/v3` for unification across platforms.

## [0.1.0] - 2026

### Added
- Initial release: Windows netstat graphical UI.
- TCP4/TCP6/UDP4/UDP6 connection listing with PID/process enrichment.
- Kill process action.
- Frameless Mica window (Win11 22621+), light/dark/auto theme, i18n (zh-CN / en-US).
- Virtual-scrolled table, address/port/process filters, search, context menu, process detail panel.
