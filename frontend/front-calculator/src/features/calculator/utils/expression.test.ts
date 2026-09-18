import { describe, expect, it } from 'vitest'
import { appendSymbol, toggleLastNumberSign } from './expression'

describe('toggleLastNumberSign', () => {
  it('envuelve el número con paréntesis y lo quita al repetir', () => {
    expect(toggleLastNumberSign('5')).toBe('(-5)')
    expect(toggleLastNumberSign('(-5)')).toBe('5')
    expect(toggleLastNumberSign('-5')).toBe('5')
  })

  it('niega el último número con paréntesis', () => {
    expect(toggleLastNumberSign('5+3')).toBe('5+(-3)')
    expect(toggleLastNumberSign('5+(-3)')).toBe('5+3')
    expect(toggleLastNumberSign('5+-3')).toBe('5+3')
    expect(toggleLastNumberSign('5-3')).toBe('5-(-3)')
    expect(toggleLastNumberSign('5-(-3)')).toBe('5-3')
    expect(toggleLastNumberSign('5--3')).toBe('5-3')
  })

  it('niega grupos y llamadas completas', () => {
    expect(toggleLastNumberSign('(2+3)')).toBe('-(2+3)')
    expect(toggleLastNumberSign('-(2+3)')).toBe('(2+3)')
    expect(toggleLastNumberSign('5+(2+3)')).toBe('5+(-(2+3))')
    expect(toggleLastNumberSign('5+(-(2+3))')).toBe('5+(2+3)')
    expect(toggleLastNumberSign('sqrt(9)')).toBe('-sqrt(9)')
    expect(toggleLastNumberSign('-sqrt(9)')).toBe('sqrt(9)')
  })

  it('niega tras operador o paréntesis', () => {
    expect(toggleLastNumberSign('(2+3)*4')).toBe('(2+3)*(-4)')
    expect(toggleLastNumberSign('2^3')).toBe('2^(-3)')
    expect(toggleLastNumberSign('percent(200, 15)')).toBe('-percent(200, 15)')
    expect(toggleLastNumberSign('percent(200, 15')).toBe('percent(200, (-15)')
  })

  it('sin operando final niega toda la expresión', () => {
    expect(toggleLastNumberSign('')).toBe('-')
    expect(toggleLastNumberSign('-')).toBe('')
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
