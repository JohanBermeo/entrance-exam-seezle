import { Icon } from '../../../components/index'
import { cn } from '../../../shared/utils/cn'
import type { KeyDef } from '../utils/keypadLayout'

export interface KeyProps {
  def: KeyDef
  onPress: (def: KeyDef) => void
  disabled?: boolean
}

/** Tecla única del keypad. Presentacional: no conoce la expresión. */
export function Key({ def, onPress, disabled = false }: KeyProps) {
  if (def.kind === 'empty') {
    return <span className="key key-empty" aria-hidden="true" />
  }
  return (
    <button
      type="button"
      className={cn('key', `key-${def.type}`, def.span === 2 && 'key-span-2')}
      disabled={disabled}
      onClick={() => onPress(def)}
      aria-label={def.label || def.id}
    >
      {def.icon ? <Icon name={def.icon} size={22} /> : def.label}
    </button>
  )
}
