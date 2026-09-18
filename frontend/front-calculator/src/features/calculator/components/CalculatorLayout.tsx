import type { ReactNode } from 'react'

export interface CalculatorLayoutProps {
  children: ReactNode
}

/** Módulo centrado de la calculadora (sombra morada, max-width 360px). */
export function CalculatorLayout({ children }: CalculatorLayoutProps) {
  return <main className="calc-module">{children}</main>
}
