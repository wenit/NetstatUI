import type { ConnRow } from './useConnections'
import type { State } from '../../bindings/github.com/zwb/network-ports/services/netstat/models'

export interface FilterState {
  search: string
  protocols: Set<string>
  states: Set<string>
  listenOnly: boolean
  externalOnly: boolean
}

export function emptyFilter(): FilterState {
  return {
    search: '',
    protocols: new Set(),
    states: new Set(),
    listenOnly: false,
    externalOnly: false,
  }
}

const LOCAL_ADDR = /^(127\.|10\.|192\.168\.|172\.(1[6-9]|2\d|3[01])\.|::1$|fe80::|::ffff:127\.|0\.0\.0\.0$|::$)/i

export function isLocal(addr: string): boolean {
  if (!addr || addr === '*') return true
  return LOCAL_ADDR.test(addr)
}

export function applyFilter(conns: Map<string, ConnRow>, f: FilterState): ConnRow[] {
  const q = f.search.trim().toLowerCase()
  const out: ConnRow[] = []
  for (const c of conns.values()) {
    if (f.protocols.size && !f.protocols.has(c.protocol)) continue
    if (f.states.size && !f.states.has(c.state)) continue
    if (f.listenOnly && c.state !== ('LISTEN' as State)) continue
    if (f.externalOnly && isLocal(c.remoteAddr)) continue
    if (q) {
      const hay = `${c.protocol} ${c.localAddr} ${c.localPort} ${c.remoteAddr} ${c.remotePort} ${c.state} ${c.pid} ${c.processName}`.toLowerCase()
      if (!hay.includes(q)) continue
    }
    out.push(c)
  }
  return out
}
