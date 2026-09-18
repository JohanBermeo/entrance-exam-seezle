import { useCallback } from 'react'

export interface ExpressionIssue {
  valid: boolean
  error?: string
  position?: number
}

const VALID_CHARS = /^[\d\s+\-*/^().,a-zA-Z]+$/
const KNOWN_FUNCTIONS = ['sqrt', 'percent']
export const MAX_EXPRESSION_LENGTH = 200

function firstInvalidCharPosition(expression: string): number {
  const match = /[^\d\s+\-*/^().,a-zA-Z]/.exec(expression)
  return match?.index ?? 0
}

function parenError(expression: string): ExpressionIssue | null {
  let balance = 0
  for (let i = 0; i < expression.length; i++) {
    if (expression[i] === '(') balance++
    if (expression[i] === ')') {
      balance--
      if (balance < 0) return { valid: false, error: 'Paréntesis de cierre sin apertura', position: i }
    }
  }
  if (balance !== 0) return { valid: false, error: 'Paréntesis sin cerrar' }
  return null
}

function unknownFunction(expression: string): ExpressionIssue | null {
  const identifiers = expression.matchAll(/[a-zA-Z_][a-zA-Z0-9_]*/g)
  for (const match of identifiers) {
    if (!KNOWN_FUNCTIONS.includes(match[0])) {
      return { valid: false, error: `Función desconocida: ${match[0]}`, position: match.index }
    }
  }
  return null
}

/** Validación cliente liviana. El backend es la autoridad (400/422). */
export function useExpressionValidation() {
  const validate = useCallback((expression: string): ExpressionIssue => {
    if (!expression.trim()) return { valid: false, error: 'Expresión vacía' }
    if (expression.length > MAX_EXPRESSION_LENGTH) {
      return { valid: false, error: `Expresión demasiado larga (máx. ${MAX_EXPRESSION_LENGTH})` }
    }
    if (!VALID_CHARS.test(expression)) {
      return {
        valid: false,
        error: 'Caracteres no permitidos',
        position: firstInvalidCharPosition(expression),
      }
    }
    return parenError(expression) ?? unknownFunction(expression) ?? { valid: true }
  }, [])

  return { validate }
}
