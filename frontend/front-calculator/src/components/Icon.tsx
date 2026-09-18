import { DeleteIcon } from '../assets/icons/index'

/** Solo iconos provistos como asset. Las teclas de función usan etiqueta de texto. */
export type IconName = 'delete'

const ICONS = {
  delete: DeleteIcon,
} as const

export interface IconProps {
  name: IconName
  size?: number
  className?: string
}

export function Icon({ name, size = 20, className }: IconProps) {
  const Cmp = ICONS[name]
  return <Cmp size={size} className={className} />
}
