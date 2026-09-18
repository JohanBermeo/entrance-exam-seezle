import { describe, expect, it } from 'vitest'
import { KEYPAD_LAYOUT } from './keypadLayout'

describe('KEYPAD_LAYOUT', () => {
  it('es un grid de 6 filas × 4 columnas sin spans ni vacíos', () => {
    expect(KEYPAD_LAYOUT).toHaveLength(6)
    for (const row of KEYPAD_LAYOUT) {
      expect(row).toHaveLength(4)
      expect(row.every((k) => (k.span ?? 1) === 1 && k.kind !== 'empty')).toBe(true)
    }
  })

  it('tiene ids únicos', () => {
    const ids = KEYPAD_LAYOUT.flat().map((k) => k.id)
    expect(new Set(ids).size).toBe(ids.length)
  })

  it('sigue el orden acordado por filas', () => {
    const ids = KEYPAD_LAYOUT.map((row) => row.map((k) => k.id))
    expect(ids[0]).toEqual(['ac', 'sign', 'open-paren', 'close-paren'])
    expect(ids[1]).toEqual(['sqrt', 'power', 'delete', 'percent'])
    expect(ids[2]).toEqual(['n7', 'n8', 'n9', 'divide'])
    expect(ids[3]).toEqual(['n4', 'n5', 'n6', 'multiply'])
    expect(ids[4]).toEqual(['n1', 'n2', 'n3', 'minus'])
    expect(ids[5]).toEqual(['n0', 'ndot', 'equals', 'plus'])
  })

  it('el = aparece solo en la fila 6', () => {
    const equals = KEYPAD_LAYOUT.flat().filter((k) => k.kind === 'equals')
    expect(equals).toHaveLength(1)
    expect(KEYPAD_LAYOUT[5]).toContain(equals[0])
  })

  it('las teclas de entrada llevan su valor de expresión', () => {
    const byId = Object.fromEntries(KEYPAD_LAYOUT.flat().map((k) => [k.id, k]))
    expect(byId.divide?.value).toBe('/')
    expect(byId.multiply?.value).toBe('*')
    expect(byId.minus?.value).toBe('-')
    expect(byId.plus?.value).toBe('+')
    expect(byId.sqrt?.value).toBe('sqrt(')
    expect(byId.power?.value).toBe('^')
    expect(byId.percent?.value).toBe('percent(')
    expect(byId['open-paren']).toMatchObject({ label: '(', value: '(', type: 'operation' })
    expect(byId['close-paren']).toMatchObject({ label: ')', value: ')', type: 'operation' })
  })
})
