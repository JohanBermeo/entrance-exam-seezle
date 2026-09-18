import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { Display } from './Display'

describe('Display', () => {
  it('muestra la expresión y el resultado formateado es-ES', () => {
    render(<Display expression="5 + 3" result={1234.56} />)
    expect(screen.getByText('5 + 3')).toBeInTheDocument()
    expect(screen.getByText('1234,56')).toBeInTheDocument()
  })

  it('vacía el renglón de operación cuando duplica al resultado', () => {
    render(<Display expression="8" result={8} />)
    expect(screen.getByText('8', { selector: '.display-value' })).toBeInTheDocument()
    const row = document.querySelector('.display-expression')
    expect(row?.textContent).toBe(' ')
  })

  it('muestra 0 cuando no hay resultado', () => {
    render(<Display expression="" result={null} />)
    expect(screen.getByText('0')).toBeInTheDocument()
  })

  it('muestra estado de carga', () => {
    render(<Display expression="5 + 3" result={null} isLoading />)
    expect(screen.getByLabelText('Calculando')).toBeInTheDocument()
    expect(screen.getByText('5 + 3').parentElement).toHaveAttribute('aria-busy', 'true')
  })
})
