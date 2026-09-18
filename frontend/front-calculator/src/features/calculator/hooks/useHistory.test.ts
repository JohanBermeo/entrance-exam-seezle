import { act, renderHook } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { HISTORY_LIMIT, useHistory } from './useHistory'

describe('useHistory', () => {
  it('acumula entradas con el más reciente primero', () => {
    const { result } = renderHook(() => useHistory())
    act(() => {
      result.current.push({ expression: '1+1', value: 2, requestId: 'a' })
      result.current.push({ expression: '2+2', value: 4, requestId: 'b' })
    })
    expect(result.current.entries).toHaveLength(2)
    expect(result.current.entries[0]?.expression).toBe('2+2')
    expect(result.current.entries[0]?.id).toBeTruthy()
  })

  it('recorta al límite (20 por defecto)', () => {
    const { result } = renderHook(() => useHistory())
    act(() => {
      for (let i = 0; i < HISTORY_LIMIT + 5; i++) {
        result.current.push({ expression: `${i}`, value: i, requestId: `${i}` })
      }
    })
    expect(result.current.entries).toHaveLength(HISTORY_LIMIT)
    expect(result.current.entries[0]?.expression).toBe(`${HISTORY_LIMIT + 4}`)
  })

  it('limpia el historial', () => {
    const { result } = renderHook(() => useHistory())
    act(() => {
      result.current.push({ expression: '1+1', value: 2, requestId: 'a' })
      result.current.clear()
    })
    expect(result.current.entries).toHaveLength(0)
  })
})
