import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { afterAll, afterEach, beforeAll, describe, expect, it } from 'vitest'
import { API_BASE_URL, calculateExpression } from './calculationApi'
import type { CalculationError, CalculationResponse } from './types'

const URL = `${API_BASE_URL}/v1/calculations`

const okBody: CalculationResponse = {
  results: { result: { id: 'result', value: 8 } },
  outputs: ['result'],
  requestId: 'req-1',
  durationMs: 3,
}

const server = setupServer(
  http.post(URL, () => HttpResponse.json(okBody)),
)

beforeAll(() => server.listen())
afterEach(() => server.resetHandlers())
afterAll(() => server.close())

describe('calculateExpression', () => {
  it('devuelve la respuesta en modo expression', async () => {
    const res = await calculateExpression('5 + 3')
    expect(res.results.result?.value).toBe(8)
    expect(res.requestId).toBe('req-1')
  })

  it('mapea 400 a error de validación con posición', async () => {
    server.use(
      http.post(URL, () =>
        HttpResponse.json({ message: 'token inesperado', position: 4 }, { status: 400 }),
      ),
    )
    const err = await calculateExpression('5 +').catch((e: CalculationError) => e)
    expect(err).toEqual({ type: 'validation', message: 'token inesperado', position: 4 })
  })

  it('mapea 422 a error de dominio con código', async () => {
    server.use(
      http.post(URL, () =>
        HttpResponse.json({ code: 'division_by_zero', message: 'divisor cero' }, { status: 422 }),
      ),
    )
    const err = await calculateExpression('5 / 0').catch((e: CalculationError) => e)
    expect(err).toEqual({ type: 'domain', code: 'division_by_zero', message: 'divisor cero' })
  })

  it('mapea 408 a timeout', async () => {
    server.use(http.post(URL, () => new HttpResponse(null, { status: 408 })))
    const err = await calculateExpression('1 + 1').catch((e: CalculationError) => e)
    expect(err).toEqual({ type: 'timeout' })
  })

  it('mapea cuerpo no-JSON a error desconocido', async () => {
    server.use(http.post(URL, () => new HttpResponse('boom', { status: 500 })))
    const err = await calculateExpression('1 + 1').catch((e: CalculationError) => e)
    expect(err).toEqual({ type: 'unknown', message: 'Error inesperado' })
  })

  it('mapea fallo de red a error de red', async () => {
    server.use(http.post(URL, () => HttpResponse.error()))
    const err = await calculateExpression('1 + 1').catch((e: CalculationError) => e)
    expect(err).toEqual({ type: 'network', message: 'Sin conexión con el servidor' })
  })

  it('propaga el abort sin normalizarlo', async () => {
    await expect(
      calculateExpression('1 + 1', { signal: AbortSignal.abort() }),
    ).rejects.toMatchObject({ name: 'AbortError' })
  })
})
