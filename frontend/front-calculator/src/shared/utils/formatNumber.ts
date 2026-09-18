const esFormatter = new Intl.NumberFormat('es-ES', { maximumFractionDigits: 10 })

/**
 * Formatea un número con locale es-ES.
 * Usa notación científica para magnitudes muy grandes/pequeñas.
 */
export function formatNumber(value: number): string {
  if (typeof value !== 'number' || Number.isNaN(value)) return '—'
  if (!Number.isFinite(value)) return 'Error'
  const abs = Math.abs(value)
  if (abs !== 0 && (abs >= 1e12 || abs < 1e-6)) {
    return value.toExponential(6).replace('.', ',')
  }
  return esFormatter.format(value)
}
