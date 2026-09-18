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
  /** Calcula y devuelve el valor de la primera salida (`null` si falla o se aborta). */
  calculate: (expression: string) => Promise<number | null>
  clearHistory: () => void
  clearError: () => void
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
        if (!mountedRef.current || controller.signal.aborted) return null
        setResult(res)
        const firstOutput = res.outputs[0] ?? 'result'
        const value = res.results[firstOutput]?.value ?? null
        push({ expression, value, requestId: res.requestId })
        return value
      } catch (e) {
        if (e instanceof DOMException && e.name === 'AbortError') return null
        if (!mountedRef.current) return null
        setError(e as CalculationError)
        return null
      } finally {
        if (mountedRef.current && abortRef.current === controller) setIsLoading(false)
      }
    },
    [push],
  )

  const clearError = useCallback(() => setError(null), [])

  return { result, error, isLoading, history, calculate, clearHistory, clearError }
}
