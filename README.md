# NetstatUI

> A graphical desktop user interface for the system `netstat` command.

[English](./README.md) · [简体中文](./README.zh-CN.md)

NetstatUI wraps the operating system's native `netstat` machinery — sockets, ports, PIDs, processes — into a live, filterable, themeable desktop window. Instead of memorising flags or piping through `grep`, you get a sortable table with one-click kill and an open-in-folder action.

Built on **Wails 3** + **Vue 3**, with **gopsutil/v3** as a unified cross-platform backend. Windows, macOS, and Linux all share the same `services/netstat/`, `services/process/`, and `services/kill/` packages — only the platform-specific implementation file differs.

---

## Screenshots

| Dark | Light |
|:----:|:-----:|
| ![Dark](./docs/main-dark.png) | ![Light](./docs/main-light.png) |

Remote IP geolocation (`Country-City`) is shown right of the remote address column. Toggle in **Settings → General → IP Geolocation** (`np.geo`).

---

## Highlights

- 🌍 **IP geolocation** — remote address column shows `Country-City` (e.g. `中国-杭州`) via offline [ip2region](https://github.com/lionsoul2014/ip2region). Embedded xdb, no network calls.

- 📡 **Full visibility** — every TCP4 / TCP6 / UDP4 / UDP6 socket, with local + remote endpoints, state, PID and resolved process path.
- 🖥️ **Three platforms, one codebase** — Windows 10/11, macOS 12+, and Linux (WebKitGTK). All features work identically across the trio.
- ⚡ **Live incremental updates** — diff-based streaming; first frame is `conn:full`, subsequent ticks push only `added` / `removed` / `updated`.
- 🚀 **Custom virtual scroll** — handles 10,000+ rows at 60 fps via absolute-positioned virtual scrolling (no third-party grid).
- 🔍 **Rich filtering** — full-text search across all columns, protocol chips, state chips, listen-only / external-only toggles.
- 🪓 **One-click kill** — confirm dialog → `TerminateProcess` (Windows) / `SIGKILL` (macOS/Linux); auto-refresh immediately afterwards.
- 🎨 **Adaptive theming** — light / dark / auto (follows OS), Mica backdrop on Windows 11 22621+, compact & comfortable density.
- 🌍 **Bilingual UI** — English (default) / Simplified Chinese, with OS-locale auto-detection.
- 💾 **Persistent settings** — theme, locale, interval, running state, geo toggle stored in localStorage (`np.*` keys).

---

## Tech Stack

| Layer       | Choice                                                                       |
| ----------- | ---------------------------------------------------------------------------- |
| Shell       | Wails 3 `v3.0.0-alpha.98` (Frameless + Mica on Win11; WebKitGTK on Linux)    |
| Backend     | Go 1.25+ — [`gopsutil/v3`](https://github.com/shirou/gopsutil) for net + process + signal |
| Frontend    | Vue 3 + TypeScript + Vite 8                                                  |
| State       | Pinia                                                                         |
| i18n        | vue-i18n `@^9` (`legacy: false`)                                             |
| Utilities   | @vueuse/core `^14`                                                           |

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
wails3 build                  # works on all three platforms
# or, on Windows only:
.\build.ps1                   # bypasses file-lock issue (see below)
```

Output:
- Windows: `bin/NetstatUI.exe`
- macOS / Linux: `bin/NetstatUI`

### Linux build dependencies

`wails3 build` on Linux requires WebKitGTK + GTK 3:

```bash
sudo apt-get install -y libgtk-3-dev libwebkit2gtk-4.1-dev pkg-config
```

For older distributions that ship `libwebkit2gtk-4.0`, replace `4.1` with `4.0` in the package name.

> **Tip:** if `wails3 build` on Windows fails with `Access is denied`, use `build.ps1` — it forces in-place regeneration of TS bindings and skips the `RemoveAll+Rename` step that Windows SearchIndexer / Defender holds open.

---

## Usage

1. Launch `NetstatUI`.
2. The table shows all live connections; the **StatsBar** at the bottom summarises total / listen / established / udp.
3. **FilterBar** — narrow by protocol, state, or type a search query (matches any visible column).
4. **Toolbar** — pick a refresh interval (5 / 15 / 30 / 60 s), pause/resume, or hit the refresh button for an immediate pull.
5. Click a row to open the **DetailPanel** (full process info + open-folder action).
6. Right-click for the context menu — **Kill process** raises a confirmation dialog.

### First run on Windows (SmartScreen)

The published `NetstatUI.exe` is **not code-signed with a paid certificate** (an EV/OV signing cert costs $300–500/year, and the project is currently distributed free of charge). As a result, Windows 10/11 shows a **Microsoft Defender SmartScreen** warning the first time the binary is launched from a fresh machine:

> "Windows protected your PC — Microsoft Defender SmartScreen prevented an unrecognized app from starting."

This is **not malware detection** — it is purely a reputation-based warning for unsigned executables. To run the app:

1. In the SmartScreen dialog, click **More info** ("更多信息").
2. A **Run anyway** ("仍要运行") button appears — click it.

SmartScreen remembers the file's hash for that machine, so subsequent launches of the **same** `NetstatUI.exe` will not prompt again. If you replace the binary with a newer build, you'll see the dialog once more.

If you want to suppress the warning permanently without buying a certificate, see the [Windows build guide](./README.md#build-from-source) and sign the binary yourself with `signtool sign /a` (self-signed certs still trigger SmartScreen, but at least the publisher name shows). For production releases, an [EV code-signing certificate from DigiCert / Sectigo](https://learn.microsoft.com/en-us/windows/security/identity-protection/access-control/access-control) is the only path to a SmartScreen-clean app.

### First run on macOS (Gatekeeper)

The published `NetstatUI` binary (or the `.app` bundle inside `NetstatUI.app`) is **not signed with an Apple Developer ID and not notarised** (the Apple Developer Program costs $99/year, and the project is currently distributed free of charge). As a result, macOS shows a **Gatekeeper** warning the first time the app is launched from a fresh machine:

> "NetstatUI cannot be opened because the developer cannot be verified."
> / "NetstatUI is from an unidentified developer."

This is **not malware detection** — it is purely a reputation-based warning for apps without a Developer ID signature. There are three ways to open it:

**Method 1 — Finder (GUI, recommended for most users):**
1. Locate `NetstatUI` in Finder.
2. **Right-click** (or Ctrl-click) the icon, choose **Open** from the context menu.
3. In the new dialog, click **Open**.

After that, macOS remembers your choice and a normal double-click will work going forward.

**Method 2 — System Settings (if double-click is blocked first):**
1. Try to open `NetstatUI` — it gets blocked with a generic dialog.
2. Open **System Settings → Privacy & Security**.
3. Scroll down — you'll see a section "NetstatUI was blocked from opening because it is not from an identified developer."
4. Click **Open Anyway**, then confirm with your Touch ID / password.

**Method 3 — CLI one-liner (for terminal users):**

```bash
xattr -d com.apple.quarantine /path/to/NetstatUI
./NetstatUI
```

`xattr -d com.apple.quarantine` removes the `com.apple.quarantine` extended attribute that macOS attaches to files downloaded from the internet. After this, Gatekeeper no longer flags the binary. (If the file isn't from a browser / curl, this attribute may not be present and the command will error — in that case use Method 1 or 2.)

If you want to suppress the warning permanently without joining the Developer Program, you can sign the binary yourself with `codesign --sign - NetstatUI` (ad-hoc signing still triggers Gatekeeper on first run, but at least the prompt is slightly less alarming). For production releases, an [Apple Developer ID + `notarytool` notarisation](https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution) is the only path to a Gatekeeper-clean app.

### Settings

Open the **⚙ Settings** dialog from the title bar:

- **General** — theme (auto / light / dark), locale (English / 简体中文), **IP geolocation (show / hide)**, density (compact / comfortable), refresh interval.
- **Advanced** — clear localStorage to reset everything.

---

## Platform Support

| Platform | Status                                                                                                          |
| -------- | --------------------------------------------------------------------------------------------------------------- |
| Windows  | ✅ Full — gopsutil reads via `GetExtendedTcpTable` / `GetExtendedUdpTable`; kill via `TerminateProcess`         |
| macOS    | ✅ Full — gopsutil reads via `sysctl`; kill via `SIGKILL`; locale via `osascript`                              |
| Linux    | ✅ Full — gopsutil reads via `/proc/net/{tcp,udp}{,6}`; kill via `SIGKILL`; locale via `$LANG`                  |

### Platform-specific notes

- **Windows** — Mica backdrop requires Windows 11 build 22621+. On older Windows the window uses an opaque background. See [First run on Windows (SmartScreen)](./README.md#first-run-on-windows-smartscreen) above for the unsigned-binary warning.
- **macOS** — first launch may prompt for Accessibility / Full Disk Access permissions (gopsutil reads process info via `libproc`). See [First run on macOS (Gatekeeper)](./README.md#first-run-on-macos-gatekeeper) above for the unsigned-binary warning.
- **Linux** — non-root users cannot see PIDs of sockets owned by other users (kernel restriction); run `sudo ./NetstatUI` for full visibility.

The architecture is pluggable: `services/netstat/`, `services/process/`, `services/kill/`, `services/system/` each define one `*_<os>.go` file per supported OS, selected at compile time via `//go:build` tags. Adding a new platform = `netstat_<os>.go` + `process_<os>.go` + `kill_<os>.go` + `system_<os>.go` + `netstat_platform_<os>.go` + entry in `main.go`'s `switch runtime.GOOS`.

---

## Known Limitations

- **Some loopback listeners may be missing on Windows 11 22H2+** — `iphlpapi.dll`'s `GetTcpTable2` / `GetExtendedUdpTable` silently drop a subset of `127.0.0.1` LISTEN entries. `netstat -ano` shows them because it uses WMI; we go through the same iphlpapi path as gopsutil, so the same limitation applies.
- **Mica backdrop is Windows-only** — the frameless translucent window falls back to an opaque background on macOS/Linux.
- **SmartScreen warning on first run (Windows)** — the published binary is not signed with a paid EV/OV certificate; see [First run on Windows](./README.md#first-run-on-windows-smartscreen) above.
- **Gatekeeper warning on first run (macOS)** — the published binary is not signed with an Apple Developer ID and not notarised; see [First run on macOS](./README.md#first-run-on-macos-gatekeeper) above.

See [`AGENTS.md` → Known pitfalls](./AGENTS.md#known-坑) for more.

---

## Project Layout

```
.
├── main.go                       # Wails app entry, registers services + events + platform switch
├── app.go                        # AppService: KillProcess / GetProcessDetail / OpenProcessFolder / GetSystemLocale
├── services/
│   ├── netstat/                  # TCP/UDP snapshot (shared adapter + per-OS provider)
│   │   ├── netstat_proto.go      # mapProto / mapState / gopsutilSnapshot (shared)
│   │   ├── netstat_windows.go    # gopsutil wrapper
│   │   ├── netstat_darwin.go     # gopsutil wrapper
│   │   └── netstat_linux.go      # gopsutil wrapper
│   ├── process/                  # PID → name/path cache (gopsutil on all OS)
│   ├── monitor/                  # Polling, diff, Wails event emit, geo lookup
│   ├── geo/                      # IP → Country-City via ip2region xdb (embedded)
│   ├── kill/                     # TerminateProcess (Win) / SIGKILL (Unix) via gopsutil
│   └── system/                   # GetSystemLocale (registry / osascript / $LANG)
├── data/                         # ip2region xdb files (gitignored, downloaded via scripts/download-xdb.ps1)
├── scripts/                      # download-xdb.ps1
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
│   ├── windows/info.json         # Windows resource metadata
│   └── darwin/Info.plist         # macOS bundle metadata
├── build.ps1                     # Safe Windows build (bypasses file-lock)
└── .github/workflows/build.yml   # CI: windows-amd64 + darwin-arm64 (Intel macOS via Rosetta 2)
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
