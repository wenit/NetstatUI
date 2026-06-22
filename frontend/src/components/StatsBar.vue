<script setup lang="ts">
import { computed } from 'vue'
import type { Stats } from '../../bindings/github.com/zwb/network-ports/services/monitor/models'

const props = defineProps<{
  stats: Stats
  filtered: number
  running: boolean
  lastRefreshedAt: number
}>()

const timeStr = computed(() => {
  const d = new Date(props.lastRefreshedAt)
  return d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
})

const items = computed(() => [
  { label: '总计', value: props.stats.total, color: 'var(--text)' },
  { label: '监听', value: props.stats.listen, color: 'var(--state-listen)' },
  { label: '已建立', value: props.stats.established, color: 'var(--state-established)' },
  { label: 'UDP', value: props.stats.udp, color: 'var(--state-syn)' },
  { label: '过滤命中', value: props.filtered, color: 'var(--accent)' },
])
</script>

<template>
  <div class="statsbar">
    <div class="items">
      <div v-for="i in items" :key="i.label" class="item">
        <span class="dot" :style="{ background: i.color }" />
        <span class="label">{{ i.label }}</span>
        <span class="value" :style="{ color: i.color }">{{ i.value }}</span>
      </div>
    </div>
    <div class="status">
      <span class="live" :class="{ on: running }">
        <span class="pulse" />{{ running ? '实时' : '已暂停' }}
      </span>
      <span class="refresh-time">上次刷新 {{ timeStr }}</span>
    </div>
  </div>
</template>

<style scoped>
.statsbar {
  height: var(--statsbar-h);
  display: flex;
  align-items: center;
  padding: 0 12px;
  border-top: 1px solid var(--border);
  background: var(--bg-elevated);
  flex-shrink: 0;
  font-size: 11px;
}
.items { display: flex; gap: 18px; }
.item { display: flex; align-items: center; gap: 5px; }
.dot { width: 6px; height: 6px; border-radius: 50%; }
.label { color: var(--text-tertiary); }
.value { font-weight: 600; font-family: var(--font-mono); min-width: 14px; }
.status { margin-left: auto; }
.live { display: flex; align-items: center; gap: 6px; color: var(--text-tertiary); }
.live.on { color: var(--state-established); }
.refresh-time { color: var(--text-tertiary); margin-left: 12px; }
.pulse {
  width: 7px; height: 7px; border-radius: 50%;
  background: currentColor;
  animation: pulse 2s ease-in-out infinite;
}
@keyframes pulse {
  0%, 100% { opacity: 0.4; }
  50% { opacity: 1; }
}
</style>
