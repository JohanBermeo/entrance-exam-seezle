/** Nombres de operador de la API (docs/backend-calculator-plan.md). */
export type Operator =
  | 'add'
  | 'subtract'
  | 'multiply'
  | 'divide'
  | 'power'
  | 'sqrt'
  | 'percent'

export interface CalculationRequest {
  expression: string
  outputs?: string[]
}

export interface OperationError {
  code: string
  message: string
  position?: number
}

export interface OperationResult {
  id: string
  value?: number
  error?: OperationError
}

export interface CalculationResponse {
  results: Record<string, OperationResult>
  outputs: string[]
  requestId: string
  durationMs: number
}

/** Error normalizado que consume la UI. */
export type CalculationError =
  | { type: 'validation'; message: string; position?: number } // HTTP 400
  | { type: 'domain'; code: string; message: string } // HTTP 422
  | { type: 'network'; message: string }
  | { type: 'timeout' }
  | { type: 'unknown'; message: string }
