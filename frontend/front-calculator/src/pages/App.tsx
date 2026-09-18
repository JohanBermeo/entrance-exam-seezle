import { useCallback, useEffect, useRef, useState } from 'react'
import {
  CalculatorLayout,
  Display,
  ErrorToast,
  HistoryPanel,
  Keypad,
  useCalculation,
  useExpressionValidation,
  useKeypad,
} from '../features/calculator/index'
import type {
  CalculationError,
  HistoryEntry,
  KeyDef,
} from '../features/calculator/index'
import './App.css'

/** F4 · Integración: cablea keypad + validación + API + historial. */
function App() {
  const {
    result,
    error: serverError,
    isLoading,
    history,
    calculate,
    clearHistory,
    clearError,
  } = useCalculation()
  const { validate } = useExpressionValidation()
  const [clientError, setClientError] = useState<CalculationError | null>(null)
  const [copied, setCopied] = useState(false)

  // El error de validación y el flag de copiado expiran al editar o recalcular.
  const markEdited = useCallback(() => {
    setClientError(null)
    setCopied(false)
  }, [])

  // Ref puente: useKeypad necesita onEquals antes de que exista handleEquals,
  // que a su vez necesita setExpression de useKeypad.
  const handleEqualsRef = useRef<(expression: string) => Promise<void>>(async () => {})

  const {
    expression,
    setExpression,
    input,
    clear,
    backspace,
    toggleSign,
    submit,
  } = useKeypad({
    onEquals: (expr) => void handleEqualsRef.current(expr),
    onEdit: markEdited,
  })

  const handleEquals = useCallback(
    async (expression: string) => {
      markEdited()
      const issue = validate(expression)
      if (!issue.valid) {
        setClientError({
          type: 'validation',
          message: issue.error ?? 'Expresión inválida',
          position: issue.position,
        })
        return
      }
      setClientError(null)
      const value = await calculate(expression)
      // Tras `=`, la expresión pasa a ser el resultado (encadenar operaciones).
      if (value !== null && value !== undefined) setExpression(String(value))
    },
    [validate, calculate, markEdited, setExpression],
  )

  useEffect(() => {
    handleEqualsRef.current = handleEquals
  }, [handleEquals])

  const handleKeyPress = useCallback(
    (def: KeyDef) => {
      switch (def.kind) {
        case 'input':
          if (def.value !== undefined) input(def.value)
          break
        case 'clear':
          clear()
          clearError()
          break
        case 'delete':
          backspace()
          break
        case 'sign':
          toggleSign()
          break
        case 'equals':
          submit()
          break
        case 'empty':
          break
      }
    },
    [input, clear, clearError, backspace, toggleSign, submit],
  )

  const handleSelectHistory = useCallback(
    (entry: HistoryEntry) => {
      markEdited()
      setExpression(entry.expression)
    },
    [setExpression, markEdited],
  )

  const firstOutput = result?.outputs[0] ?? 'result'
  const value = result ? (result.results[firstOutput]?.value ?? null) : null
  const error = clientError ?? serverError

  const dismissError = useCallback(() => {
    setClientError(null)
    clearError()
  }, [clearError])

  const handleCopy = useCallback(async () => {
    if (value === null || value === undefined) return
    try {
      await navigator.clipboard.writeText(String(value))
      setCopied(true)
    } catch {
      setCopied(false)
    }
  }, [value])

  return (
    <CalculatorLayout>
      <h1 className="sr-only">Calculadora</h1>
      <div className="calc-stack">
        <Display expression={expression} result={value} isLoading={isLoading} />
        <div className="calc-actions">
          {result ? (
            <span className="calc-meta">
              {result.requestId} · {result.durationMs} ms
            </span>
          ) : (
            <span className="calc-meta" aria-hidden="true" />
          )}
          <button
            type="button"
            className="copy-btn"
            onClick={() => void handleCopy()}
            disabled={value === null || value === undefined}
            title="Copiar resultado"
          >
            {copied ? '¡Copiado!' : 'Copiar'}
          </button>
        </div>
        <ErrorToast error={error} onDismiss={dismissError} />
        <Keypad onKeyPress={handleKeyPress} disabled={isLoading} />
        <HistoryPanel entries={history} onSelect={handleSelectHistory} onClear={clearHistory} />
      </div>
    </CalculatorLayout>
  )
}

export default App
