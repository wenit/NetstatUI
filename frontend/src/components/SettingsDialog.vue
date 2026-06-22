<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useSettingsStore } from '../stores/settings'
import { useI18n } from 'vue-i18n'
import type { AppLocale } from '../locales'
import type { Theme } from '../stores/settings'

const settings = useSettingsStore()
const { t } = useI18n()

const emit = defineEmits<{ close: [] }>()

const locales: { value: AppLocale; labelKey: string }[] = [
  { value: 'zh-CN', labelKey: 'settings.langZh' },
  { value: 'en-US', labelKey: 'settings.langEn' },
]

const themes: { value: Theme; labelKey: string }[] = [
  { value: 'auto', labelKey: 'settings.themeAuto' },
  { value: 'light', labelKey: 'settings.themeLight' },
  { value: 'dark', labelKey: 'settings.themeDark' },
]

function onKey(e: KeyboardEvent) {
  if (e.key === 'Escape') emit('close')
}

onMounted(() => document.addEventListener('keydown', onKey))
onUnmounted(() => document.removeEventListener('keydown', onKey))
</script>

<template>
  <div class="overlay" @click="emit('close')">
    <div class="dialog" @click.stop>
      <div class="header">
        <span class="title">{{ t('settings.title') }}</span>
        <button class="close-btn" @click="emit('close')">×</button>
      </div>
      <div class="body">
        <div class="section">
          <div class="sec-title">{{ t('settings.language') }}</div>
          <div class="options">
            <button
              v-for="l in locales"
              :key="l.value"
              class="opt"
              :class="{ active: settings.locale === l.value }"
              @click="settings.setLocale(l.value)"
            >{{ t(l.labelKey) }}</button>
          </div>
        </div>
        <div class="section">
          <div class="sec-title">{{ t('settings.theme') }}</div>
          <div class="options">
            <button
              v-for="th in themes"
              :key="th.value"
              class="opt"
              :class="{ active: settings.theme === th.value }"
              @click="settings.setTheme(th.value)"
            >{{ t(th.labelKey) }}</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 3000;
  backdrop-filter: blur(2px);
}
.dialog {
  width: 320px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-strong);
  border-radius: var(--radius-md);
  box-shadow: 0 12px 36px rgba(0, 0, 0, 0.3);
  overflow: hidden;
}
.header {
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 14px;
  border-bottom: 1px solid var(--border);
}
.title {
  font-size: 13px;
  font-weight: 600;
  color: var(--text);
}
.close-btn {
  border: none;
  background: transparent;
  color: var(--text-tertiary);
  font-size: 18px;
  cursor: pointer;
  width: 24px;
  height: 24px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.close-btn:hover {
  background: var(--bg-hover);
  color: var(--text);
}
.body {
  padding: 16px 14px;
}
.section {
  margin-bottom: 16px;
}
.section:last-child {
  margin-bottom: 0;
}
.sec-title {
  font-size: 11px;
  font-weight: 600;
  color: var(--text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 8px;
}
.options {
  display: flex;
  gap: 6px;
}
.opt {
  flex: 1;
  height: 30px;
  border: 1px solid var(--border);
  background: var(--bg-active);
  color: var(--text-secondary);
  font-size: 12px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all var(--transition-fast);
}
.opt:hover {
  border-color: var(--border-strong);
  color: var(--text);
}
.opt.active {
  background: var(--accent);
  border-color: var(--accent);
  color: #fff;
}
</style>
