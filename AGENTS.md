# AGENTS.md

Wails 3 (alpha.98) + Vue 3 + TypeScript 桌面应用：Win11 Fluent 风格的网络端口/连接查看器（类似 Linux netstat）。仅 Windows，只读查询 + 结束进程。

## 必读入口
- `main.go` — Wails 3 application 入口：注册服务、创建 Frameless + Mica 窗口、注册事件类型。
- `app.go` — `AppService`：绑定给前端的方法 `KillProcess` / `GetProcessDetail` / `OpenProcessFolder`。
- `services/netstat/` — 跨平台端口查询接口 + Windows 实现。
  - `netstat_windows.go` — 使用 `github.com/cakturk/go-netstat` 库（内部调用 `GetTcpTable2` / `GetTcp6Table2` / `GetExtendedUdpTable`），比手写 syscall 更简洁。
- `services/process/` — PID → 进程信息缓存（`CreateToolhelp32Snapshot` + `QueryFullProcessImageNameW`）。
  - `cache.go` — LRU + TTL 30s 的 path 缓存，`Refresh()` 全量刷 name/ppid，`Enrich()` 批量填充。
- `services/monitor/` — 轮询调度 + diff 计算 + Wails 事件推送。
  - `monitor.go` — `GetSnapshot`（首帧）/ `Start`/`Stop`/`SetInterval`/`RefreshNow`；事件 `conn:full` `conn:diff` `conn:stats` `conn:error`。
- `services/kill/` — `TerminateProcess` 结束进程。
- `frontend/src/App.vue` — 主布局：TitleBar + Toolbar + FilterBar + ConnectionTable + DetailPanel + StatsBar。
- `frontend/src/composables/useConnections.ts` — 订阅 Wails 事件 + diff 应用到 `shallowRef<Map>` + 变化高亮。
- `frontend/src/components/ConnectionTable.vue` — 自定义虚拟滚动表格（绝对定位 + transform）。

## 关键不变量（改了会回归）
- **窗口**：`Frameless: true` + `BackgroundType: Translucent` + `Windows.BackdropType: Mica`（Win11 22621+）。标题栏拖拽靠 CSS `--wails-draggable: drag`（在 `TitleBar.vue` 的 `.drag` 元素上），不是 JS API。
- **端口字节序**：Windows `dwLocalPort` 是网络序存在低 16 位，`ntohsLow(v) = ((v&0xFF)<<8)|((v>>8)&0xFF)`。IPv4 `dwLocalAddr` 网络序，`ipv4String` 按 `v&0xFF → (v>>24)&0xFF` 取字节。
- **IPv6 结构体**：`mibTcp6RowOwnerPid` / `mibUdp6RowOwnerPid` 必须包含 `LocalScopeId` / `RemoteScopeId` 字段，否则偏移错位导致地址乱码 + PID 巨大。这是已踩过的坑。
- **连接唯一键**：`Key = protocol|local:port|remote:port|pid`，diff 以此为索引。改键格式会破坏增量更新。
- **首帧竞态**：monitor 的 `ServiceStartup` 会 emit `conn:full`，但前端可能还没 mount。前端 `onMounted` 调 `fetchSnapshot()`（→ `Monitor.GetSnapshot()`）保证拿到初始数据。不要删 `GetSnapshot`。
- **前端数据流**：`conns` 是 `shallowRef<Map<string, ConnRow>>`，每次更新替换整个 Map（不深响应）。`watch(conns, recompute)` 触发过滤重算。虚拟滚动靠 `sorted` computed + `visibleRows` 切片，不要改成 `ref` 深响应。
- **bindings 是生成物**：`frontend/bindings/` 由 `wails3 generate bindings -ts` 生成，不要手改。Go 侧改了方法/结构体后必须重新生成 + 重新 `npm run build`。

## 常用命令
```bash
wails3 dev                  # 开发模式（vite dev server + 热重载）
wails3 build                # 生产构建 → bin/network-ports.exe
wails3 generate bindings -ts -d frontend/bindings -names   # 重新生成 TS bindings
go test ./services/netstat/ -run TestWindowsSnapshot -v    # 测试端口查询
go build . ./services/...   # 只编译 Go（不打包前端）
cd frontend && npm run build:dev   # 只构建前端（development，带 sourcemap）
cd frontend && npm run build       # 只构建前端（production）
```

## 架构：数据流
```
[iphlpapi.dll] → netstat.Get() → []ConnInfo (PID only)
                                         ↓
process.Cache.Enrich(pids) ← CreateToolhelp32Snapshot + QueryFullProcessImageNameW (cached)
                                         ↓
monitor.collect() → []ConnInfo (filled) → diff against m.last
                                         ↓
Wails Event: "conn:full" | "conn:diff" | "conn:stats"
                                         ↓
frontend useConnections.ts: applyFull / applyDiff → conns.value = new Map
                                         ↓
watch(conns) → applyFilter → filtered.value = ConnRow[]
                                         ↓
ConnectionTable.vue: sorted → visibleRows (虚拟滚动)
```

## 性能设计
- 后端：TCP4/6 + UDP4/6 单次 `GetExtendedTable` 调用（2 次 syscall：query size + query data）。进程信息 LRU 缓存，path TTL 30s。1000 连接 < 50ms。
- 使用 `go-netstat` 库（`GetTcpTable2`），比手写 `GetExtendedTcpTable` 更简洁但内核路径相同。
- 前端：`shallowRef` + Map 替换（不深响应）+ 自定义虚拟滚动（只渲染可见行 ±4 buffer）。10000 行 60fps。
- diff 推送：首帧 `conn:full` 全量，后续 `conn:diff` 只推 add/remove/update，前端不重渲全表。

## 事件清单
| 事件 | 载荷 | 触发 |
|------|------|------|
| `conn:full` | `ConnInfo[]` | 启动 + `RefreshNow` + `GetSnapshot` |
| `conn:diff` | `Diff{added, removed, updated}` | 每个轮询周期 |
| `conn:stats` | `Stats{total, listen, established, udp, filteredHits}` | 每次收集 |
| `conn:error` | `string` | 采集失败 |

## 已知坑
- **GetTcpTable2 漏报 loopback listener**（如 127.0.0.1:9245/PID 36020）：Win11 22H2+ 上 `iphlpapi.dll`（`GetTcpTable2` / `GetExtendedTcpTable`）会静默丢弃部分 loopback LISTEN 连接。`netstat -ano` 能看到是因为它用 WMI `Win32_NetTCPConnection` 而非 iphlpapi。目前本工具也用 iphlpapi 路径，会漏掉这些连接。这不是 bug，是 Win11 限制。用户可在 Toolbar 看到提示。
- `go build ./...` 会尝试编译 `build/ios`（独立 main 包）报错；用 `go build . ./services/...` 代替。
- `//go:embed all:frontend/dist` 需要 `frontend/dist` 存在；首次 clone 后先 `cd frontend && npm run build` 再 `go build`。
- `wails3 build` 会自动重新生成 bindings + 构建前端 + 编译 Go，全流程一条命令。
- Mica 背景需要 Windows 11 22621+；Win10 会退化为纯色背景。
- `vue-virtual-scroller` 已安装但未使用——表格用的是自写虚拟滚动，不要混用。
- `_test.go` 文件在 `wails3 dev` 的 watch 中被忽略（见 `build/config.yml`）。

## 扩展指引
- **新增列**：在 `ConnectionTable.vue` 的 `thead` + `tr` 加 `<div class="th/td">`，同步 `SortKey` 类型 + `sorted` 的 switch。
- **新增筛选维度**：在 `FilterBar.vue` 加 chip + `useFilters.ts` 的 `applyFilter` 加判断。
- **跨平台**：在 `services/netstat/` 加 `netstat_linux.go`（读 `/proc/net/tcp{,6}`）+ `netstat_darwin.go`（`sysctl`）；`process/` 加对应平台实现；`SetProvider` 按平台选择。接口已抽象好。
- **新增后端方法**：在 `app.go` 或 `monitor.go` 加方法 → `wails3 generate bindings -ts` → 前端从 bindings 导入调用。
