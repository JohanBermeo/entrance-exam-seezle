import { cn } from '../shared/utils/cn'

export type ToastTone = 'info' | 'success' | 'error'

export interface ToastProps {
  tone?: ToastTone
  title?: string
  message?: string
  onDismiss?: () => void
}

const TONES: Record<ToastTone, string> = {
  info: 'toast-info',
  success: 'toast-success',
  error: 'toast-error',
}

export function Toast({ tone = 'info', title, message, onDismiss }: ToastProps) {
  return (
    <div role={tone === 'error' ? 'alert' : 'status'} className={cn('toast', TONES[tone])}>
      <div className="toast-body">
        {title ? <strong className="toast-title">{title}</strong> : null}
        {message ? <p className="toast-message">{message}</p> : null}
      </div>
      {onDismiss ? (
        <button type="button" className="toast-close" onClick={onDismiss} aria-label="Cerrar aviso">
          ×
        </button>
      ) : null}
    </div>
  )
}
