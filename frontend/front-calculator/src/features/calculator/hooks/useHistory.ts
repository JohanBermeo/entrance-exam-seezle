import { useCallback, useRef, useState } from 'react'

export interface HistoryEntry {
  id: string
  expression: string
  value: number | null
  requestId: string
  createdAt: number
}

export type NewHistoryEntry = Omit<HistoryEntry, 'id' | 'createdAt'>

export const HISTORY_LIMIT = 20

/** Historial en memoria (sin localStorage en v1). El más reciente primero. */
export function useHistory(limit: number = HISTORY_LIMIT) {
  const [entries, setEntries] = useState<HistoryEntry[]>([])
  const counterRef = useRef(0)

  const push = useCallback(
    (entry: NewHistoryEntry) => {
      counterRef.current += 1
      const full: HistoryEntry = {
        ...entry,
        id: `${Date.now()}-${counterRef.current}`,
        createdAt: Date.now(),
      }
      setEntries((prev) => [full, ...prev].slice(0, limit))
    },
    [limit],
  )

  const clear = useCallback(() => setEntries([]), [])

  return { entries, push, clear }
}
