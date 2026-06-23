import { shallowRef, readonly } from 'vue'

export interface ErrorPayload {
  title: string
  body: string
}

const current = shallowRef<ErrorPayload | null>(null)
let seq = 0

export function showError(title: string, body: string) {
  current.value = { title, body }
}

export function dismissError() {
  current.value = null
}

export function useErrorDialog() {
  return {
    current: readonly(current),
    dismiss: dismissError,
  }
}