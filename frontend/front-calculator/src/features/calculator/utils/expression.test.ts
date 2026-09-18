import { describe, expect, it } from 'vitest'
import { appendSymbol, toggleLastNumberSign } from './expression'

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

describe('appendSymbol', () => {
  it('inserta × explícito ante abridor tras dígito o cierre', () => {
    expect(appendSymbol('9', 'sqrt(')).toBe('9*sqrt(')
    expect(appendSymbol('5', 'percent(')).toBe('5*percent(')
    expect(appendSymbol('(2+3)', '(')).toBe('(2+3)*(')
  })

  it('no inserta × en el resto de casos', () => {
    expect(appendSymbol('', 'sqrt(')).toBe('sqrt(')
    expect(appendSymbol('5+', 'sqrt(')).toBe('5+sqrt(')
    expect(appendSymbol('9', '5')).toBe('95')
    expect(appendSymbol('9', '+')).toBe('9+')
  })
})
