<script setup lang="ts">
import { computed } from 'vue'
import { useSettingsStore } from '../stores/settings'
import { refreshNow, setIntervalMs, startMonitor, stopMonitor } from '../composables/useConnections'

const props = defineProps<{ lastRefreshedAt: number }>()

const settings = useSettingsStore()

const intervals = [
  { label: '5s', ms: 5000 },
  { label: '15s', ms: 15000 },
  { label: '30s', ms: 30000 },
  { label: '60s', ms: 60000 },
]

const paused = computed(() => !settings.running)

async function pickInterval(ms: number) {
  settings.setInterval(ms)
  if (settings.running) await setIntervalMs(ms)
}

async function togglePause() {
  if (settings.running) {
    await stopMonitor()
    settings.setRunning(false)
  } else {
    await startMonitor(settings.intervalMs)
    settings.setRunning(true)
  }
}

async function doRefresh() {
  await refreshNow()
}
</script>

<template>
  <div class="toolbar">
    <div class="left">
      <span class="label">刷新</span>
      <div class="seg">
        <button
          v-for="i in intervals"
          :key="i.ms"
          class="seg-btn"
          :class="{ active: settings.intervalMs === i.ms && !paused }"
          @click="pickInterval(i.ms)"
        >{{ i.label }}</button>
      </div>
      <button class="icon-btn" :class="{ active: paused }" :title="paused ? '继续' : '暂停'" @click="togglePause">
        <svg v-if="paused" width="14" height="14" viewBox="0 0 14 14"><path d="M3 2v10l8-5z" fill="currentColor" /></svg>
        <svg v-else width="14" height="14" viewBox="0 0 14 14"><rect x="3" y="2" width="3" height="10" fill="currentColor" /><rect x="8" y="2" width="3" height="10" fill="currentColor" /></svg>
      </button>
      <button class="icon-btn" title="立即刷新" @click="doRefresh">
        <svg width="14" height="14" viewBox="0 0 14 14">
          <path d="M12 7a5 5 0 1 1-1.46-3.54" stroke="currentColor" stroke-width="1.2" fill="none" />
          <path d="M12 2v3h-3" stroke="currentColor" stroke-width="1.2" fill="none" stroke-linecap="round" />
        </svg>
      </button>
      <span class="refresh-time">{{ new Date(props.lastRefreshedAt).toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit', second: '2-digit' }) }}</span>
    </div>
    <div class="right">
      <slot />
    </div>
  </div>
</template>

<style scoped>
.toolbar {
  height: var(--toolbar-h);
  display: flex;
  align-items: center;
  padding: 0 12px;
  gap: 12px;
  border-bottom: 1px solid var(--border);
  flex-shrink: 0;
}
.left { display: flex; align-items: center; gap: 6px; }
.right { margin-left: auto; display: flex; align-items: center; gap: 8px; }
.label { font-size: 12px; color: var(--text-tertiary); margin-right: 2px; }
.seg {
  display: flex;
  background: var(--bg-active);
  border-radius: var(--radius-md);
  padding: 2px;
}
.seg-btn {
  height: 24px;
  padding: 0 10px;
  border: none;
  background: transparent;
  color: var(--text-secondary);
  font-size: 12px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all var(--transition-fast);
}
.seg-btn:hover { color: var(--text); }
.seg-btn.active { background: var(--accent); color: #fff; }
.icon-btn {
  height: 28px;
  width: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  color: var(--text-secondary);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all var(--transition-fast);
}
.icon-btn:hover { background: var(--bg-hover); color: var(--text); }
.icon-btn.active { color: var(--state-listen); }
.refresh-time { font-size: 11px; color: var(--text-tertiary); font-family: var(--font-mono); white-space: nowrap; }
</style>
