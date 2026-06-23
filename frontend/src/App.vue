<script setup lang="ts">
import { computed, onMounted, ref, shallowRef, watch } from 'vue'
import { usePreferredDark } from '@vueuse/core'
import { applyTheme, useSettingsStore } from './stores/settings'
import { setI18nLocale, matchLocale, type AppLocale } from './locales'
import { AppService } from '../bindings/github.com/zwb/network-ports'
import {
  initConnections, useConnections,
  startMonitor, stopMonitor, setIntervalMs, fetchSnapshot,
  type ConnRow,
} from './composables/useConnections'
import { applyFilter, emptyFilter, type FilterState } from './composables/useFilters'
import TitleBar from './components/TitleBar.vue'
import Toolbar from './components/Toolbar.vue'
import FilterBar from './components/FilterBar.vue'
import ConnectionTable from './components/ConnectionTable.vue'
import StatsBar from './components/StatsBar.vue'
import DetailPanel from './components/DetailPanel.vue'
import ContextMenu from './components/ContextMenu.vue'
import SettingsDialog from './components/SettingsDialog.vue'
import ErrorDialog from './components/ErrorDialog.vue'
import ConfirmDialog from './components/ConfirmDialog.vue'

const settings = useSettingsStore()
applyTheme(settings.theme)
setI18nLocale(settings.locale)

const { conns, stats, highlights, error, lastRefreshedAt } = useConnections()

const filter = ref<FilterState>(emptyFilter())
const filtered = shallowRef<ConnRow[]>([])
const settingsOpen = ref(false)

const selected = ref<ConnRow | null>(null)
const ctxMenu = ref<{ x: number; y: number; row: ConnRow } | null>(null)

function recompute() {
  filtered.value = applyFilter(conns.value, filter.value)
}

watch(conns, recompute, { flush: 'post' })
watch(filter, recompute, { deep: true, flush: 'post' })

function onSelect(row: ConnRow | null) {
  selected.value = row
}
function onContextmenu(ev: MouseEvent, row: ConnRow) {
  ev.preventDefault()
  ctxMenu.value = { x: ev.clientX, y: ev.clientY, row }
}

const preferredDark = usePreferredDark()
watch(preferredDark, () => {
  if (settings.theme === 'auto') applyTheme('auto')
})

onMounted(async () => {
  if (!localStorage.getItem('np.locale')) {
    try {
      const sysLocale = await AppService.GetSystemLocale()
      if (sysLocale) {
        const matched = matchLocale(sysLocale)
        if (matched !== settings.locale) {
          settings.setLocale(matched as AppLocale)
        }
      }
    } catch {}
  }

  initConnections()
  await fetchSnapshot()
  if (settings.running) {
    await startMonitor(settings.intervalMs)
  }
})

watch(() => settings.intervalMs, async (ms) => {
  if (settings.running) await setIntervalMs(ms)
})

watch(() => settings.running, async (r) => {
  if (r) { await startMonitor(settings.intervalMs) }
  else { await stopMonitor() }
})

watch(() => settings.locale, (l) => {
  setI18nLocale(l)
})

const filteredCount = computed(() => filtered.value.length)
const errorText = computed(() => error.value)

const displayStats = computed(() => {
  const rows = filtered.value
  let listen = 0, established = 0, udp = 0
  for (const c of rows) {
    if (c.state === 'LISTEN') listen++
    else if (c.state === 'ESTABLISHED') established++
    if (c.protocol === 'udp4' || c.protocol === 'udp6') udp++
  }
  return {
    total: conns.value.size,
    listen,
    established,
    udp,
  }
})
</script>

<template>
  <div class="app-root">
    <TitleBar @open-settings="settingsOpen = true" />
    <Toolbar :last-refreshed-at="lastRefreshedAt" />
    <FilterBar v-model="filter" />
    <div class="main">
      <ConnectionTable
        :rows="filtered"
        :highlights="highlights"
        :density="settings.density"
        @select="onSelect"
        @contextmenu="onContextmenu"
      />
      <DetailPanel :row="selected" @close="onSelect(null)" />
    </div>
    <StatsBar :stats="displayStats" :filtered="filteredCount" :running="settings.running" />
    <ContextMenu
      v-if="ctxMenu"
      :x="ctxMenu.x" :y="ctxMenu.y" :row="ctxMenu.row"
      @close="ctxMenu = null"
    />
    <SettingsDialog v-if="settingsOpen" @close="settingsOpen = false" />
    <ErrorDialog />
    <ConfirmDialog />
    <transition name="fade">
      <div v-if="errorText" class="error-toast">{{ errorText }}</div>
    </transition>
  </div>
</template>

<style scoped>
.app-root {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: var(--bg);
  color: var(--text);
  min-height: 0;
  position: relative;
}
.main {
  flex: 1;
  display: flex;
  min-height: 0;
  position: relative;
}
.error-toast {
  position: fixed;
  bottom: 40px; left: 50%;
  transform: translateX(-50%);
  background: var(--danger);
  color: #fff;
  padding: 8px 16px;
  border-radius: var(--radius-md);
  font-size: 12px;
  z-index: 2000;
  box-shadow: 0 4px 12px rgba(0,0,0,0.2);
}
.fade-enter-active, .fade-leave-active { transition: opacity var(--transition); }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
