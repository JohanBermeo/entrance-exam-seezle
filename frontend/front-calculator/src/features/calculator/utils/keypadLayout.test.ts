import { describe, expect, it } from 'vitest'
import { KEYPAD_LAYOUT } from './keypadLayout'

describe('KEYPAD_LAYOUT', () => {
  it('es un grid de 6 filas que suman 4 columnas (el 0 ocupa 2)', () => {
    expect(KEYPAD_LAYOUT).toHaveLength(6)
    for (const row of KEYPAD_LAYOUT) {
      const span = row.reduce((acc, k) => acc + (k.span ?? 1), 0)
      expect(span).toBe(4)
    }
  })

  it('tiene ids únicos', () => {
    const ids = KEYPAD_LAYOUT.flat().map((k) => k.id)
    expect(new Set(ids).size).toBe(ids.length)
  })

  it('el 0 ocupa 2 columnas y la última celda está vacía', () => {
    const last = KEYPAD_LAYOUT[5]!
    expect(last[0]).toMatchObject({ id: 'n0', span: 2 })
    expect(last[2]).toMatchObject({ kind: 'empty' })
  })

  it('el = aparece solo en la fila 5', () => {
    const equals = KEYPAD_LAYOUT.flat().filter((k) => k.kind === 'equals')
    expect(equals).toHaveLength(1)
    expect(KEYPAD_LAYOUT[4]).toContain(equals[0])
  })

  it('las teclas de entrada llevan su valor de expresión', () => {
    const byId = Object.fromEntries(KEYPAD_LAYOUT.flat().map((k) => [k.id, k]))
    expect(byId.divide?.value).toBe('/')
    expect(byId.multiply?.value).toBe('*')
    expect(byId.sqrt?.value).toBe('sqrt(')
    expect(byId.power?.value).toBe('^')
    expect(byId.percent?.value).toBe('percent(')
  })
})
