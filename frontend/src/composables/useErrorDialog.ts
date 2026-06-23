import { ref, type Ref } from 'vue'

export interface ErrorPayload {
  title: string
  body: string
}

const current: Ref<ErrorPayload | null> = ref(null)

export function showError(title: string, body: string) {
  current.value = { title, body }
}

export function dismissError() {
  current.value = null
}

export function useErrorDialog() {
  return {
    current,
    dismiss: dismissError,
  }
}