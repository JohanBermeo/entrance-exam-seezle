import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import type { CalculationError } from '../api/types'
import { ErrorToast } from './ErrorToast'

describe('ErrorToast', () => {
  it('no renderiza sin error', () => {
    const { container } = render(<ErrorToast error={null} onDismiss={() => {}} />)
    expect(container).toBeEmptyDOMElement()
  })

  it('muestra validación con posición y rol alert', () => {
    const error: CalculationError = {
      type: 'validation',
      message: 'token inesperado',
      position: 4,
    }
    render(<ErrorToast error={error} onDismiss={() => {}} />)
    expect(screen.getByRole('alert')).toHaveTextContent('Expresión inválida')
    expect(screen.getByRole('alert')).toHaveTextContent('posición 4')
  })

  it('muestra dominio con código', () => {
    const error: CalculationError = {
      type: 'domain',
      code: 'division_by_zero',
      message: 'divisor cero',
    }
    render(<ErrorToast error={error} onDismiss={() => {}} />)
    expect(screen.getByRole('alert')).toHaveTextContent('division_by_zero')
  })

  it('notifica el cierre', () => {
    const onDismiss = vi.fn()
    const error: CalculationError = { type: 'timeout' }
    render(<ErrorToast error={error} onDismiss={onDismiss} />)
    fireEvent.click(screen.getByRole('button', { name: 'Cerrar aviso' }))
    expect(onDismiss).toHaveBeenCalled()
  })
})
