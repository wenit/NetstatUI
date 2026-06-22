<script setup lang="ts">
import { Window } from '@wailsio/runtime'
import { useSettingsStore } from '../stores/settings'

const settings = useSettingsStore()

async function minimise() {
  await Window.Minimise()
}
async function toggleMax() {
  await Window.ToggleMaximise()
}
async function close() {
  await Window.Close()
}
</script>

<template>
  <div class="titlebar">
    <div class="drag" />
    <div class="brand">
      <span class="brand-icon" />
      <span class="brand-text">Network Ports</span>
    </div>
    <div class="spacer" />
    <button class="tb-btn theme" title="切换主题" @click="settings.toggleTheme()">
      <svg v-if="settings.theme === 'dark'" width="16" height="16" viewBox="0 0 16 16">
        <circle cx="8" cy="8" r="3.5" fill="currentColor" />
        <g stroke="currentColor" stroke-width="1.2" stroke-linecap="round">
          <line x1="8" y1="1.5" x2="8" y2="3" /><line x1="8" y1="13" x2="8" y2="14.5" />
          <line x1="1.5" y1="8" x2="3" y2="8" /><line x1="13" y1="8" x2="14.5" y2="8" />
          <line x1="3.2" y1="3.2" x2="4.3" y2="4.3" /><line x1="11.7" y1="11.7" x2="12.8" y2="12.8" />
          <line x1="3.2" y1="12.8" x2="4.3" y2="11.7" /><line x1="11.7" y1="4.3" x2="12.8" y2="3.2" />
        </g>
      </svg>
      <svg v-else width="16" height="16" viewBox="0 0 16 16">
        <path d="M13.5 8.8A5.5 5.5 0 1 1 7.2 2.5a4.3 4.3 0 0 0 6.3 6.3z" fill="currentColor" />
      </svg>
    </button>
    <button class="tb-btn" title="最小化" @click="minimise">
      <svg width="10" height="10" viewBox="0 0 10 10"><line x1="1" y1="5" x2="9" y2="5" stroke="currentColor" stroke-width="1" /></svg>
    </button>
    <button class="tb-btn" title="最大化" @click="toggleMax">
      <svg width="10" height="10" viewBox="0 0 10 10"><rect x="1.5" y="1.5" width="7" height="7" stroke="currentColor" stroke-width="1" fill="none" /></svg>
    </button>
    <button class="tb-btn close" title="关闭" @click="close">
      <svg width="10" height="10" viewBox="0 0 10 10">
        <line x1="1.5" y1="1.5" x2="8.5" y2="8.5" stroke="currentColor" stroke-width="1" />
        <line x1="8.5" y1="1.5" x2="1.5" y2="8.5" stroke="currentColor" stroke-width="1" />
      </svg>
    </button>
  </div>
</template>

<style scoped>
.titlebar {
  height: var(--titlebar-h);
  display: flex;
  align-items: center;
  padding: 0 4px 0 12px;
  background: transparent;
  position: relative;
  flex-shrink: 0;
}
.drag {
  position: absolute;
  inset: 0;
  --wails-draggable: drag;
}
.brand {
  display: flex;
  align-items: center;
  gap: 8px;
  z-index: 1;
  pointer-events: none;
}
.brand-icon {
  width: 14px;
  height: 14px;
  border-radius: 3px;
  background: linear-gradient(135deg, var(--accent), var(--state-established));
}
.brand-text {
  font-size: 12px;
  font-weight: 500;
  color: var(--text-secondary);
  letter-spacing: 0.02em;
}
.spacer { flex: 1; }
.tb-btn {
  height: 28px;
  width: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  color: var(--text-secondary);
  cursor: pointer;
  border-radius: var(--radius-sm);
  transition: background var(--transition-fast), color var(--transition-fast);
  z-index: 1;
}
.tb-btn:hover { background: var(--bg-hover); color: var(--text); }
.tb-btn.close:hover { background: var(--danger); color: #fff; }
.tb-btn.theme { width: 36px; }
</style>
