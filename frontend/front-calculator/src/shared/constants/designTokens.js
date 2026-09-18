/**
 * Design tokens v1 — fuente JS de los tokens acordados en
 * docs/frontend-calculator-plan.md. El CSS en `src/index.css`
 * es la representación efectiva; este módulo existe para que
 * componentes/utils puedan importar valores (p. ej. límites,
 * tamaños) sin hardcodearlos.
 */
export const designTokens = {
  color: {
    accent: '#833AED',
    accentLight: 'rgba(131, 58, 237, 0.12)',
    accentBorder: 'rgba(131, 58, 237, 0.2)',
    accentDark: '#6B21A8',
    accentDarkHover: '#581C87',
    white: '#FFFFFF',
    textPrimary: '#1A1A1A',
    textSecondary: '#6B6375',
    borderLight: '#E5E4E7',
    success: '#10B981',
    warning: '#F59E0B',
    error: '#EF4444',
    info: '#3B82F6',
  },
  font: {
    mono: "'JetBrains Mono', ui-monospace, Consolas, monospace",
    displaySize: 48,
    displaySizeMobile: 36,
    keySize: 20,
    keySizeSmall: 16,
  },
  space: {
    keypadGap: 8,
    modulePadding: 24,
  },
  radius: {
    btn: 12,
    display: 16,
    module: 20,
  },
  layout: {
    moduleMaxWidth: 360,
    displayHeight: 100,
    keySize: 64,
    keySizeMobile: 56,
  },
  breakpoints: {
    mobile: 480,
    tablet: 768,
    desktop: 1024,
  },
}
