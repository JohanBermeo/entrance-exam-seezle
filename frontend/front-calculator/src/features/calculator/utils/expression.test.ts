import { describe, expect, it } from 'vitest'
import { toggleLastNumberSign } from './expression'

describe('toggleLastNumberSign', () => {
  it('niega un número simple', () => {
    expect(toggleLastNumberSign('5')).toBe('-5')
    expect(toggleLastNumberSign('-5')).toBe('5')
  })

  it('niega el último número tras un operador', () => {
    expect(toggleLastNumberSign('5+3')).toBe('5+-3')
    expect(toggleLastNumberSign('5+-3')).toBe('5+3')
  })

  it('no confunde resta con signo', () => {
    expect(toggleLastNumberSign('5-3')).toBe('5--3')
  })

  it('niega tras paréntesis', () => {
    expect(toggleLastNumberSign('(2+3)*4')).toBe('(2+3)*-4')
  })

  it('sin número al final niega toda la expresión', () => {
    expect(toggleLastNumberSign('')).toBe('-')
    expect(toggleLastNumberSign('(2+3)')).toBe('-(2+3)')
  })
})
