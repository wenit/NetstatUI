# AGENTS.md

Wails 3 (alpha.98) + Vue 3 + TypeScript 桌面应用：Win11 Fluent 风格的网络端口/连接查看器（类似 Linux netstat）。仅 Windows，只读查询 + 结束进程。

## 必读入口
- `main.go` — Wails 3 application 入口：注册服务、创建 Frameless + Mica 窗口、注册事件类型（`conn:full/conn:diff/conn:stats/conn:error`）。
- `app.go` — `AppService`：绑定给前端的方法 `KillProcess` / `GetProcessDetail` / `OpenProcessFolder` / `GetSystemLocale`。
- `services/netstat/` — 跨平台端口查询抽象（`Provider` 接口）+ Windows 实现（`go-netstat` 库，内部 `GetTcpTable2` / `GetTcp6Table2` / `GetExtendedUdpTable`）。
- `services/process/` — PID → 进程信息缓存（`CreateToolhelp32Snapshot` + `QueryFullProcessImageNameW`）。LRU + path TTL 30s。
- `services/monitor/` — 轮询调度 + diff 计算 + Wails 事件推送。`GetSnapshot`（首帧）/ `Start`/`Stop`/`SetInterval`/`RefreshNow`。
- `services/kill/` — `windows.TerminateProcess` 结束进程（仅 `PROCESS_TERMINATE` 权限）。
- `services/system/` — `GetSystemLocale()` 返回 OS locale 字符串。
- `frontend/src/composables/useConnections.ts` — 订阅 Wails 事件 + diff 应用到 `shallowRef<Map>` + 变化高亮。
- `frontend/src/components/ConnectionTable.vue` — 自定义虚拟滚动表格（绝对定位 + transform）。

## 关键不变量
- **窗口**：`Frameless: true` + `BackgroundType: Translucent` + `Windows.BackdropType: Mica`（Win11 22621+）。标题栏拖拽靠 CSS `--wails-draggable: drag`（在 `TitleBar.vue` 的 `.drag` 元素上），不是 JS API。
- **端口字节序**：`dwLocalPort` 是网络序存低 16 位，`ntohsLow(v) = ((v&0xFF)<<8)|((v>>8)&0xFF)`。IPv4 `dwLocalAddr` 网络序，`ipv4String` 按 `v&0xFF → (v>>24)&0xFF` 取字节。
- **IPv6 结构体**：`mibTcp6RowOwnerPid` / `mibUdp6RowOwnerPid` 必须含 `LocalScopeId` / `RemoteScopeId`，否则偏移错位导致地址乱码 + PID 巨大。
- **连接唯一键**：`Key = protocol|local:port|remote:port|pid`，diff 以此为索引。改键格式会破坏增量更新。
- **首帧竞态**：monitor `ServiceStartup` 会 emit `conn:full`，但前端可能还没 mount。前端 `onMounted` 调 `fetchSnapshot()`（→ `Monitor.GetSnapshot()`）保证拿到初始数据。不要删 `GetSnapshot`。
- **前端数据流**：`conns` 是 `shallowRef<Map<string, ConnRow>>`，每次更新替换整个 Map（不深响应）。`watch(conns, recompute)` 触发过滤重算。虚拟滚动靠 `sorted` computed + `visibleRows` 切片。
- **bindings 是生成物**：`frontend/bindings/` 由 `wails3 generate bindings -ts` 生成，不要手改。Go 侧改了方法/结构体后必须重新生成 + 重新 npm run build。
- **diff 比较**：monitor `diff()` 只对比 `State` 和 `ProcessName` 变化，其他字段变化不触发 update 事件。

## 常用命令
```bash
wails3 dev                  # 开发模式（vite dev server + 热重载）
wails3 build                # 生产构建 → bin/NetstatUI.exe
wails3 generate bindings -ts -d frontend/bindings -names   # 重新生成 TS bindings
go test ./services/netstat/ -run TestWindowsSnapshot -v    # 测试端口查询
go build . ./services/...   # 只编译 Go（不打包前端）
cd frontend && npm run build:dev   # vue-tsc + vite (development, sourcemap)
cd frontend && npm run build       # vue-tsc + vite (production)
.\build.ps1                 # Windows 安全构建（绕过 SearchIndexer/Defender 锁文件问题）
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

## 事件清单
| 事件 | 载荷 | 触发 |
|------|------|------|
| `conn:full` | `ConnInfo[]` | 启动 + `RefreshNow` + `GetSnapshot` |
| `conn:diff` | `Diff{added, removed, updated}` | 每个轮询周期 |
| `conn:stats` | `Stats{total, listen, established, udp, filteredHits}` | 每次收集 |
| `conn:error` | `string` | 采集失败 |

## 已知坑
- **GetTcpTable2 漏报 loopback listener**（如 127.0.0.1:9245）：Win11 22H2+ 上 iphlpapi.dll 会静默丢弃部分 loopback LISTEN 连接。不是 bug，是 Win11 限制。
- `go build ./...` 会尝试编译 `build/ios`（独立 main 包）报错；用 `go build . ./services/...` 代替。
- `//go:embed all:frontend/dist` 需要 `frontend/dist` 存在；首次 clone 后先 `npm run build` 再 `go build`。
- `wails3 build` 全流程（bindings → 前端 → Go）。Windows 上可能因 SearchIndexer/Defender 锁文件导致 `RemoveAll+Rename` 失败。用 `build.ps1` 以 `-clean=false` 绕过。
- `vue-virtual-scroller` 已安装但未使用——表格用的是自写虚拟滚动，不要混用。
- `_test.go` 文件在 `wails3 dev` 的 watch 中被忽略（见 `build/config.yml`）。
- localStorage 存储键前缀 `np.`（`np.theme / np.locale / np.interval / np.running`），清除 localStorage 可重置设置。
- vite.config.ts 有 `@vueuse/core` INVALID_ANNOTATION 警告抑制（Rolldown 解析器比 esbuild 更严格）。

## 扩展指引
- **新增列**：在 `ConnectionTable.vue` 的 `thead` + `tr` 加 `<div class="th/td">`，同步 `SortKey` 类型 + `sorted` 的 switch。
- **新增筛选维度**：在 `FilterBar.vue` 加 chip + `useFilters.ts` 的 `applyFilter` 加判断。
- **跨平台**：在 `services/netstat/` 加 `netstat_linux.go`（读 `/proc/net/tcp{,6}`）+ `netstat_darwin.go`（`sysctl`）；`process/` 加对应平台实现；`SetProvider` 按平台选择。接口已抽象好。
- **新增后端方法**：在 `app.go` 或 `monitor.go` 加方法 → `wails3 generate bindings -ts` → 前端从 bindings 导入调用。
- **新增 locale**：在 `frontend/src/locales/` 加 `<locale>.ts`，`index.ts` 的 `messages` 注册，`SUPPORTED_LOCALES` 加类型。
