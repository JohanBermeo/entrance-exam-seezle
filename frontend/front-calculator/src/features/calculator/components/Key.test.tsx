import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import type { KeyDef } from '../utils/keypadLayout'
import { Key } from './Key'

const def: KeyDef = { id: 'n7', label: '7', kind: 'input', type: 'number', value: '7' }

describe('Key', () => {
  it('renderiza la etiqueta y notifica la pulsación con su def', () => {
    const onPress = vi.fn()
    render(<Key def={def} onPress={onPress} />)
    fireEvent.click(screen.getByRole('button', { name: '7' }))
    expect(onPress).toHaveBeenCalledWith(def)
  })

  it('aplica la clase de variante y el span', () => {
    const onPress = vi.fn()
    render(
      <Key
        def={{ id: 'n0', label: '0', kind: 'input', type: 'number', value: '0', span: 2 }}
        onPress={onPress}
      />,
    )
    expect(screen.getByRole('button', { name: '0' })).toHaveClass('key-number', 'key-span-2')
  })

  it('respeta disabled', () => {
    const onPress = vi.fn()
    render(<Key def={def} onPress={onPress} disabled />)
    const btn = screen.getByRole('button', { name: '7' })
    expect(btn).toBeDisabled()
    fireEvent.click(btn)
    expect(onPress).not.toHaveBeenCalled()
  })

  it('la celda vacía no es un botón', () => {
    const onPress = vi.fn()
    render(<Key def={{ id: 'empty', label: '', kind: 'empty', type: 'empty' }} onPress={onPress} />)
    expect(screen.queryByRole('button')).toBeNull()
  })
})
