import { useState } from 'react'
import { formatNumber } from '../../../shared/utils/formatNumber'
import type { HistoryEntry } from '../hooks/useHistory'

export interface HistoryPanelProps {
  entries: HistoryEntry[]
  onSelect: (entry: HistoryEntry) => void
  onClear: () => void
}

/** Historial en memoria, colapsable. Clic para reutilizar la expresión. */
export function HistoryPanel({ entries, onSelect, onClear }: HistoryPanelProps) {
  const [open, setOpen] = useState(true)

  return (
    <section className="history-panel" aria-label="Historial">
      <div className="history-header">
        <button
          type="button"
          className="history-toggle"
          aria-expanded={open}
          onClick={() => setOpen((v) => !v)}
        >
          Historial ({entries.length})
        </button>
        {entries.length > 0 ? (
          <button type="button" className="history-clear" onClick={onClear}>
            Limpiar
          </button>
        ) : null}
      </div>
      {open ? (
        entries.length === 0 ? (
          <p className="history-empty">Sin cálculos todavía</p>
        ) : (
          <ul className="history-list">
            {entries.map((entry) => (
              <li key={entry.id}>
                <button
                  type="button"
                  className="history-item"
                  onClick={() => onSelect(entry)}
                  title="Reutilizar expresión"
                >
                  <span className="history-expression">{entry.expression}</span>
                  <span className="history-value">
                    = {entry.value !== null ? formatNumber(entry.value) : 'Error'}
                  </span>
                </button>
              </li>
            ))}
          </ul>
        )
      ) : null}
    </section>
  )
}
