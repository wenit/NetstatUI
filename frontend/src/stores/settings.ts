import { defineStore } from 'pinia'
import { setI18nLocale, matchLocale, type AppLocale } from '../locales'

export type Theme = 'auto' | 'dark' | 'light'
export type ResolvedTheme = 'dark' | 'light'

interface SettingsState {
  theme: Theme
  locale: AppLocale
  intervalMs: number
  running: boolean
  density: 'compact' | 'comfortable'
}

export function resolveTheme(t: Theme): ResolvedTheme {
  if (t === 'auto') {
    if (typeof window !== 'undefined' && window.matchMedia) {
      return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
    }
    return 'light'
  }
  return t
}

export function applyTheme(t: Theme) {
  document.documentElement.setAttribute('data-theme', resolveTheme(t))
}

export const useSettingsStore = defineStore('settings', {
  state: (): SettingsState => ({
    theme: (localStorage.getItem('np.theme') as Theme) || 'auto',
    locale: matchLocale(localStorage.getItem('np.locale')) || matchLocale(navigator.language),
    intervalMs: Number(localStorage.getItem('np.interval')) || 5000,
    running: true,
    density: 'comfortable',
  }),
  getters: {
    resolvedTheme(state): ResolvedTheme {
      return resolveTheme(state.theme)
    },
  },
  actions: {
    setTheme(t: Theme) {
      this.theme = t
      localStorage.setItem('np.theme', t)
      applyTheme(t)
    },
    toggleTheme() {
      this.setTheme(this.resolvedTheme === 'dark' ? 'light' : 'dark')
    },
    setLocale(l: AppLocale) {
      this.locale = l
      localStorage.setItem('np.locale', l)
      setI18nLocale(l)
    },
    setInterval(ms: number) {
      this.intervalMs = ms
      localStorage.setItem('np.interval', String(ms))
    },
    setRunning(r: boolean) {
      this.running = r
    },
  },
})
