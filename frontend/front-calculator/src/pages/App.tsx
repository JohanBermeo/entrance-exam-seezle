import { useState } from 'react'
import { Button, Card, Icon } from '../components/index'
import type { IconName } from '../components/index'
import { Toast } from '../components/index'
import './App.css'

const ICON_NAMES: IconName[] = ['delete']

/**
 * Shell visual de F1 · Fundaciones: verifica tokens, components e iconos.
 * El layout real de la calculadora (Display/Keypad) llega en F3.
 */
function App() {
  const [showToast, setShowToast] = useState(true)

  return (
    <main className="calc-module app-shell">
      <header className="app-header">
        <p className="app-kicker">Calculadora · F1 Fundaciones</p>
        <h1 className="app-title">components + tokens</h1>
      </header>

      <Card className="app-preview">
        <div className="calc-display app-display" aria-label="Display de ejemplo">
          0
        </div>
        <div className="app-row">
          <Button variant="primary" size="md">Primary</Button>
          <Button variant="secondary" size="md">Secondary</Button>
          <Button variant="ghost" size="md">Ghost</Button>
          <Button variant="danger" size="sm">Danger</Button>
        </div>
        <div className="app-row app-icons">
          {ICON_NAMES.map((name) => (
            <span key={name} className="app-icon-chip" title={name}>
              <Icon name={name} size={22} />
            </span>
          ))}
        </div>
        <div className="app-row">
          <button type="button" className="key key-number">7</button>
          <button type="button" className="key key-operation">÷</button>
          <button type="button" className="key key-function">AC</button>
          <button type="button" className="key key-equals">=</button>
        </div>
      </Card>

      {showToast ? (
        <Toast
          tone="info"
          title="F1 lista para revisión visual"
          message="Tokens, Button, Card, Icon, Toast y el icono provisto como componente."
          onDismiss={() => setShowToast(false)}
        />
      ) : null}
    </main>
  )
}

export default App
