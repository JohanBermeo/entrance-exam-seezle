# Uso de inteligencia artificial en el proyecto

> Documento explicativo del uso de IA durante el desarrollo. Registros de
> sesión: [dialogo-sesion-hitos-01-05.md](./dialogo-sesion-hitos-01-05.md) y
> [dialogo-sesion-frontend-f1-f5.md](./dialogo-sesion-frontend-f1-f5.md).

## 1. Enfoque general

Con el tiempo disponible para la prueba, usé la IA como **acelerador con
supervisión**, no como piloto automático. Combiné dos tipos de IA:

- **IA generativa** → idea y referencia visual del diseño.
- **IA agéntica** (Codex y OpenCode) → planificación, guardrails y código,
  siempre dentro de un protocolo con aprobación humana por hito.

## 2. IA generativa: diseño con Figma Make

Realizar un diseño adecuado a mano con el tiempo disponible era complicado.
Como la aplicación no requiere un diseño complejo y la IA de Figma se
especializa justo en esos casos, decidí usar **Figma Make** para obtener una
idea de cómo quería que se viera la aplicación.

Del mockup «Calculadora web estética» extraje la referencia visual del
proyecto: módulo blanco centrado con sombreado morado, acento `#833AED`,
keypad con orden definido y tipografía JetBrains Mono. Las decisiones finales
(colores, distribución, iconos) las aprobé yo sobre esa base.

- Referencia: mockup «Calculadora web estética» (Figma Make).

## 3. IA agéntica: guardrails antes que código

Usé IAs agénticas como **Codex** y **OpenCode** para generar unos guardrails
que permitieran acelerar la fase de desarrollo y evitar equivocaciones en el
mayor porcentaje posible:

- Planes canónicos versionados: [backend-calculator-plan.md](./backend-calculator-plan.md),
  [frontend-calculator-plan.md](./frontend-calculator-plan.md).
- Protocolo de trabajo: [development-protocol.md](./development-protocol.md)
  (un hito a la vez, ramas cortas, verificación obligatoria:
  lint, typecheck, tests y build en verde).
- [AGENTS.md](../AGENTS.md) como punto de entrada para que cualquier agente
  (Codex, OpenCode, Copilot) siga la misma fuente de verdad sin derivas
  silenciosas.

## 4. Recorrido de modelos

- **gpt 5.6 terra alta** → planificación de todo el desarrollo del backend.
  Lo elegí para aprovechar sus altas capacidades de procesamiento sin
  malgastar tokens, ya que tengo una cantidad limitada.
- **luna medio** → realización del código en sus primeras instancias, hasta
  que se acabó la cantidad de tokens disponibles.
- **OpenCode + nemotron 3 ultra (modo plan)** → definir una forma eficiente
  de migrar el canvas generado por Codex a un esquema de archivos más
  estandarizado para el uso general de las IAs agénticas, y posteriormente
  para la planificación de toda la fase frontend (F1–F5).
- **muse spark 1.3 (modo build)** → implementación: aproveché sus bondades
  como la optimización en la generación de código para llevar a cabo la
  misma tarea (frontend F1–F5, fixes de keypad y display, Docker y `run.ps1`).

Evidencia del trabajo agéntico: chat compartido de Codex
([https://chatgpt.com/s/cx_6aace32fb38c8191bbab0f0feb29fd34](https://chatgpt.com/s/cx_6aace32fb38c8191bbab0f0feb29fd34)) y los
registros de sesión citados arriba.

## 5. Humano en el loop

Todas las decisiones de producto las tomé yo: aprobación del diseño y sus
tokens, aprobación de cada hito antes de continuar, merges y creación de PRs,
y decisiones finas (ubicación de paréntesis, comportamiento tras `=`,
`+/-` con paréntesis). La IA ejecutó dentro de esos guardrails.
