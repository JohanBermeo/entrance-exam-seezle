import { act, renderHook } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { afterAll, afterEach, beforeAll, describe, expect, it } from 'vitest'
import { API_BASE_URL } from '../api/calculationApi'
import type { CalculationResponse } from '../api/types'
import { useCalculation } from './useCalculation'

const URL = `${API_BASE_URL}/v1/calculations`

const okBody: CalculationResponse = {
  results: { result: { id: 'result', value: 8 } },
  outputs: ['result'],
  requestId: 'req-1',
  durationMs: 3,
}

const server = setupServer(http.post(URL, () => HttpResponse.json(okBody)))

beforeAll(() => server.listen())
afterEach(() => server.resetHandlers())
afterAll(() => server.close())

describe('useCalculation', () => {
  it('calcula, guarda resultado e historial', async () => {
    const { result } = renderHook(() => useCalculation())
    let value: number | null | undefined
    await act(async () => {
      value = await result.current.calculate('5 + 3')
    })
    expect(value).toBe(8)
    expect(result.current.isLoading).toBe(false)
    expect(result.current.error).toBeNull()
    expect(result.current.result?.results.result?.value).toBe(8)
    expect(result.current.history).toHaveLength(1)
    expect(result.current.history[0]).toMatchObject({
      expression: '5 + 3',
      value: 8,
      requestId: 'req-1',
    })
  })

  it('expone el error de dominio y no guarda historial', async () => {
    server.use(
      http.post(URL, () =>
        HttpResponse.json({ code: 'division_by_zero', message: 'divisor cero' }, { status: 422 }),
      ),
    )
    const { result } = renderHook(() => useCalculation())
    let value: number | null | undefined
    await act(async () => {
      value = await result.current.calculate('5 / 0')
    })
    expect(value).toBeNull()
    expect(result.current.result).toBeNull()
    expect(result.current.error).toEqual({
      type: 'domain',
      code: 'division_by_zero',
      message: 'divisor cero',
    })
    expect(result.current.history).toHaveLength(0)
  })

  it('marca isLoading durante la petición', async () => {
    let release!: () => void
    const gate = new Promise<void>((resolve) => {
      release = resolve
    })
    server.use(
      http.post(URL, async () => {
        await gate
        return HttpResponse.json(okBody)
      }),
    )
    const { result } = renderHook(() => useCalculation())
    let pending: Promise<number | null> | undefined
    act(() => {
      pending = result.current.calculate('1 + 1')
    })
    expect(result.current.isLoading).toBe(true)
    await act(async () => {
      release()
      await pending
    })
    expect(result.current.isLoading).toBe(false)
    expect(result.current.result?.requestId).toBe('req-1')
  })

  it('limpia el historial', async () => {
    const { result } = renderHook(() => useCalculation())
    await act(async () => {
      await result.current.calculate('1 + 1')
    })
    expect(result.current.history).toHaveLength(1)
    act(() => result.current.clearHistory())
    expect(result.current.history).toHaveLength(0)
  })
})
