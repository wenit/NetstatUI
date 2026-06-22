<script setup lang="ts">
import { reactive, computed } from 'vue'
import type { FilterState } from '../composables/useFilters'

const props = defineProps<{ modelValue: FilterState }>()
const emit = defineEmits<{ 'update:modelValue': [FilterState] }>()

const protoOptions = ['tcp4', 'tcp6', 'udp4', 'udp6']
const stateOptions = ['LISTEN', 'ESTABLISHED', 'TIME_WAIT', 'CLOSE_WAIT', 'FIN_WAIT_2', 'SYN_SENT', 'CLOSING']

const local = reactive({
  search: props.modelValue.search,
  proto: new Set(props.modelValue.protocols),
  states: new Set(props.modelValue.states),
  listenOnly: props.modelValue.listenOnly,
  externalOnly: props.modelValue.externalOnly,
})

const activeChips = computed(() => {
  const n = local.proto.size + local.states.size + (local.listenOnly ? 1 : 0) + (local.externalOnly ? 1 : 0)
  return n
})

function toggle(set: Set<string>, v: string) {
  if (set.has(v)) set.delete(v)
  else set.add(v)
  emitUpdate()
}
function onSearch() { emitUpdate() }
function toggleListen() { local.listenOnly = !local.listenOnly; emitUpdate() }
function toggleExternal() { local.externalOnly = !local.externalOnly; emitUpdate() }
function clearAll() {
  local.search = ''
  local.proto.clear()
  local.states.clear()
  local.listenOnly = false
  local.externalOnly = false
  emitUpdate()
}
function emitUpdate() {
  emit('update:modelValue', {
    search: local.search,
    protocols: new Set(local.proto),
    states: new Set(local.states),
    listenOnly: local.listenOnly,
    externalOnly: local.externalOnly,
  })
}
</script>

<template>
  <div class="filterbar">
    <div class="search">
      <svg width="14" height="14" viewBox="0 0 14 14" class="search-icon">
        <circle cx="6" cy="6" r="4" stroke="currentColor" stroke-width="1.2" fill="none" />
        <line x1="9.5" y1="9.5" x2="12.5" y2="12.5" stroke="currentColor" stroke-width="1.2" stroke-linecap="round" />
      </svg>
      <input
        v-model="local.search"
        type="text"
        placeholder="搜索 地址 / 端口 / 进程 / PID…"
        spellcheck="false"
        @input="onSearch"
      />
      <button v-if="local.search" class="clear" @click="local.search = ''; onSearch()">×</button>
    </div>

    <div class="chips">
      <button
        v-for="p in protoOptions" :key="p"
        class="chip" :class="{ on: local.proto.has(p) }"
        @click="toggle(local.proto, p)"
      >{{ p }}</button>
      <span class="sep" />
      <button
        v-for="s in stateOptions" :key="s"
        class="chip state" :class="{ on: local.states.has(s) }"
        @click="toggle(local.states, s)"
      >{{ s }}</button>
      <span class="sep" />
      <button class="chip" :class="{ on: local.listenOnly }" @click="toggleListen">仅监听</button>
      <button class="chip" :class="{ on: local.externalOnly }" @click="toggleExternal">仅外部</button>
      <button v-if="activeChips || local.search" class="chip clear-btn" @click="clearAll">清除</button>
    </div>
  </div>
</template>

<style scoped>
.filterbar {
  height: var(--filterbar-h);
  display: flex;
  align-items: center;
  padding: 0 12px;
  gap: 12px;
  border-bottom: 1px solid var(--border);
  flex-shrink: 0;
}
.search {
  display: flex;
  align-items: center;
  width: 280px;
  height: 28px;
  background: var(--bg-active);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  padding: 0 8px;
  gap: 6px;
  transition: border-color var(--transition-fast);
}
.search:focus-within { border-color: var(--accent); }
.search-icon { color: var(--text-tertiary); flex-shrink: 0; }
.search input {
  flex: 1;
  border: none;
  outline: none;
  background: transparent;
  color: var(--text);
  font-size: 12px;
  font-family: var(--font-mono);
}
.search input::placeholder { color: var(--text-tertiary); }
.clear {
  border: none; background: transparent; color: var(--text-tertiary);
  cursor: pointer; font-size: 16px; line-height: 1;
}
.clear:hover { color: var(--text); }
.chips { display: flex; align-items: center; gap: 4px; flex-wrap: wrap; }
.chip {
  height: 22px;
  padding: 0 8px;
  border: 1px solid var(--border);
  background: transparent;
  color: var(--text-secondary);
  font-size: 11px;
  border-radius: 12px;
  cursor: pointer;
  transition: all var(--transition-fast);
  font-family: var(--font-mono);
}
.chip:hover { border-color: var(--border-strong); color: var(--text); }
.chip.on { background: var(--accent); border-color: var(--accent); color: #fff; }
.chip.state.on { background: var(--state-established); border-color: var(--state-established); }
.clear-btn { color: var(--danger); border-color: transparent; }
.clear-btn:hover { background: var(--row-remove); }
.sep { width: 1px; height: 16px; background: var(--border); margin: 0 4px; }
</style>
