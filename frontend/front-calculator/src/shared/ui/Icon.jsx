import {
  DeleteIcon,
  SqrtIcon,
  PowerIcon,
  PlusMinusIcon,
  PercentIcon,
} from '../../assets/icons/index.js'

const ICONS = {
  delete: DeleteIcon,
  sqrt: SqrtIcon,
  power: PowerIcon,
  'plus-minus': PlusMinusIcon,
  percent: PercentIcon,
}

/** Wrapper tipado para los SVGs de `assets/icons`. */
export function Icon({ name, size = 20, className }) {
  const Cmp = ICONS[name]
  if (!Cmp) return null
  return <Cmp size={size} className={className} />
}
