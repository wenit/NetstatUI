<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { Clipboard } from '@wailsio/runtime'
import type { ConnRow } from '../composables/useConnections'
import { AppService } from '../../bindings/github.com/zwb/network-ports'
import type { Info } from '../../bindings/github.com/zwb/network-ports/services/process/models'

const { t } = useI18n()

const props = defineProps<{ row: ConnRow | null }>()
const emit = defineEmits<{ close: [] }>()

const detail = ref<Info | null>(null)
const loading = ref(false)

watch(() => props.row, async (r) => {
  detail.value = null
  if (!r || !r.pid) return
  loading.value = true
  try {
    detail.value = await AppService.GetProcessDetail(r.pid)
  } finally {
    loading.value = false
  }
}, { immediate: true })

async function copyPath() {
  if (detail.value?.path) await Clipboard.SetText(detail.value.path)
}
async function killProc() {
  if (!props.row?.pid) return
  const r = await AppService.KillProcess(props.row.pid)
  if (!r.ok) alert(t('error.killFailed', { pid: props.row.pid, reason: r.reason }))
  else emit('close')
}
async function openFolder() {
  if (props.row?.pid) await AppService.OpenProcessFolder(props.row.pid)
}
</script>

<template>
  <transition name="slide">
    <div v-if="row" class="panel">
      <div class="header">
        <span class="title">{{ t('detail.title') }}</span>
        <button class="close-btn" @click="emit('close')">×</button>
      </div>
      <div class="body">
        <div class="section">
          <div class="sec-title">{{ t('detail.connection') }}</div>
          <div class="grid">
            <span class="k">{{ t('detail.protocol') }}</span><span class="v mono">{{ row.protocol }}</span>
            <span class="k">{{ t('detail.local') }}</span><span class="v mono">{{ row.localAddr }}:{{ row.localPort }}</span>
            <span class="k">{{ t('detail.remote') }}</span><span class="v mono">{{ row.remoteAddr || '*' }}{{ row.remotePort ? ':' + row.remotePort : '' }}</span>
            <span class="k">{{ t('detail.state') }}</span><span class="v"><span class="state-pill">{{ row.state }}</span></span>
          </div>
        </div>
        <div class="section">
          <div class="sec-title">{{ t('detail.process') }}</div>
          <div v-if="loading" class="loading">{{ t('detail.loading') }}</div>
          <div v-else-if="detail" class="grid">
            <span class="k">{{ t('detail.name') }}</span><span class="v">{{ detail.name || '—' }}</span>
            <span class="k">PID</span><span class="v mono">{{ detail.pid }}</span>
            <span class="k">{{ t('detail.ppid') }}</span><span class="v mono">{{ detail.ppid }}</span>
            <span class="k">{{ t('detail.path') }}</span>
            <span class="v path" :title="detail.path" @click="copyPath">
              {{ detail.path || '—' }}
              <span v-if="detail.path" class="copy-hint">{{ t('detail.clickToCopy') }}</span>
            </span>
          </div>
          <div v-else class="loading">{{ t('detail.noInfo', { pid: row.pid }) }}</div>
        </div>
        <div class="actions" v-if="row.pid">
          <button class="act" @click="openFolder">{{ t('detail.openFolder') }}</button>
          <button class="act danger" @click="killProc">{{ t('detail.kill') }}</button>
        </div>
      </div>
    </div>
  </transition>
</template>

<style scoped>
.panel {
  position: absolute;
  right: 0; top: 0; bottom: 0;
  width: 300px;
  background: var(--bg-elevated);
  backdrop-filter: blur(20px);
  border-left: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  z-index: 10;
}
.header {
  height: var(--filterbar-h);
  display: flex; align-items: center; justify-content: space-between;
  padding: 0 12px;
  border-bottom: 1px solid var(--border);
  flex-shrink: 0;
}
.title { font-size: 12px; font-weight: 600; color: var(--text); }
.close-btn {
  border: none; background: transparent; color: var(--text-tertiary);
  font-size: 18px; cursor: pointer; width: 24px; height: 24px;
  border-radius: 4px; display: flex; align-items: center; justify-content: center;
}
.close-btn:hover { background: var(--bg-hover); color: var(--text); }
.body { flex: 1; overflow-y: auto; padding: 12px; }
.section { margin-bottom: 18px; }
.sec-title {
  font-size: 11px; font-weight: 600; color: var(--text-tertiary);
  text-transform: uppercase; letter-spacing: 0.05em; margin-bottom: 8px;
}
.grid { display: grid; grid-template-columns: 70px 1fr; gap: 6px 8px; }
.k { color: var(--text-tertiary); font-size: 11px; align-self: center; }
.v { color: var(--text); font-size: 12px; word-break: break-all; }
.mono { font-family: var(--font-mono); }
.path { cursor: pointer; position: relative; }
.path:hover { color: var(--accent); }
.copy-hint { display: block; font-size: 10px; color: var(--text-tertiary); margin-top: 2px; }
.state-pill {
  font-size: 10px; padding: 2px 8px; border-radius: 10px;
  background: var(--bg-active); color: var(--text-secondary);
  font-family: var(--font-mono);
}
.loading { color: var(--text-tertiary); font-size: 12px; padding: 8px 0; }
.actions { display: flex; gap: 8px; margin-top: 12px; }
.act {
  flex: 1; height: 30px; border: 1px solid var(--border);
  background: var(--bg-active); color: var(--text);
  border-radius: var(--radius-md); font-size: 12px; cursor: pointer;
  transition: all var(--transition-fast);
}
.act:hover { background: var(--bg-hover); }
.act.danger { color: var(--danger); border-color: transparent; }
.act.danger:hover { background: var(--row-remove); }
.slide-enter-active, .slide-leave-active { transition: transform var(--transition); }
.slide-enter-from, .slide-leave-to { transform: translateX(100%); }
</style>
