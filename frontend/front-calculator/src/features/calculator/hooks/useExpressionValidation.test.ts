import { renderHook } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { MAX_EXPRESSION_LENGTH, useExpressionValidation } from './useExpressionValidation'

describe('useExpressionValidation', () => {
  it('acepta una expresión válida con funciones', () => {
    const { result } = renderHook(() => useExpressionValidation())
    expect(result.current.validate('sqrt(percent(200, 15)) + 4 ^ 2')).toEqual({ valid: true })
  })

  it('rechaza la expresión vacía', () => {
    const { result } = renderHook(() => useExpressionValidation())
    expect(result.current.validate('   ')).toEqual({ valid: false, error: 'Expresión vacía' })
  })

  it('rechaza caracteres no permitidos con posición', () => {
    const { result } = renderHook(() => useExpressionValidation())
    expect(result.current.validate('5 + §')).toEqual({
      valid: false,
      error: 'Caracteres no permitidos',
      position: 4,
    })
  })

  it('rechaza paréntesis sin cerrar', () => {
    const { result } = renderHook(() => useExpressionValidation())
    expect(result.current.validate('((2 + 3)')).toEqual({
      valid: false,
      error: 'Paréntesis sin cerrar',
    })
  })

  it('rechaza cierre sin apertura con posición', () => {
    const { result } = renderHook(() => useExpressionValidation())
    expect(result.current.validate('2 + 3)')).toEqual({
      valid: false,
      error: 'Paréntesis de cierre sin apertura',
      position: 5,
    })
  })

  it('rechaza funciones desconocidas con posición', () => {
    const { result } = renderHook(() => useExpressionValidation())
    expect(result.current.validate('foo(2)')).toEqual({
      valid: false,
      error: 'Función desconocida: foo',
      position: 0,
    })
  })

  it('rechaza expresiones que superan la longitud máxima', () => {
    const { result } = renderHook(() => useExpressionValidation())
    const issue = result.current.validate('1+'.repeat(MAX_EXPRESSION_LENGTH))
    expect(issue.valid).toBe(false)
    expect(issue.error).toContain('demasiado larga')
  })
})
