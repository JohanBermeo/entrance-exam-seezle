import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import type { HistoryEntry } from '../hooks/useHistory'
import { HistoryPanel } from './HistoryPanel'

const entries: HistoryEntry[] = [
  { id: '1', expression: '5 + 3', value: 8, requestId: 'a', createdAt: 1 },
  { id: '2', expression: '2 ^ 10', value: 1024, requestId: 'b', createdAt: 2 },
]

describe('HistoryPanel', () => {
  it('muestra estado vacío', () => {
    render(<HistoryPanel entries={[]} onSelect={() => {}} onClear={() => {}} />)
    expect(screen.getByText('Sin cálculos todavía')).toBeInTheDocument()
  })

  it('lista entradas y notifica la selección', () => {
    const onSelect = vi.fn()
    render(<HistoryPanel entries={entries} onSelect={onSelect} onClear={() => {}} />)
    fireEvent.click(screen.getByRole('button', { name: /5 \+ 3/ }))
    expect(onSelect).toHaveBeenCalledWith(entries[0])
  })

  it('notifica limpiar', () => {
    const onClear = vi.fn()
    render(<HistoryPanel entries={entries} onSelect={() => {}} onClear={onClear} />)
    fireEvent.click(screen.getByRole('button', { name: 'Limpiar' }))
    expect(onClear).toHaveBeenCalled()
  })

  it('se colapsa con el toggle', () => {
    render(<HistoryPanel entries={entries} onSelect={() => {}} onClear={() => {}} />)
    const toggle = screen.getByRole('button', { name: /Historial/ })
    fireEvent.click(toggle)
    expect(toggle).toHaveAttribute('aria-expanded', 'false')
    expect(screen.queryByText('5 + 3')).toBeNull()
  })
})
