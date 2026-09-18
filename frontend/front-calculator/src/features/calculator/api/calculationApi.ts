import type {
  CalculationError,
  CalculationResponse,
} from './types'

export const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080'

export interface CalculateOptions {
  outputs?: string[]
  signal?: AbortSignal
}

interface ErrorBody {
  code?: string
  message?: string
  position?: number
}

function mapHttpError(status: number, body: ErrorBody): CalculationError {
  switch (status) {
    case 400:
      return {
        type: 'validation',
        message: body.message ?? 'Expresión inválida',
        position: body.position,
      }
    case 422:
      return {
        type: 'domain',
        code: body.code ?? 'domain_error',
        message: body.message ?? 'Error de cálculo',
      }
    case 408:
    case 499:
      return { type: 'timeout' }
    default:
      return { type: 'unknown', message: body.message ?? 'Error inesperado' }
  }
}

/** Ejecuta una expresión contra `POST /v1/calculations` (modo expression). */
export async function calculateExpression(
  expression: string,
  options: CalculateOptions = {},
): Promise<CalculationResponse> {
  const { outputs = ['result'], signal } = options
  let res: Response
  try {
    res = await fetch(`${API_BASE_URL}/v1/calculations`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ expression, outputs }),
      signal,
    })
  } catch (e) {
    if (signal?.aborted) {
      const reason = signal.reason
      throw reason instanceof Error ? reason : new DOMException('Aborted', 'AbortError')
    }
    if (e instanceof DOMException && e.name === 'AbortError') throw e
    throw { type: 'network', message: 'Sin conexión con el servidor' } as CalculationError
  }

  if (!res.ok) {
    const body = (await res.json().catch(() => ({}))) as ErrorBody
    throw mapHttpError(res.status, body)
  }
  return (await res.json()) as CalculationResponse
}
