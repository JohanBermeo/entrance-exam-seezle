import { cn } from '../utils/cn.js'

export function Card({ className, children, ...rest }) {
  return (
    <section className={cn('card', className)} {...rest}>
      {children}
    </section>
  )
}
