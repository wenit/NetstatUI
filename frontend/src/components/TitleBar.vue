<script setup lang="ts">
import { Window } from '@wailsio/runtime'
import { useSettingsStore } from '../stores/settings'
import { useI18n } from 'vue-i18n'

const settings = useSettingsStore()
const { t } = useI18n()

const emit = defineEmits<{ 'open-settings': [] }>()

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
    <button class="tb-btn settings" :title="t('titlebar.settings')" @click="emit('open-settings')">
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
        <path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z" />
        <circle cx="12" cy="12" r="3" />
      </svg>
    </button>
    <button class="tb-btn theme" :title="t('titlebar.toggleTheme')" @click="settings.toggleTheme()">
      <svg v-if="settings.resolvedTheme === 'dark'" width="16" height="16" viewBox="0 0 16 16">
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
    <button class="tb-btn" :title="t('titlebar.minimize')" @click="minimise">
      <svg width="10" height="10" viewBox="0 0 10 10"><line x1="1" y1="5" x2="9" y2="5" stroke="currentColor" stroke-width="1" /></svg>
    </button>
    <button class="tb-btn" :title="t('titlebar.maximize')" @click="toggleMax">
      <svg width="10" height="10" viewBox="0 0 10 10"><rect x="1.5" y="1.5" width="7" height="7" stroke="currentColor" stroke-width="1" fill="none" /></svg>
    </button>
    <button class="tb-btn close" :title="t('titlebar.close')" @click="close">
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
.tb-btn.settings { width: 36px; }
</style>
