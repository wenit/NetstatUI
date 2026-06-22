<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useSettingsStore } from '../stores/settings'
import { useI18n } from 'vue-i18n'
import type { AppLocale } from '../locales'
import type { Theme } from '../stores/settings'

const settings = useSettingsStore()
const { t } = useI18n()

const emit = defineEmits<{ close: [] }>()

type Tab = 'general' | 'advanced'
const activeTab = ref<Tab>('general')

const locales: { value: AppLocale; labelKey: string }[] = [
  { value: 'zh-CN', labelKey: 'settings.langZh' },
  { value: 'en-US', labelKey: 'settings.langEn' },
]

const themes: { value: Theme; labelKey: string }[] = [
  { value: 'auto', labelKey: 'settings.themeAuto' },
  { value: 'light', labelKey: 'settings.themeLight' },
  { value: 'dark', labelKey: 'settings.themeDark' },
]

const cleared = ref(false)

const cacheKeys = ['np.locale', 'np.theme', 'np.interval', 'np.running']

function cacheValue(key: string): string {
  return localStorage.getItem(key) ?? '(not set)'
}

function clearData() {
  localStorage.clear()
  cleared.value = true
  setTimeout(() => location.reload(), 800)
}

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
      <div class="layout">
        <div class="tabs">
          <button
            v-for="tab in (['general', 'advanced'] as Tab[])"
            :key="tab"
            class="tab"
            :class="{ active: activeTab === tab }"
            @click="activeTab = tab"
          >{{ t(`settings.${tab}`) }}</button>
        </div>
        <div class="body">
          <div v-if="activeTab === 'general'">
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
          <div v-if="activeTab === 'advanced'" class="adv-body">
            <div class="section">
              <div class="sec-title">{{ t('settings.clearData') }}</div>
              <p class="desc">{{ t('settings.clearDataDesc') }}</p>
              <p v-if="cleared" class="cleared">{{ t('settings.cleared') }}</p>
              <button v-else class="danger-btn" @click="clearData">{{ t('settings.clearBtn') }}</button>
            </div>
            <div class="section">
              <div class="sec-title">localStorage</div>
              <table class="cache-table">
                <tr v-for="k in cacheKeys" :key="k">
                  <td class="ck">{{ k }}</td>
                  <td class="cv">{{ cacheValue(k) }}</td>
                </tr>
              </table>
            </div>
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
  width: 520px;
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
.layout {
  display: flex;
  min-height: 300px;
}
.tabs {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: 12px 0;
  border-right: 1px solid var(--border);
  width: 100px;
  flex-shrink: 0;
}
.tab {
  height: 32px;
  padding: 0 14px;
  border: none;
  background: transparent;
  color: var(--text-secondary);
  font-size: 12px;
  text-align: left;
  cursor: pointer;
  transition: all var(--transition-fast);
  border-radius: 0;
}
.tab:hover {
  color: var(--text);
  background: var(--bg-hover);
}
.tab.active {
  color: var(--accent);
  background: var(--bg-active);
  font-weight: 600;
}
.body {
  flex: 1;
  padding: 16px 14px;
}
.adv-body {
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-height: 120px;
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
.desc {
  font-size: 11px;
  color: var(--text-tertiary);
  margin: 0 0 10px;
  line-height: 1.5;
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
.danger-btn {
  height: 30px;
  padding: 0 14px;
  border: 1px solid var(--danger);
  background: transparent;
  color: var(--danger);
  font-size: 12px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all var(--transition-fast);
}
.danger-btn:hover {
  background: var(--danger);
  color: #fff;
}
.cleared {
  font-size: 12px;
  color: var(--state-established);
  margin: 0;
}
.cache-table {
  width: 100%;
  font-size: 11px;
  font-family: var(--font-mono);
  border-collapse: collapse;
}
.cache-table td {
  padding: 4px 0;
  border-bottom: 1px solid var(--border);
}
.ck {
  color: var(--text-tertiary);
  width: 40%;
}
.cv {
  color: var(--text);
  text-align: right;
}
</style>
