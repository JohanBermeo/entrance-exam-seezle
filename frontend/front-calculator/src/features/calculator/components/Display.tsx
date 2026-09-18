import { formatNumber } from '../../../shared/utils/formatNumber'
import { cn } from '../../../shared/utils/cn'

export interface DisplayProps {
  expression: string
  result: number | null
  isLoading?: boolean
}

/**
 * Display superior: expresión en curso + resultado formateado es-ES.
 * Tras `=`, la expresión pasa a ser el resultado para encadenar; en ese caso
 * el renglón de operación queda vacío y solo se muestra el resultado.
 */
export function Display({ expression, result, isLoading = false }: DisplayProps) {
  const showExpression = result === null || expression !== String(result)
  return (
    <div className="calc-display" aria-live="polite" aria-busy={isLoading}>
      <span className="display-expression">{showExpression ? expression || ' ' : ' '}</span>
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
