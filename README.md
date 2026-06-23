# NetstatUI

> Win11 Fluent-style desktop network port & connection inspector (like `netstat`).

[English](./README.md) · [简体中文](./README.zh-CN.md)

A read-only Windows desktop tool that lists every TCP/UDP connection with its owning PID, process name, and executable path — and lets you kill the process with a single click. Built on **Wails 3** + **Vue 3** with a Fluent **Mica**-backed frameless window.

---

## Highlights

- 📡 **Full visibility** — every TCP4 / TCP6 / UDP4 / UDP6 socket, with local + remote endpoints, state, PID and resolved process path.
- 🎨 **Win11 Fluent UI** — Mica backdrop, light / dark / auto theme, custom title bar, compact & comfortable density.
- ⚡ **Live updates** — diff-based streaming; first frame is `conn:full`, subsequent ticks push only `added` / `removed` / `updated`.
- 🚀 **Virtual scroll** — handles 10,000+ rows at 60 fps via custom absolute-positioned virtual scrolling.
- 🔍 **Rich filtering** — search across all fields, protocol chips, state chips, listen-only / external-only toggles.
- 🪓 **One-click kill** — confirm dialog → `TerminateProcess`; auto-refresh immediately afterwards.
- 🌍 **Bilingual** — English (default) / Simplified Chinese, with OS-locale auto-detection.
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
.\build.ps1                  # Windows (recommended, bypasses file-lock issue)
# or
wails3 build                 # macOS / Linux
```

Output binary: `bin/NetstatUI.exe` (Windows) / `bin/NetstatUI` (Unix).

> **Tip:** if `wails3 build` on Windows fails with `Access is denied`, use `build.ps1` — it forces in-place regeneration of TS bindings and avoids the `RemoveAll+Rename` step that Windows SearchIndexer / Defender holds open.

---

## Usage

1. Launch `NetstatUI.exe`.
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

## Known Limitations

- **Some loopback listeners may be missing** — Windows 11 22H2+ silently drops a subset of `127.0.0.1` LISTEN entries from `GetTcpTable2` / `GetExtendedTcpTable`. `netstat -ano` shows them because it uses WMI. We use the same iphlpapi path as `go-netstat`, so the same limitation applies.
- **Windows only** — Linux/macOS providers are stubbed (`services/netstat/provider.go` interface) but not implemented. The UI and `kill` service are also Windows-specific.

See [`AGENTS.md` → Known pitfalls](./AGENTS.md#known-坑) for more.

---

## Project Layout

```
.
├── main.go                       # Wails app entry, registers services + events
├── app.go                        # AppService: KillProcess / GetProcessDetail / OpenProcessFolder / GetSystemLocale
├── services/
│   ├── netstat/                  # TCP/UDP snapshot via go-netstat
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
- Step-by-step extension guides (add a column, a filter, a locale, a backend method)

To regenerate TypeScript bindings after changing Go service signatures:

```bash
wails3 generate bindings -ts -clean=false
```

---

## License

MIT — see [`LICENSE`](./LICENSE).