# NetstatUI

> A graphical desktop user interface for the system `netstat` command.

[English](./README.md) · [简体中文](./README.zh-CN.md)

NetstatUI wraps the operating system's native `netstat` machinery — sockets, ports, PIDs, processes — into a live, filterable, themeable desktop window. Instead of memorising flags or piping through `grep`, you get a sortable table with one-click kill and an open-in-folder action.

Built on **Wails 3** + **Vue 3** with a cross-platform architecture: Windows is the primary target today (with Win11 Fluent **Mica** styling), and the `Provider` interface makes macOS / Linux straightforward to add.

---

## Highlights

- 📡 **Full visibility** — every TCP4 / TCP6 / UDP4 / UDP6 socket, with local + remote endpoints, state, PID and resolved process path.
- ⚡ **Live incremental updates** — diff-based streaming; first frame is `conn:full`, subsequent ticks push only `added` / `removed` / `updated`.
- 🚀 **Custom virtual scroll** — handles 10,000+ rows at 60 fps via absolute-positioned virtual scrolling (no third-party grid).
- 🔍 **Rich filtering** — full-text search across all columns, protocol chips, state chips, listen-only / external-only toggles.
- 🪓 **One-click kill** — confirm dialog → terminate; auto-refresh immediately afterwards.
- 🎨 **Adaptive theming** — light / dark / auto (follows OS), Mica backdrop on Windows 11 22621+, compact & comfortable density.
- 🌍 **Bilingual UI** — English (default) / Simplified Chinese, with OS-locale auto-detection.
- 💾 **Persistent settings** — theme, locale, interval, running state stored in localStorage (`np.*` keys).

---

## Tech Stack

| Layer       | Choice                                                   |
| ----------- | -------------------------------------------------------- |
| Shell       | Wails 3 `v3.0.0-alpha.98` (Frameless + Mica + WebView2)  |
| Backend     | Go 1.25+ (`go-netstat`, `windows.TerminateProcess`, ...) |
| Frontend    | Vue 3 + TypeScript + Vite 8                              |
| State       | Pinia                                                     |
| i18n        | vue-i18n `@^9` (`legacy: false`)                         |
| Utilities   | @vueuse/core `^14`                                       |

See [`AGENTS.md`](./AGENTS.md) for the full architecture, data-flow diagram and invariants.

---

## Build from Source

Requires **Go 1.25+**, **Node.js 20+**, and the **Wails 3** CLI.

```bash
# one-time
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0-alpha.98

# dev mode (hot reload)
wails3 dev

# production build
wails3 build                # works on all platforms
# or, on Windows only:
.\build.ps1                 # bypasses file-lock issue (see below)
```

Output: `bin/NetstatUI.exe` (Windows) / `bin/NetstatUI` (macOS / Linux).

> **Tip:** if `wails3 build` on Windows fails with `Access is denied`, use `build.ps1` — it forces in-place regeneration of TS bindings and skips the `RemoveAll+Rename` step that Windows SearchIndexer / Defender holds open.

---

## Usage

1. Launch `NetstatUI`.
2. The table shows all live connections; the **StatsBar** at the bottom summarises total / listen / established / udp.
3. **FilterBar** — narrow by protocol, state, or type a search query (matches any visible column).
4. **Toolbar** — pick a refresh interval (5 / 15 / 30 / 60 s), pause/resume, or hit the refresh button for an immediate pull.
5. Click a row to open the **DetailPanel** (full process info + open-folder action).
6. Right-click for the context menu — **Kill process** raises a confirmation dialog.

### Settings

Open the **⚙ Settings** dialog from the title bar:

- **General** — theme (auto / light / dark), locale (English / 简体中文), density (compact / comfortable), refresh interval.
- **Advanced** — clear localStorage to reset everything.

---

## Platform Support

| Platform | Status                                                          |
| -------- | --------------------------------------------------------------- |
| Windows  | ✅ Fully implemented — uses `go-netstat` (GetTcpTable2/6 + GetExtendedUdpTable) and Windows `TerminateProcess` |
| macOS    | 🟡 Stubbed — backend returns "not supported"; see `services/netstat/provider.go` for the `Provider` interface to implement |
| Linux    | 🟡 Stubbed — same as macOS; `/proc/net/tcp{,6}` is the natural data source |

The cross-platform architecture is in place — `services/netstat/`, `services/process/`, `services/kill/`, `services/system/` all use `//go:build` tags and an injectable `Provider`. Adding a new platform is `netstat_<os>.go` + a `process_<os>.go` for PID resolution. See [`AGENTS.md` → Extension guide](./AGENTS.md#扩展指引).

---

## Known Limitations

- **Some loopback listeners may be missing on Windows 11 22H2+** — `iphlpapi.dll`'s `GetTcpTable2` / `GetExtendedUdpTable` silently drop a subset of `127.0.0.1` LISTEN entries. `netstat -ano` shows them because it uses WMI; we use the same iphlpapi path as `go-netstat`, so the same limitation applies.
- **Windows-only UI polish** — Mica backdrop, snap layouts, and Win11 Fluent controls are Windows-specific. macOS / Linux will get the generic WebView chrome until native styling is added.

See [`AGENTS.md` → Known pitfalls](./AGENTS.md#已知坑) for more.

---

## Project Layout

```
.
├── main.go                       # Wails app entry, registers services + events
├── app.go                        # AppService: KillProcess / GetProcessDetail / OpenProcessFolder / GetSystemLocale
├── services/
│   ├── netstat/                  # TCP/UDP snapshot (Provider interface + Windows impl)
│   ├── process/                  # PID → name/path cache (Toolhelp32Snapshot + QueryFullProcessImageNameW)
│   ├── monitor/                  # Polling, diff, Wails event emit
│   ├── kill/                     # TerminateProcess wrapper
│   └── system/                   # GetSystemLocale (registry read)
├── frontend/
│   ├── bindings/                 # GENERATED — wails3 generate bindings -ts
│   └── src/
│       ├── App.vue               # Layout + initial fetchSnapshot
│       ├── components/           # TitleBar, Toolbar, FilterBar, ConnectionTable, DetailPanel, ...
│       ├── composables/          # useConnections (event subscription + diff), useFilters
│       ├── locales/              # en-US, zh-CN
│       └── stores/settings.ts    # Pinia store (theme/locale/interval/running/density)
├── build/
│   ├── config.yml                # Wails build metadata (company, product, identifier)
│   └── windows/info.json         # Windows resource metadata
├── build.ps1                     # Safe Windows build (bypasses file-lock)
└── .github/workflows/build.yml   # CI: build windows-amd64 + darwin-{arm64,amd64}
```

---

## Development

See [`AGENTS.md`](./AGENTS.md) for:

- Entry points and key invariants
- Data-flow diagram and event payload schema
- Critical gotchas (byte order, IPv6 struct offsets, first-frame race, embed requirement)
- Step-by-step extension guides (add a column, a filter, a locale, a backend method, a platform)

To regenerate TypeScript bindings after changing Go service signatures:

```bash
wails3 generate bindings -ts -clean=false
```

---

## License

MIT — see [`LICENSE`](./LICENSE).