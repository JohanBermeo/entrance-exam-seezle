import { cn } from '../utils/cn.js'

const TONES = { info: 'toast-info', success: 'toast-success', error: 'toast-error' }

export function Toast({ tone = 'info', title, message, onDismiss }) {
  return (
    <div role={tone === 'error' ? 'alert' : 'status'} className={cn('toast', TONES[tone] ?? TONES.info)}>
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
