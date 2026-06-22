import { defineStore } from 'pinia'

export type Theme = 'dark' | 'light'

interface SettingsState {
  theme: Theme
  intervalMs: number
  running: boolean
  density: 'compact' | 'comfortable'
}

export const useSettingsStore = defineStore('settings', {
  state: (): SettingsState => ({
    theme: (localStorage.getItem('np.theme') as Theme) || 'dark',
    intervalMs: Number(localStorage.getItem('np.interval')) || 5000,
    running: true,
    density: 'comfortable',
  }),
  actions: {
    setTheme(t: Theme) {
      this.theme = t
      localStorage.setItem('np.theme', t)
      applyTheme(t)
    },
    toggleTheme() {
      this.setTheme(this.theme === 'dark' ? 'light' : 'dark')
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

export function applyTheme(t: Theme) {
  document.documentElement.setAttribute('data-theme', t)
}
