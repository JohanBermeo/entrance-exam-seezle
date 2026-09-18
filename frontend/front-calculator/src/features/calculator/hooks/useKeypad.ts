import { useCallback, useEffect, useRef, useState } from 'react'
import { appendSymbol, toggleLastNumberSign } from '../utils/expression'

const DIRECT_KEYS = new Set('0123456789+-*/().^%'.split(''))

/** Tras `=`, estos símbolos empiezan una expresión nueva (el resto encadena). */
const FRESH_CLEAR = /^[0-9().]$/

export interface UseKeypadOptions {
  initialExpression?: string
  onEquals?: (expression: string) => void
  /** Se invoca en cada edición (virtual o física), no en submit. */
  onEdit?: () => void
}

/** Construye la expresión desde el keypad virtual y el teclado físico. */
export function useKeypad(options: UseKeypadOptions = {}) {
  const { initialExpression = '', onEquals, onEdit } = options
  const [expression, setExpressionState] = useState(initialExpression)
  /** Resultado recién calculado: número/paréntesis/punto lo borran, operación encadena. */
  const freshRef = useRef<string | null>(null)
  const onEqualsRef = useRef(onEquals)
  useEffect(() => {
    onEqualsRef.current = onEquals
  }, [onEquals])
  const onEditRef = useRef(onEdit)
  useEffect(() => {
    onEditRef.current = onEdit
  }, [onEdit])

  const setExpression = useCallback((value: string | ((prev: string) => string)) => {
    freshRef.current = null
    setExpressionState(value)
  }, [])

  /** Fija la expresión al resultado de `=` (habilita borrado fresco y encadenado). */
  const commitResult = useCallback((value: number) => {
    const text = String(value)
    freshRef.current = text
    setExpressionState(text)
  }, [])

  const input = useCallback((symbol: string) => {
    onEditRef.current?.()
    const fresh = freshRef.current
    freshRef.current = null
    setExpressionState((prev) =>
      fresh !== null && FRESH_CLEAR.test(symbol) && prev === fresh
        ? symbol
        : appendSymbol(prev, symbol),
    )
  }, [])

  const clear = useCallback(() => {
    onEditRef.current?.()
    freshRef.current = null
    setExpressionState('')
  }, [])

  const backspace = useCallback(() => {
    onEditRef.current?.()
    freshRef.current = null
    setExpressionState((prev) => prev.slice(0, -1))
  }, [])

  const toggleSign = useCallback(() => {
    onEditRef.current?.()
    freshRef.current = null
    setExpressionState((prev) => toggleLastNumberSign(prev))
  }, [])

  const submit = useCallback(() => {
    setExpression((current) => {
      onEqualsRef.current?.(current)
      return current
    })
  }, [setExpression])

  useEffect(() => {
    const onKeyDown = (e: KeyboardEvent) => {
      if (e.metaKey || e.ctrlKey || e.altKey) return
      if (e.key === 'Enter' || e.key === '=') {
        e.preventDefault()
        submit()
      } else if (e.key === 'Backspace') {
        e.preventDefault()
        backspace()
      } else if (e.key === 'Escape') {
        clear()
      } else if (DIRECT_KEYS.has(e.key)) {
        input(e.key)
      }
    }
    window.addEventListener('keydown', onKeyDown)
    return () => window.removeEventListener('keydown', onKeyDown)
  }, [input, backspace, clear, submit])

  return { expression, setExpression, input, clear, backspace, toggleSign, submit, commitResult }
}
