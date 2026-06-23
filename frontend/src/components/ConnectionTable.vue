<script setup lang="ts">
import { computed, ref, shallowRef, watch, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import type { ConnRow } from '../composables/useConnections'
import type { HighlightState } from '../composables/useConnections'
import { AppService } from '../../bindings/github.com/zwb/network-ports'
import type { State } from '../../bindings/github.com/zwb/network-ports/services/netstat/models'
import { showError } from '../composables/useErrorDialog'
import { showConfirm } from '../composables/useConfirmDialog'

const { t } = useI18n()

const props = defineProps<{
  rows: ConnRow[]
  highlights: HighlightState[]
  density: 'compact' | 'comfortable'
}>()

const emit = defineEmits<{
  select: [ConnRow | null]
  contextmenu: [event: MouseEvent, row: ConnRow]
}>()

type SortKey = 'protocol' | 'localAddr' | 'localPort' | 'remoteAddr' | 'remotePort' | 'state' | 'pid' | 'processName'
const sortKey = ref<SortKey>('localPort')
const sortDesc = ref(false)

const rowH = computed(() => props.density === 'compact' ? 26 : 30)
const headerH = 30

const containerEl = ref<HTMLElement | null>(null)
const scrollTop = ref(0)
const viewportH = ref(600)

const highlightMap = computed(() => {
  const m = new Map<string, 'add' | 'remove' | 'update'>()
  const now = Date.now()
  for (const h of props.highlights) {
    if (h.until > now) m.set(h.key, h.type)
  }
  return m
})

const sorted = computed(() => {
  const arr = props.rows
  const k = sortKey.value
  const desc = sortDesc.value
  if (arr.length <= 1) return arr
  const idx = arr.slice()
  idx.sort((a, b) => {
    let r = 0
    switch (k) {
      case 'protocol': r = a.protocol.localeCompare(b.protocol); break
      case 'localAddr': r = a.localAddr.localeCompare(b.localAddr); break
      case 'localPort': r = a.localPort - b.localPort; break
      case 'remoteAddr': r = a.remoteAddr.localeCompare(b.remoteAddr); break
      case 'remotePort': r = a.remotePort - b.remotePort; break
      case 'state': r = a.state.localeCompare(b.state); break
      case 'pid': r = a.pid - b.pid; break
      case 'processName': r = (a.processName || '').localeCompare(b.processName || ''); break
    }
    return desc ? -r : r
  })
  return idx
})

const totalH = computed(() => sorted.value.length * rowH.value + 2)
const startIdx = computed(() => Math.max(0, Math.floor(scrollTop.value / rowH.value) - 4))
const endIdx = computed(() => Math.min(sorted.value.length, Math.ceil((scrollTop.value + viewportH.value) / rowH.value) + 4))
const visibleRows = computed(() => {
  const out: { row: ConnRow; top: number; index: number }[] = []
  const list = sorted.value
  for (let i = startIdx.value; i < endIdx.value; i++) {
    out.push({ row: list[i], top: i * rowH.value, index: i })
  }
  return out
})

function onScroll() {
  if (containerEl.value) scrollTop.value = containerEl.value.scrollTop
}

function updateViewport() {
  if (containerEl.value) viewportH.value = containerEl.value.clientHeight
}

let ro: ResizeObserver | null = null
onMounted(() => {
  updateViewport()
  if (containerEl.value && typeof ResizeObserver !== 'undefined') {
    ro = new ResizeObserver(updateViewport)
    ro.observe(containerEl.value)
  }
})
onUnmounted(() => { ro?.disconnect() })

function toggleSort(k: SortKey) {
  if (sortKey.value === k) sortDesc.value = !sortDesc.value
  else { sortKey.value = k; sortDesc.value = false }
}

function stateColor(s: State): string {
  const map: Record<string, string> = {
    'LISTEN': 'var(--state-listen)',
    'ESTABLISHED': 'var(--state-established)',
    'TIME_WAIT': 'var(--state-time-wait)',
    'CLOSE_WAIT': 'var(--state-close-wait)',
    'FIN_WAIT_1': 'var(--state-close-wait)',
    'FIN_WAIT_2': 'var(--state-time-wait)',
    'SYN_SENT': 'var(--state-syn)',
    'SYN_RECEIVED': 'var(--state-syn)',
    'CLOSING': 'var(--state-time-wait)',
  }
  return map[s as string] || 'var(--text-tertiary)'
}

function rowBg(row: ConnRow): string {
  const hl = highlightMap.value.get(row.key)
  if (hl === 'add') return 'var(--row-add)'
  if (hl === 'remove') return 'var(--row-remove)'
  if (hl === 'update') return 'var(--row-update)'
  return ''
}

const selectedKey = ref<string | null>(null)
function selectRow(row: ConnRow) {
  selectedKey.value = row.key
  emit('select', row)
}

let lastKill = 0
async function killRow(row: ConnRow, ev: MouseEvent) {
  ev.stopPropagation()
  const now = Date.now()
  if (now - lastKill < 800) return
  lastKill = now
  const ok = await showConfirm(t('confirm.killTitle'), t('confirm.killBody', { pid: row.pid }))
  if (!ok) return
  const r = await AppService.KillProcess(row.pid)
  if (!r.ok) {
    showError(t('error.killFailedTitle'), t('error.killFailed', { pid: row.pid, reason: r.reason }))
  }
}

watch(() => props.rows, () => {
  if (containerEl.value && scrollTop.value > totalH.value) {
    containerEl.value.scrollTop = 0
    scrollTop.value = 0
  }
})
</script>

<template>
  <div class="table-wrap">
    <div class="thead">
      <div class="th protocol" @click="toggleSort('protocol')">
        {{ t('table.colProtocol') }}<span class="arrow" :class="{ show: sortKey==='protocol', desc: sortDesc }">▾</span>
      </div>
      <div class="th flex" @click="toggleSort('localAddr')">
        {{ t('table.colLocalAddr') }}<span class="arrow" :class="{ show: sortKey==='localAddr', desc: sortDesc }">▾</span>
      </div>
      <div class="th port" @click="toggleSort('localPort')">
        {{ t('table.colLocalPort') }}<span class="arrow" :class="{ show: sortKey==='localPort', desc: sortDesc }">▾</span>
      </div>
      <div class="th flex" @click="toggleSort('remoteAddr')">
        {{ t('table.colRemoteAddr') }}<span class="arrow" :class="{ show: sortKey==='remoteAddr', desc: sortDesc }">▾</span>
      </div>
      <div class="th port" @click="toggleSort('remotePort')">
        {{ t('table.colRemotePort') }}<span class="arrow" :class="{ show: sortKey==='remotePort', desc: sortDesc }">▾</span>
      </div>
      <div class="th state" @click="toggleSort('state')">
        {{ t('table.colState') }}<span class="arrow" :class="{ show: sortKey==='state', desc: sortDesc }">▾</span>
      </div>
      <div class="th pid" @click="toggleSort('pid')">
        PID<span class="arrow" :class="{ show: sortKey==='pid', desc: sortDesc }">▾</span>
      </div>
      <div class="th flex proc" @click="toggleSort('processName')">
        {{ t('table.colProcess') }}<span class="arrow" :class="{ show: sortKey==='processName', desc: sortDesc }">▾</span>
      </div>
      <div class="th kill" />
    </div>
    <div ref="containerEl" class="tbody" @scroll="onScroll">
      <div class="spacer" :style="{ height: totalH + 'px' }">
        <div
          v-for="v in visibleRows"
          :key="v.row.key"
          class="tr"
          :class="{ alt: v.index % 2, selected: selectedKey === v.row.key }"
          :style="{ transform: `translateY(${v.top}px)`, background: rowBg(v.row) }"
          @click="selectRow(v.row)"
          @dblclick="killRow(v.row, $event)"
          @contextmenu="emit('contextmenu', $event, v.row)"
        >
          <div class="td protocol" :class="v.row.protocol">{{ v.row.protocol }}</div>
          <div class="td flex mono">{{ v.row.localAddr }}</div>
          <div class="td port mono">{{ v.row.localPort }}</div>
          <div class="td flex mono">{{ v.row.remoteAddr || '*' }}<span v-if="v.row.remotePort" class="dim">:{{ v.row.remotePort }}</span></div>
          <div class="td port mono">{{ v.row.remotePort || '—' }}</div>
          <div class="td state"><span class="dot" :style="{ background: stateColor(v.row.state) }" />{{ v.row.state }}</div>
          <div class="td pid mono">{{ v.row.pid || '—' }}</div>
          <div class="td flex proc" :title="v.row.processPath">{{ v.row.processName || (v.row.pid ? t('table.system') : '—') }}</div>
          <div class="td kill" @click.stop="killRow(v.row, $event)">
            <span class="kill-btn" :title="t('table.kill')">×</span>
          </div>
        </div>
        <div v-if="rows.length === 0" class="empty">
          <span>{{ t('table.empty') }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.table-wrap { flex: 1; display: flex; flex-direction: column; min-height: 0; }
.thead {
  display: flex;
  align-items: center;
  height: 30px;
  border-bottom: 1px solid var(--border);
  background: var(--bg-elevated);
  flex-shrink: 0;
  font-size: 11px;
  color: var(--text-secondary);
  user-select: none;
}
.th {
  display: flex;
  align-items: center;
  padding: 0 8px;
  height: 100%;
  cursor: pointer;
  white-space: nowrap;
  gap: 4px;
}
.th:hover { color: var(--text); }
.th.flex { flex: 1; min-width: 0; }
.th.protocol { width: 64px; }
.th.port { width: 72px; justify-content: flex-end; }
.th.state { width: 120px; }
.th.pid { width: 76px; justify-content: flex-end; }
.th.kill { width: 32px; padding: 0; justify-content: center; cursor: default; }
.th.kill:hover { color: var(--text-secondary); }
.arrow { opacity: 0; font-size: 9px; transition: transform var(--transition-fast); }
.arrow.show { opacity: 0.7; }
.arrow.desc { transform: rotate(180deg); }
.tbody {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  position: relative;
}
.spacer { position: relative; width: 100%; }
.tr {
  position: absolute;
  left: 0; right: 0;
  height: v-bind('rowH + "px"');
  display: flex;
  align-items: center;
  font-size: 12px;
  color: var(--text);
  cursor: pointer;
  border-bottom: 1px solid transparent;
  will-change: transform;
}
.tr:hover { background: var(--bg-hover) !important; }
.tr.alt { background: var(--bg-row-alt); }
.tr.selected { box-shadow: inset 2px 0 0 var(--accent); }
.td {
  padding: 0 8px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  height: 100%;
  display: flex;
  align-items: center;
}
.td.flex { flex: 1; min-width: 0; }
.td.protocol { width: 64px; font-size: 11px; font-weight: 500; }
.td.protocol.tcp4, .td.protocol.tcp6 { color: var(--state-syn); }
.td.protocol.udp4, .td.protocol.udp6 { color: var(--state-listen); }
.td.port { width: 72px; justify-content: flex-end; }
.td.state { width: 120px; gap: 6px; font-size: 11px; }
.td.state .dot { width: 7px; height: 7px; border-radius: 50%; flex-shrink: 0; }
.td.pid { width: 76px; justify-content: flex-end; }
.td.kill { width: 32px; justify-content: center; padding: 0; }
.mono { font-family: var(--font-mono); font-size: 11.5px; }
.dim { color: var(--text-tertiary); }
.proc { color: var(--text-secondary); }
.kill-btn {
  width: 18px; height: 18px;
  display: flex; align-items: center; justify-content: center;
  border-radius: 4px;
  color: var(--text-tertiary);
  font-size: 16px; line-height: 1;
  cursor: pointer;
  transition: all var(--transition-fast);
}
.kill-btn:hover { background: var(--danger); color: #fff; }
.empty {
  position: absolute; inset: 0;
  display: flex; align-items: center; justify-content: center;
  color: var(--text-tertiary); font-size: 13px;
}
</style>
