import { shallowRef, readonly, type ShallowRef } from 'vue'
import { Events } from '@wailsio/runtime'
import type { ConnInfo } from '../../bindings/github.com/zwb/network-ports/services/netstat/models'
import { Stats } from '../../bindings/github.com/zwb/network-ports/services/monitor/models'
import type { Diff } from '../../bindings/github.com/zwb/network-ports/services/monitor/models'
import { Monitor } from '../../bindings/github.com/zwb/network-ports/services/monitor'
export type ConnRow = ConnInfo & {
  _addedAt: number
  _changedAt: number
  _removed?: boolean
}

export interface HighlightState {
  key: string
  type: 'add' | 'remove' | 'update'
  until: number
}

const conns = shallowRef<Map<string, ConnRow>>(new Map())
const stats = shallowRef<Stats>(new Stats())
const highlights = shallowRef<HighlightState[]>([])
const error = shallowRef<string>('')

let cleaned = false

function applyFull(list: ConnInfo[]) {
  const map = new Map<string, ConnRow>()
  const now = Date.now()
  const hl: HighlightState[] = []
  const prevKeys = new Set(conns.value.keys())
  for (const c of list) {
    const existed = prevKeys.has(c.key)
    map.set(c.key, { ...c, _addedAt: now, _changedAt: now })
    if (!existed) hl.push({ key: c.key, type: 'add', until: now + 3000 })
  }
  for (const k of prevKeys) {
    if (!map.has(k)) hl.push({ key: k, type: 'remove', until: now + 1500 })
  }
  conns.value = map
  if (hl.length) highlights.value = hl
}

function applyDiff(diff: Diff) {
  const now = Date.now()
  const map = new Map(conns.value)
  const hl: HighlightState[] = []

  for (const c of diff.added ?? []) {
    map.set(c.key, { ...c, _addedAt: now, _changedAt: now })
    hl.push({ key: c.key, type: 'add', until: now + 3000 })
  }
  for (const c of diff.updated ?? []) {
    const prev = map.get(c.key)
    map.set(c.key, { ...c, _addedAt: prev?._addedAt ?? now, _changedAt: now })
    hl.push({ key: c.key, type: 'update', until: now + 3000 })
  }
  for (const key of diff.removed ?? []) {
    if (map.delete(key)) hl.push({ key, type: 'remove', until: now + 1500 })
  }

  conns.value = map
  if (hl.length) highlights.value = hl
}

function cleanupHighlights() {
  if (cleaned) return
  const now = Date.now()
  const cur = highlights.value
  if (cur.length === 0) return
  const next = cur.filter((h) => h.until > now)
  if (next.length !== cur.length) highlights.value = next
}

export function useConnections() {
  return {
    conns: readonly(conns) as ShallowRef<Map<string, ConnRow>>,
    stats: readonly(stats) as ShallowRef<Stats>,
    highlights: readonly(highlights) as ShallowRef<HighlightState[]>,
    error: readonly(error) as ShallowRef<string>,
  }
}

export function initConnections() {
  Events.On('conn:full', (ev) => {
    applyFull(ev.data as unknown as ConnInfo[])
  })
  Events.On('conn:diff', (ev) => {
    const d = ev.data as any
    applyDiff({
      added: d.added ?? [],
      removed: d.removed ?? [],
      updated: d.updated ?? [],
    } as Diff)
  })
  Events.On('conn:stats', (ev) => {
    stats.value = Object.assign(new Stats(), ev.data)
  })
  Events.On('conn:error', (ev) => {
    error.value = ev.data as unknown as string
  })

  if (typeof window !== 'undefined') {
    window.setInterval(cleanupHighlights, 500)
  }
}

export async function startMonitor(intervalMs: number): Promise<void> {
  await Monitor.Start(intervalMs)
}

export async function stopMonitor(): Promise<void> {
  await Monitor.Stop()
}

export async function setIntervalMs(ms: number): Promise<void> {
  await Monitor.SetInterval(ms)
}

export async function refreshNow(): Promise<void> {
  await Monitor.RefreshNow()
}

export async function fetchSnapshot(): Promise<void> {
  const r = await Monitor.GetSnapshot()
  if (r && r.conns) {
    applyFull(r.conns)
    if (r.stats) stats.value = r.stats
  }
}
