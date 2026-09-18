import { act, fireEvent, renderHook } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { useKeypad } from './useKeypad'

describe('useKeypad', () => {
  it('construye la expresión con input/clear/backspace', () => {
    const { result } = renderHook(() => useKeypad())
    act(() => {
      result.current.input('5')
      result.current.input('+')
      result.current.input('3')
    })
    expect(result.current.expression).toBe('5+3')
    act(() => result.current.backspace())
    expect(result.current.expression).toBe('5+')
    act(() => result.current.clear())
    expect(result.current.expression).toBe('')
  })

  it('inserta × explícito al abrir función tras dígito', () => {
    const { result } = renderHook(() => useKeypad({ initialExpression: '9' }))
    act(() => result.current.input('sqrt('))
    expect(result.current.expression).toBe('9*sqrt(')
  })

  it('toggleSign envuelve el último número con paréntesis', () => {
    const { result } = renderHook(() => useKeypad({ initialExpression: '5+3' }))
    act(() => result.current.toggleSign())
    expect(result.current.expression).toBe('5+(-3)')
    act(() => result.current.toggleSign())
    expect(result.current.expression).toBe('5+3')
  })

  it('tras commitResult el número/paréntesis/punto borran y la operación encadena', () => {
    const { result } = renderHook(() => useKeypad())
    act(() => result.current.commitResult(8))
    expect(result.current.expression).toBe('8')
    act(() => result.current.input('5'))
    expect(result.current.expression).toBe('5')
    act(() => result.current.commitResult(8))
    act(() => result.current.input('+'))
    expect(result.current.expression).toBe('8+')
    act(() => result.current.commitResult(8))
    act(() => result.current.input('('))
    expect(result.current.expression).toBe('(')
    act(() => result.current.commitResult(8))
    act(() => result.current.input('.'))
    expect(result.current.expression).toBe('.')
  })

  it('editar invalida el borrado fresco', () => {
    const { result } = renderHook(() => useKeypad())
    act(() => result.current.commitResult(8))
    act(() => result.current.backspace())
    act(() => result.current.input('8'))
    act(() => result.current.input('5'))
    expect(result.current.expression).toBe('85')
  })

  it('submit llama a onEquals con la expresión actual', () => {
    const onEquals = vi.fn()
    const { result } = renderHook(() => useKeypad({ initialExpression: '2*4', onEquals }))
    act(() => result.current.submit())
    expect(onEquals).toHaveBeenCalledWith('2*4')
  })

  it('escucha el teclado físico', () => {
    const onEquals = vi.fn()
    const { result } = renderHook(() => useKeypad({ onEquals }))
    act(() => {
      fireEvent.keyDown(window, { key: '7' })
      fireEvent.keyDown(window, { key: '+' })
    })
    expect(result.current.expression).toBe('7+')
    act(() => {
      fireEvent.keyDown(window, { key: 'Enter' })
    })
    expect(onEquals).toHaveBeenCalledWith('7+')
    act(() => {
      fireEvent.keyDown(window, { key: 'Escape' })
    })
    expect(result.current.expression).toBe('')
  })

  it('ignora teclas con modificadores', () => {
    const { result } = renderHook(() => useKeypad())
    act(() => {
      fireEvent.keyDown(window, { key: '5', ctrlKey: true })
    })
    expect(result.current.expression).toBe('')
  })

  it('notifica onEdit en cada edición pero no en submit', () => {
    const onEdit = vi.fn()
    const onEquals = vi.fn()
    const { result } = renderHook(() => useKeypad({ onEdit, onEquals }))
    act(() => {
      result.current.input('5')
      result.current.backspace()
      result.current.toggleSign()
      result.current.clear()
    })
    expect(onEdit).toHaveBeenCalledTimes(4)
    act(() => result.current.submit())
    expect(onEdit).toHaveBeenCalledTimes(4)
    expect(onEquals).toHaveBeenCalledTimes(1)
  })
})
