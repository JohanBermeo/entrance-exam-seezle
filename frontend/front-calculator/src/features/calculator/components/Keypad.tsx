import { KEYPAD_LAYOUT } from '../utils/keypadLayout'
import type { KeyDef } from '../utils/keypadLayout'
import { Key } from './Key'

export interface KeypadProps {
  onKeyPress: (def: KeyDef) => void
  /** Durante la carga solo AC sigue habilitado. */
  disabled?: boolean
}

/** Grid 4×6 del keypad. Presentacional: delega la acción en `onKeyPress`. */
export function Keypad({ onKeyPress, disabled = false }: KeypadProps) {
  return (
    <div className="calc-keypad" role="group" aria-label="Teclado">
      {KEYPAD_LAYOUT.flat().map((def) => (
        <Key
          key={def.id}
          def={def}
          onPress={onKeyPress}
          disabled={disabled && def.kind !== 'clear'}
        />
      ))}
    </div>
  )
}
