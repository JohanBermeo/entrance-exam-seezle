import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { CalculatorLayout } from './CalculatorLayout'

describe('CalculatorLayout', () => {
  it('envuelve el contenido en el módulo', () => {
    render(
      <CalculatorLayout>
        <p>contenido</p>
      </CalculatorLayout>,
    )
    expect(screen.getByRole('main')).toHaveClass('calc-module')
    expect(screen.getByText('contenido')).toBeInTheDocument()
  })
})
