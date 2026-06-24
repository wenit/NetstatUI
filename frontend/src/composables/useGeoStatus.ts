import { ref, readonly, type Ref } from 'vue'
import { Events } from '@wailsio/runtime'
import { Monitor } from '../../bindings/github.com/wenit/NetstatUI/services/monitor'
import { GeoStatus } from '../../bindings/github.com/wenit/NetstatUI/services/monitor/models'

export type GeoState = 'loading' | 'ready' | 'error' | 'disabled' | 'unknown'

const status = ref<GeoStatus>(new GeoStatus({ state: 'unknown' }))
let inited = false

function apply(s: GeoStatus) {
  status.value = s
}

export function useGeoStatus() {
  return {
    status: readonly(status) as Ref<GeoStatus>,
  }
}

export function initGeoStatus() {
  if (inited) return
  inited = true

  Events.On('geo:status', (ev) => {
    const data = ev.data as any
    apply(new GeoStatus({ state: data?.state ?? 'unknown', error: data?.error }))
  })

  // Pull the current state on startup so the loading indicator is
  // correct even if we missed the initial event (e.g. resolved before
  // the listener was registered).
  Monitor.GetGeoStatus()
    .then((s) => apply(s))
    .catch(() => apply(new GeoStatus({ state: 'error', error: 'GetGeoStatus failed' })))
}
