const LAST_NUMBER = /(\d+(?:[.,]\d+)?)\s*$/
const CALL_NAME_END = /[a-zA-Z_][a-zA-Z0-9_]*$/

/** Símbolos que abren función o grupo y exigen `×` explícito tras dígito o `)`. */
const IMPLICIT_MULT_OPENERS = new Set(['sqrt(', 'percent(', '('])

function isBalanced(text: string): boolean {
  let depth = 0
  for (const ch of text) {
    if (ch === '(') depth++
    else if (ch === ')') {
      depth--
      if (depth < 0) return false
    }
  }
  return depth === 0
}

/** Grupo balanceado final (`prefix` + `group`) o null si no termina en grupo. */
function trailingGroup(s: string): { prefix: string; group: string } | null {
  if (!s.endsWith(')')) return null
  let depth = 0
  for (let i = s.length - 1; i >= 0; i--) {
    if (s[i] === ')') depth++
    else if (s[i] === '(') {
      depth--
      if (depth === 0) return { prefix: s.slice(0, i), group: s.slice(i) }
    }
  }
  return null
}

/** Un único operando: número, grupo total o llamada `nombre(...)`. */
function isSingleOperand(text: string): boolean {
  if (/^\d+(?:[.,]\d+)?$/.test(text)) return true
  const group = trailingGroup(text)
  if (!group) return false
  if (group.prefix === '') return true
  return /^[a-zA-Z_][a-zA-Z0-9_]*$/.test(group.prefix)
}

/** ¿El `-` en `index` es unario (no una resta binaria)? */
function isUnaryMinusAt(s: string, index: number): boolean {
  let i = index - 1
  while (i >= 0 && s[i] === ' ') i--
  if (i < 0) return true
  const prev = s[i]
  return prev !== ')' && !/\d/.test(prev ?? '')
}

/**
 * Niega el último operando con paréntesis para lógica matemática:
 * `5+3` → `5+(-3)`, `5+(2+3)` → `5+(-(2+3))`, `sqrt(9)` → `-sqrt(9)`.
 * Aplicado de nuevo, quita la negación.
 */
export function toggleLastNumberSign(expression: string): string {
  const s = expression.trimEnd()

  // 1. Quitar negación total sobre un único operando: "-sqrt(9)" → "sqrt(9)".
  const whole = /^-(.+)$/.exec(s)
  if (whole && whole[1] !== undefined && isSingleOperand(whole[1])) {
    return whole[1]
  }

  // 2. Operando final: grupo balanceado o número.
  let prefix: string
  let operand: string
  const group = trailingGroup(s)
  if (group) {
    prefix = group.prefix
    operand = group.group
  } else {
    const numMatch = LAST_NUMBER.exec(s)
    if (!numMatch || numMatch.index === undefined || numMatch[1] === undefined) {
      return s.startsWith('-') ? s.slice(1) : `-${s}`
    }
    prefix = s.slice(0, numMatch.index)
    operand = numMatch[1]
  }

  // 3. Quitar envoltura "(-X)": "5+(-3)" → "5+3".
  const wrapped = /^\(-(.*)\)$/.exec(operand)
  if (wrapped && wrapped[1] !== undefined && isBalanced(wrapped[1])) {
    return prefix + wrapped[1]
  }

  // 4. Operando = toda la expresión: número → "(-N)", grupo → "-(...)".
  if (prefix === '') {
    return operand.startsWith('(') ? `-${operand}` : `(-${operand})`
  }

  // 5. Llamada a función: negar la llamada completa: "sqrt(9)" → "-sqrt(9)".
  if (CALL_NAME_END.test(prefix.trimEnd())) {
    return `-${s}`
  }

  // 6. Quitar '-' unario previo: "5+-3" → "5+3".
  const minus = /-\s*$/.exec(prefix)
  if (minus && minus.index !== undefined && isUnaryMinusAt(prefix, minus.index)) {
    return prefix.slice(0, minus.index) + operand
  }

  // 7. Envolver el operando: "5+3" → "5+(-3)".
  return `${prefix}(-${operand})`
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
