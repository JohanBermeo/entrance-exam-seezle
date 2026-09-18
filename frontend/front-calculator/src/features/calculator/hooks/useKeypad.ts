import { useCallback, useEffect, useRef, useState } from 'react'
import { toggleLastNumberSign } from '../utils/expression'

const DIRECT_KEYS = new Set('0123456789+-*/().^%'.split(''))

export interface UseKeypadOptions {
  initialExpression?: string
  onEquals?: (expression: string) => void
}

/** Construye la expresión desde el keypad virtual y el teclado físico. */
export function useKeypad(options: UseKeypadOptions = {}) {
  const { initialExpression = '', onEquals } = options
  const [expression, setExpression] = useState(initialExpression)
  const onEqualsRef = useRef(onEquals)
  useEffect(() => {
    onEqualsRef.current = onEquals
  }, [onEquals])

  const input = useCallback((symbol: string) => {
    setExpression((prev) => prev + symbol)
  }, [])

  const clear = useCallback(() => setExpression(''), [])

  const backspace = useCallback(() => {
    setExpression((prev) => prev.slice(0, -1))
  }, [])

  const toggleSign = useCallback(() => {
    setExpression((prev) => toggleLastNumberSign(prev))
  }, [])

  const submit = useCallback(() => {
    setExpression((current) => {
      onEqualsRef.current?.(current)
      return current
    })
  }, [])

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

  return { expression, setExpression, input, clear, backspace, toggleSign, submit }
}
