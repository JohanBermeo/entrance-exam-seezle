import type { HTMLAttributes } from 'react'
import { cn } from '../shared/utils/cn'

export type CardProps = HTMLAttributes<HTMLElement>

export function Card({ className, children, ...rest }: CardProps) {
  return (
    <section className={cn('card', className)} {...rest}>
      {children}
    </section>
  )
}
