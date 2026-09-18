import { Toast } from '../../../components/index'
import type { CalculationError } from '../api/types'

export interface ErrorToastProps {
  error: CalculationError | null
  onDismiss: () => void
}

function describeError(error: CalculationError): { title: string; message: string } {
  switch (error.type) {
    case 'validation':
      return {
        title: 'Expresión inválida',
        message:
          error.position !== undefined
            ? `${error.message} (posición ${error.position})`
            : error.message,
      }
    case 'domain':
      return { title: 'Error de cálculo', message: `${error.message} (${error.code})` }
    case 'network':
      return { title: 'Sin conexión', message: error.message }
    case 'timeout':
      return {
        title: 'Tiempo agotado',
        message: 'El cálculo tardó demasiado. Inténtalo de nuevo.',
      }
    case 'unknown':
      return { title: 'Error inesperado', message: error.message }
  }
}

/** Traduce CalculationError a Toast visible. Nada si no hay error. */
export function ErrorToast({ error, onDismiss }: ErrorToastProps) {
  if (!error) return null
  const { title, message } = describeError(error)
  return <Toast tone="error" title={title} message={message} onDismiss={onDismiss} />
}
