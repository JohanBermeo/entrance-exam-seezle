import { describe, expect, it } from 'vitest'
import { formatNumber } from './formatNumber'

describe('formatNumber', () => {
  it('formatea con locale es-ES', () => {
    expect(formatNumber(1234.56)).toBe('1234,56')
  })

  it('usa notación científica para magnitudes extremas', () => {
    expect(formatNumber(1e13)).toContain('e+')
    expect(formatNumber(0.0000001)).toContain('e-')
  })

  it('maneja no-finitos', () => {
    expect(formatNumber(NaN)).toBe('—')
    expect(formatNumber(Infinity)).toBe('Error')
  })
})
