const LAST_NUMBER = /(\d+(?:[.,]\d+)?)\s*$/

/** Símbolos que abren función o grupo y exigen `×` explícito tras dígito o `)`. */
const IMPLICIT_MULT_OPENERS = new Set(['sqrt(', 'percent(', '('])

/**
 * Niega el último número de la expresión (tecla +/-).
 * Sin número al final, niega toda la expresión.
 */
export function toggleLastNumberSign(expression: string): string {
  const match = LAST_NUMBER.exec(expression)
  if (!match || match.index === undefined) {
    return expression.startsWith('-') ? expression.slice(1) : `-${expression}`
  }
  const numStart = match.index
  let i = numStart - 1
  while (i >= 0 && expression[i] === ' ') i--
  if (expression[i] === '-' && (i === 0 || '+-*/^('.includes(expression[i - 1] ?? ''))) {
    return expression.slice(0, i) + expression.slice(i + 1)
  }
  return `${expression.slice(0, numStart)}-${match[1]}${expression.slice(numStart + match[0].length)}`
}

/**
 * Añade un símbolo a la expresión. La gramática del backend no acepta
 * multiplicación implícita (`9sqrt(4)` → parse error), así que ante un
 * abridor (`sqrt(`, `percent(`, `(`) tras dígito o `)` se inserta `*`.
 */
export function appendSymbol(expression: string, symbol: string): string {
  if (IMPLICIT_MULT_OPENERS.has(symbol) && /[\d)]$/.test(expression)) {
    return `${expression}*${symbol}`
  }
  return expression + symbol
}
