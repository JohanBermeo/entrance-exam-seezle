const LAST_NUMBER = /(\d+(?:[.,]\d+)?)\s*$/

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
