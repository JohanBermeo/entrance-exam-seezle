import { fireEvent, render, screen, within } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { Keypad } from './Keypad'

describe('Keypad', () => {
  it('renderiza las 24 teclas sin celda vacía', () => {
    render(<Keypad onKeyPress={() => {}} />)
    const group = screen.getByRole('group', { name: 'Teclado' })
    expect(within(group).getAllByRole('button')).toHaveLength(24)
  })

  it('los paréntesis son tipo operación y notifican su valor', () => {
    const onKeyPress = vi.fn()
    render(<Keypad onKeyPress={onKeyPress} />)
    const open = screen.getByRole('button', { name: '(' })
    expect(open).toHaveClass('key-operation')
    expect(open).not.toHaveClass('key-span-2')
    fireEvent.click(open)
    expect(onKeyPress).toHaveBeenCalledWith(
      expect.objectContaining({ id: 'open-paren', value: '(' }),
    )
  })

  it('notifica la tecla pulsada con su def', () => {
    const onKeyPress = vi.fn()
    render(<Keypad onKeyPress={onKeyPress} />)
    fireEvent.click(screen.getByRole('button', { name: '÷' }))
    expect(onKeyPress).toHaveBeenCalledWith(
      expect.objectContaining({ id: 'divide', value: '/' }),
    )
  })

  it('con disabled solo AC sigue habilitado', () => {
    render(<Keypad onKeyPress={() => {}} disabled />)
    expect(screen.getByRole('button', { name: 'AC' })).toBeEnabled()
    expect(screen.getByRole('button', { name: '7' })).toBeDisabled()
    expect(screen.getByRole('button', { name: '=' })).toBeDisabled()
  })
})
