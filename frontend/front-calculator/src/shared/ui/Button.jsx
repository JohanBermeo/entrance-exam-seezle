import { cn } from '../utils/cn.js'

const VARIANTS = {
  primary: 'btn-primary',
  secondary: 'btn-secondary',
  ghost: 'btn-ghost',
  danger: 'btn-danger',
}

const SIZES = { sm: 'btn-sm', md: 'btn-md', lg: 'btn-lg' }

export function Button({
  variant = 'primary',
  size = 'md',
  className,
  type = 'button',
  children,
  ...rest
}) {
  return (
    <button
      type={type}
      className={cn('btn', VARIANTS[variant] ?? VARIANTS.primary, SIZES[size] ?? SIZES.md, className)}
      {...rest}
    >
      {children}
    </button>
  )
}
