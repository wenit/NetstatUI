import { ref, type Ref } from 'vue'

export interface ConfirmPayload {
  title: string
  body: string
  resolve: (ok: boolean) => void
}

const current: Ref<ConfirmPayload | null> = ref(null)

export function showConfirm(title: string, body: string): Promise<boolean> {
  return new Promise<boolean>((resolve) => {
    current.value = { title, body, resolve }
  })
}

export function resolveConfirm(ok: boolean) {
  const c = current.value
  if (!c) return
  current.value = null
  c.resolve(ok)
}

export function useConfirmDialog() {
  return {
    current,
    confirm: () => resolveConfirm(true),
    cancel: () => resolveConfirm(false),
  }
}