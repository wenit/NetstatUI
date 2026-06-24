# Changelog

## [0.3.0] - 2026-06-24

### Added
- Remote IP geolocation column (`Country-City`, e.g. `中国-杭州`) via offline [ip2region](https://github.com/lionsoul2014/ip2region) v3. New `services/geo` package, xdb files embedded in binary, `data/*.xdb` gitignored, refresh via `scripts/download-xdb.ps1`.
- Settings → General → "IP 地理位置" toggle (`np.geo`, default on). Geo also shown in detail panel and included in search.

### Performance
- LRU 4096 + non-public IP skip in resolver. `monitor.diff()` unchanged (Geo never triggers spurious updates).

## [0.2.1]

Patch release.

## [0.2.0]

- macOS and Linux backends via `gopsutil/v3`; Windows backend migrated to the same library. Cross-platform CI matrix.

## [0.1.0]

- Initial release. Windows netstat GUI: TCP/UDP listing, PID/process enrichment, kill action, Mica window, light/dark/auto theme, i18n (zh-CN/en-US), virtual-scrolled table, filters, context menu, process detail panel.
