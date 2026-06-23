<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Clipboard } from '@wailsio/runtime'
import type { ConnRow } from '../composables/useConnections'
import { AppService } from '../../bindings/github.com/zwb/network-ports'
import { showError } from '../composables/useErrorDialog'
import { showConfirm } from '../composables/useConfirmDialog'

const { t } = useI18n()

const props = defineProps<{
  x: number
  y: number
  row: ConnRow
}>()
const emit = defineEmits<{ close: [] }>()

function rowCSV(r: ConnRow): string {
  return [r.protocol, r.localAddr, r.localPort, r.remoteAddr, r.remotePort, r.state, r.pid, r.processName, r.processPath].join(',')
}

async function copy(text: string) {
  await Clipboard.SetText(text)
  emit('close')
}
function copyCell() { copy(`${props.row.localAddr}:${props.row.localPort}`) }
function copyRow() { copy(rowCSV(props.row)) }
async function copyAll() {
  const r = props.row
  await Clipboard.SetText(`Protocol,LocalAddr,LocalPort,RemoteAddr,RemotePort,State,PID,Process,Path\n${rowCSV(r)}`)
  emit('close')
}
async function openFolder() {
  await AppService.OpenProcessFolder(props.row.pid)
  emit('close')
}
async function killProc() {
  const ok = await showConfirm(t('confirm.killTitle'), t('confirm.killBody', { pid: props.row.pid }))
  if (!ok) return
  const r = await AppService.KillProcess(props.row.pid)
  if (!r.ok) showError(t('error.killFailedTitle'), t('error.killFailed', { pid: props.row.pid, reason: r.reason }))
  emit('close')
}

function onDocClick() { emit('close') }
onMounted(() => {
  setTimeout(() => document.addEventListener('click', onDocClick), 0)
  document.addEventListener('contextmenu', onDocClick)
})
onUnmounted(() => {
  document.removeEventListener('click', onDocClick)
  document.removeEventListener('contextmenu', onDocClick)
})
</script>

<template>
  <div class="ctx-menu" :style="{ left: x + 'px', top: y + 'px' }" @click.stop>
    <button class="mi" @click="copyCell"><svg width="14" height="14" viewBox="0 0 14 14"><rect x="3" y="3" width="6" height="6" stroke="currentColor" fill="none" /><rect x="5" y="5" width="6" height="6" stroke="currentColor" fill="none" /></svg>{{ t('context.copyAddrPort') }}</button>
    <button class="mi" @click="copyRow"><svg width="14" height="14" viewBox="0 0 14 14"><rect x="2" y="3" width="10" height="8" stroke="currentColor" fill="none" rx="1" /></svg>{{ t('context.copyRowCsv') }}</button>
    <button class="mi" @click="copyAll"><svg width="14" height="14" viewBox="0 0 14 14"><rect x="2" y="3" width="10" height="8" stroke="currentColor" fill="none" rx="1" /><line x1="4" y1="6" x2="10" y2="6" stroke="currentColor" /><line x1="4" y1="8" x2="10" y2="8" stroke="currentColor" /></svg>{{ t('context.copyWithHeader') }}</button>
    <div class="divider" />
    <button class="mi" :disabled="!row.pid" @click="openFolder"><svg width="14" height="14" viewBox="0 0 14 14"><path d="M2 4h3l1.5 1.5H12v6H2z" stroke="currentColor" fill="none" /></svg>{{ t('context.openFolder') }}</button>
    <button class="mi danger" :disabled="!row.pid" @click="killProc"><svg width="14" height="14" viewBox="0 0 14 14"><circle cx="7" cy="7" r="5" stroke="currentColor" fill="none" /><line x1="5" y1="5" x2="9" y2="9" stroke="currentColor" /><line x1="9" y1="5" x2="5" y2="9" stroke="currentColor" /></svg>{{ t('context.kill') }}</button>
  </div>
</template>

<style scoped>
.ctx-menu {
  position: fixed;
  z-index: 1000;
  min-width: 180px;
  background: var(--bg-elevated);
  backdrop-filter: blur(20px);
  border: 1px solid var(--border-strong);
  border-radius: var(--radius-md);
  padding: 4px;
  box-shadow: 0 8px 24px rgba(0,0,0,0.28);
}
.mi {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  height: 30px;
  padding: 0 10px;
  border: none;
  background: transparent;
  color: var(--text);
  font-size: 12px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  text-align: left;
}
.mi:hover:not(:disabled) { background: var(--bg-hover); }
.mi:disabled { opacity: 0.4; cursor: default; }
.mi.danger { color: var(--danger); }
.mi.danger:hover:not(:disabled) { background: var(--row-remove); }
.mi svg { color: var(--text-secondary); flex-shrink: 0; }
.mi.danger svg { color: var(--danger); }
.divider { height: 1px; background: var(--border); margin: 4px 0; }
</style>
