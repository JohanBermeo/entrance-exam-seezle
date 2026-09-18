import { useCallback, useEffect, useRef, useState } from 'react'
import { calculateExpression } from '../api/calculationApi'
import type { CalculationError, CalculationResponse } from '../api/types'
import { useHistory } from './useHistory'
import type { HistoryEntry } from './useHistory'

export interface UseCalculationResult {
  result: CalculationResponse | null
  error: CalculationError | null
  isLoading: boolean
  history: HistoryEntry[]
  calculate: (expression: string) => Promise<void>
  clearHistory: () => void
}

/** Estado + fetch + historial para el modo expression. Cancela la petición anterior. */
export function useCalculation(): UseCalculationResult {
  const [result, setResult] = useState<CalculationResponse | null>(null)
  const [error, setError] = useState<CalculationError | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const { entries: history, push, clear: clearHistory } = useHistory()
  const abortRef = useRef<AbortController | null>(null)
  const mountedRef = useRef(true)

  useEffect(() => {
    mountedRef.current = true
    return () => {
      mountedRef.current = false
      abortRef.current?.abort()
    }
  }, [])

  const calculate = useCallback(
    async (expression: string) => {
      abortRef.current?.abort()
      const controller = new AbortController()
      abortRef.current = controller
      setIsLoading(true)
      setError(null)
      try {
        const res = await calculateExpression(expression, { signal: controller.signal })
        if (!mountedRef.current || controller.signal.aborted) return
        setResult(res)
        const firstOutput = res.outputs[0] ?? 'result'
        push({
          expression,
          value: res.results[firstOutput]?.value ?? null,
          requestId: res.requestId,
        })
      } catch (e) {
        if (e instanceof DOMException && e.name === 'AbortError') return
        if (!mountedRef.current) return
        setError(e as CalculationError)
      } finally {
        if (mountedRef.current && abortRef.current === controller) setIsLoading(false)
      }
    },
    [push],
  )

  return { result, error, isLoading, history, calculate, clearHistory }
}
