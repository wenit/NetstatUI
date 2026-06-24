# Changelog

## [0.4.0] - 2026-06-24

### Changed
- xdb files no longer embedded in binary; loaded from `<exe-dir>/data/` at startup (side-load). Binary drops from 54 MB to 9.6 MB.
- Remove embed/extract infrastructure; CI packages binary + `data/` together in release archives.

### Added
- Click outside DetailPanel to close it (transparent backdrop covers full window when panel is open).

## [0.3.1] - 2026-06-24

### Performance
- Restore `CreateToolhelp32Snapshot` for Windows process enumeration (replace gopsutil) — `Refresh()` from 3.9s → 10ms on busy machines (400+ processes).

### Changed
- Geo init made fully async: `geo.New()` is instant, searcher pool built in background goroutine via `InitAsync`. Geo column shows `…` while loading, batch-filled on ready.
- Reduce geo searcher pool from 20 to 4.

## [0.3.0] - 2026-06-24

### Added
- Remote IP geolocation column (`Country-City`, e.g. `中国-杭州`) via offline [ip2region](https://github.com/lionsoul2014/ip2region) v3. New `services/geo` package, xdb files embedded in binary, committed to repo, refresh via `scripts/download-xdb.ps1`.
- Settings → General → "IP 地理位置" toggle (`np.geo`, default on). Geo also shown in detail panel and included in search.

### Performance
- LRU 4096 + non-public IP skip in resolver. `monitor.diff()` unchanged (Geo never triggers spurious updates).

## [0.2.1]

Patch release.

## [0.2.0]

- macOS and Linux backends via `gopsutil/v3`; Windows backend migrated to the same library. Cross-platform CI matrix.

## [0.1.0]

- Initial release. Windows netstat GUI: TCP/UDP listing, PID/process enrichment, kill action, Mica window, light/dark/auto theme, i18n (zh-CN/en-US), virtual-scrolled table, filters, context menu, process detail panel.
