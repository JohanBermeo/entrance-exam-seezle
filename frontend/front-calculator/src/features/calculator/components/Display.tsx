import { formatNumber } from '../../../shared/utils/formatNumber'
import { cn } from '../../../shared/utils/cn'

export interface DisplayProps {
  expression: string
  result: number | null
  isLoading?: boolean
}

/** Display superior: expresión en curso + resultado formateado es-ES. */
export function Display({ expression, result, isLoading = false }: DisplayProps) {
  return (
    <div className="calc-display" aria-live="polite" aria-busy={isLoading}>
      <span className="display-expression">{expression || ' '}</span>
      {isLoading ? (
        <span className={cn('display-value', 'display-loading')} aria-label="Calculando">
          …
        </span>
      ) : result !== null ? (
        <span className="display-value">{formatNumber(result)}</span>
      ) : (
        <span className={cn('display-value', 'display-placeholder')}>0</span>
      )}
    </div>
  )
}
