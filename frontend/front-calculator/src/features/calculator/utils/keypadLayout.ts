import type { IconName } from '../../../components/index'

export type KeyType = 'number' | 'operation' | 'function' | 'equals' | 'empty'
export type KeyKind = 'input' | 'clear' | 'delete' | 'sign' | 'equals' | 'empty'

export interface KeyDef {
  id: string
  label: string
  kind: KeyKind
  type: KeyType
  /** Símbolo a añadir a la expresión (solo kind 'input'). */
  value?: string
  /** Columnas que ocupa en el grid. */
  span?: number
  /** Icono a renderizar en vez de la etiqueta. */
  icon?: IconName
}

function input(id: string, label: string, value: string, type: KeyType): KeyDef {
  return { id, label, kind: 'input', type, value }
}

function num(label: string): KeyDef {
  return { id: `n${label === '.' ? 'dot' : label}`, label, kind: 'input', type: 'number', value: label }
}

/** Layout definitivo 4 cols × 6 filas (docs/frontend-calculator-plan.md). */
export const KEYPAD_LAYOUT: KeyDef[][] = [
  [
    { id: 'ac', label: 'AC', kind: 'clear', type: 'function' },
    { id: 'sign', label: '+/-', kind: 'sign', type: 'function' },
    input('percent', '%', 'percent(', 'function'),
    input('divide', '÷', '/', 'operation'),
  ],
  [
    input('sqrt', '√', 'sqrt(', 'function'),
    input('power', 'xʸ', '^', 'function'),
    { id: 'delete', label: '', kind: 'delete', type: 'function', icon: 'delete' },
    input('multiply', '×', '*', 'operation'),
  ],
  [num('7'), num('8'), num('9'), input('minus', '−', '-', 'operation')],
  [num('4'), num('5'), num('6'), input('plus', '+', '+', 'operation')],
  [num('1'), num('2'), num('3'), { id: 'equals', label: '=', kind: 'equals', type: 'equals' }],
  [
    { ...num('0'), id: 'n0', span: 2 },
    num('.'),
    { id: 'empty', label: '', kind: 'empty', type: 'empty' },
  ],
]
