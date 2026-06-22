import { createI18n } from 'vue-i18n'
import zhCN from './zh-CN'
import enUS from './en-US'

export type AppLocale = 'zh-CN' | 'en-US'

export const SUPPORTED_LOCALES: AppLocale[] = ['zh-CN', 'en-US']

export function matchLocale(raw: string | undefined | null): AppLocale {
  if (!raw) return 'en-US'
  const lower = raw.toLowerCase()
  if (lower.startsWith('zh')) return 'zh-CN'
  return 'en-US'
}

export function detectInitialLocale(): AppLocale {
  const stored = localStorage.getItem('np.locale')
  if (stored) {
    const matched = matchLocale(stored)
    if (matched) return matched
  }
  if (typeof navigator !== 'undefined' && navigator.language) {
    return matchLocale(navigator.language)
  }
  return 'en-US'
}

const i18n = createI18n({
  legacy: false,
  locale: detectInitialLocale(),
  fallbackLocale: 'en-US',
  messages: {
    'zh-CN': zhCN,
    'en-US': enUS,
  },
})

export default i18n

export function setI18nLocale(locale: AppLocale) {
  i18n.global.locale.value = locale
}
