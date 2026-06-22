<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

interface DisplayStats {
  total: number
  listen: number
  established: number
  udp: number
}

const props = defineProps<{
  stats: DisplayStats
  filtered: number
  running: boolean
}>()

const items = computed(() => [
  { label: t('stats.total'), value: props.stats.total, color: 'var(--text)' },
  { label: t('stats.listen'), value: props.stats.listen, color: 'var(--state-listen)' },
  { label: t('stats.established'), value: props.stats.established, color: 'var(--state-established)' },
  { label: t('stats.udp'), value: props.stats.udp, color: 'var(--state-syn)' },
  { label: t('stats.filtered'), value: props.filtered, color: 'var(--accent)' },
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
        <span class="pulse" />{{ running ? t('stats.live') : t('stats.paused') }}
      </span>
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
