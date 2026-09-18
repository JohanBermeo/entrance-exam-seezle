# Plan de frontend · Calculadora React

**Fuente canónica versionada** del plan frontend (equivalente a `backend-calculator-plan.md`).
Cualquier agente de IA debe tratar este documento como referencia principal para el frontend.
Si el plan y el código discrepan, **no derivar en silencio**: alinear con el plan, preguntar al usuario o actualizar este archivo en el mismo cambio.

Stack: **React 19 + Vite**, `fetch` nativo + `useState`, sin librerías de estado/data-fetching en v1.

## Visión

UI de calculadora centrada en página: **display arriba + keypad abajo**, en un módulo blanco con sombreado morado leve.
El usuario construye una **expresión** con teclado físico o virtual; el frontend hace **una sola llamada HTTP** (`POST /v1/calculations`, modo `expression`) y el backend la compila al mismo DAG que ejecuta el scheduler. El frontend no ejecuta operadores ni compila DAGs.

```text
Keypad / teclado físico → expression string → validación cliente (liviana) → POST /v1/calculations {expression} → Display + errores → Historial en memoria
```

## Alcance v1 (decisiones acordadas)

| Aspecto                               | Decisión                                                                                                                                            |
| ------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| Modo API                              | Solo`expression` (sin builder visual de DAG en v1)                                                                                                 |
| Comunicación                         | `fetch` nativo + `useState` + `AbortController` (sin React Query / Zustand)                                                                    |
| Componentización                     | Por feature:`Display`, `Keypad`, `Key`, `HistoryPanel`, `ErrorToast` bajo `features/calculator/`                                         |
| Historial                             | Solo memoria (máx. 20),**sin** `localStorage` en v1                                                                                         |
| Fuente                                | `JetBrains Mono` vía CDN (Google Fonts `@import`)                                                                                               |
| Formato numérico                     | Locale`es-ES` (`1.234,56`; notación científica para magnitudes grandes)                                                                        |
| Keypad                                | Grid 4 cols × 6 filas; fila 1: `AC`,`+/-`,`(`,`)`; fila 2: `√`,`xʸ`,`⌫`,`%`; filas 3–5: dígitos con `/`,`×`,`−`; fila 6: `0`,`.`,`=`,`+` (sin celdas vacías ni spans) |
| Sin cuentas, sin persistencia, sin DB | Igual que backend v1                                                                                                                                 |

### Orden del keypad (definitivo)

```text
Row 1:  AC      +/-     (       )
Row 2:  √      xʸ      ⌫(del)   %
Row 3:  7       8       9       ÷  (/)
Row 4:  4       5       6       ×  (*)
Row 5:  1       2       3       −
Row 6:  0       .       =       +
```

Mapeo a expresión backend: `÷→/`, `×→*`, `−→-`, `√→sqrt(`, `xʸ→^`, `%→percent(,)` según helper (ver `expressionHelpers`).

## Contrato con el backend

Fuente de contrato: `docs/backend-calculator-plan.md` + `backend/api/openapi.yaml` (cuando exista; OpenAPI manda si hay conflicto de campos).

```ts
// features/calculator/api/types.ts
export interface CalculationRequest {
  expression: string;
  outputs?: string[]; // default: ["result"]
}
export interface OperationResult {
  id: string;
  value?: number;
  error?: { code: string; message: string; position?: number };
}
export interface CalculationResponse {
  results: Record<string, OperationResult>;
  outputs: string[];
  requestId: string;
  durationMs: number;
}
export type CalculationError =
  | { type: 'validation'; message: string; position?: number } // HTTP 400
  | { type: 'domain'; code: string; message: string }          // HTTP 422
  | { type: 'network'; message: string }
  | { type: 'timeout' }
  | { type: 'unknown'; message: string };
```

`POST {VITE_API_BASE_URL}/v1/calculations` con `{ expression, outputs: ["result"] }`.
Mapeo HTTP: `400→validation`, `422→domain`, `408/499→timeout`, fallo de red→`network`, resto→`unknown`.

## Tokens de diseño (acordados)

Paleta derivada del icono provisto (`DeleteIcon`, `#833AED`) + blanco.

| Token                                               | Valor                                                                             |
| --------------------------------------------------- | --------------------------------------------------------------------------------- |
| `accent` (símbolos operación/función)          | `#833AED`                                                                       |
| `accentLight` (fondo botones operación/función) | `rgba(131, 58, 237, 0.12)` + borde `rgba(131, 58, 237, 0.2)`                  |
| `accentDark` (fondo `=`)                        | `#6B21A8`, hover `#581C87`, texto `#FFF`                                    |
| `number bg / fg`                                  | `#FFF` / `#1A1A1A`                                                            |
| `module bg / page bg`                             | `#FFFFFF`                                                                       |
| `shadow module`                                   | `0 20px 40px -10px rgba(131,58,237,0.15), 0 8px 16px -4px rgba(131,58,237,0.1)` |
| `radius btn / display / module`                   | `12px` / `16px` / `20px`                                                    |
| `display`                                         | `JetBrains Mono 48px` (36px móvil), altura `100px`, alineado derecha         |
| `key`                                             | `JetBrains Mono 600 20px` (16px funciones), `64px` (56px móvil), gap `8px` |
| `module`                                          | Centrado,`max-width 360px`, padding `24px`                                    |
| `breakpoints`                                     | `480px` móvil / `768px` tablet / `1024px` desktop                          |

Implementación: custom properties en `src/index.css` + `shared/constants/designTokens.ts` tipado. Ver detalle en el [README raíz](../../README.md) (sección Frontend).

## Arquitectura por capas

| Capa                                    | Responsabilidad                                                                                                                                                                                                                           |
| --------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `features/calculator/api`             | Wrapper`fetch` (`calculationApi.ts`), tipos (`types.ts`). Único lugar que conoce la URL y el mapeo de errores HTTP                                                                                                                 |
| `features/calculator/hooks`           | `useCalculation` (fetch + loading/error/result/historial + abort), `useExpressionValidation` (regex cliente liviana), `useKeypad` (teclado físico + virtual → string), `useHistory` (memoria, máx. 20)                         |
| `features/calculator/components`      | Presentacionales:`CalculatorLayout`, `Display`, `Keypad`, `Key`, `HistoryPanel`, `ErrorToast`                                                                                                                                 |
| `features/calculator/utils`           | `keypadLayout.ts` (definición del grid), `expressionHelpers.ts` (símbolos UI → sintaxis backend), `formatResult.ts` (locale ES)                                                                                                  |
| `components/`                         | Base reutilizable:`Button`, `Card`, `Icon`, `Toast` (variantes por props, sin lógica de negocio)                                                                                                                                 |
| `shared/utils` + `shared/constants` | `cn`, `formatNumber`, `designTokens.ts`                                                                                                                                                                                             |
| `assets/icons`                        | Solo el icono provisto como componente React:`DeleteIcon` (convertido 1:1 de `DeletIcon.svg`, con `currentColor`). Las teclas de función (`√`, `xʸ`, `+/-`, `%`) usan etiqueta de texto, no se crean iconos nuevos en v1 |

Principio: añadir una tecla o un icono no obliga a tocar la capa API ni los hooks.

## Distribución de carpetas (objetivo)

```text
frontend/front-calculator/
├── public/fonts/               # solo si se abandona CDN
├── src/
│   ├── features/calculator/
│   │   ├── api/{calculationApi.ts,types.ts}
│   │   ├── components/{CalculatorLayout,Display,Keypad,Key,HistoryPanel,ErrorToast}.tsx
│   │   ├── hooks/{useCalculation,useExpressionValidation,useKeypad,useHistory}.ts
│   │   ├── utils/{keypadLayout,expressionHelpers,formatResult}.ts
│   │   └── index.ts
│   ├── components/{Button,Card,Icon,Toast}.tsx
│   ├── shared/
│   │   ├── utils/{cn,formatNumber}.ts
│   │   └── constants/designTokens.ts
│   ├── assets/icons/DeleteIcon.tsx   # único icono; resto de teclas con etiqueta de texto
│   ├── pages/App.tsx
│   ├── main.tsx
│   └── index.css
└── .env                        # VITE_API_BASE_URL=http://localhost:8080 (ver README raíz)
```

## Validación cliente (liviana, no sustituye al backend)

`useExpressionValidation`: expresión no vacía, caracteres permitidos `[\d\s+\-*/^().,a-zA-Z]`, paréntesis balanceados, funciones conocidas (`sqrt`, `percent`). El backend es la autoridad (errores 400/422 con `position` se muestran en `ErrorToast`).

## Manejo de errores y estados

- `isLoading`: skeleton en `Display`, `Key` deshabilitadas salvo `AC`.
- `ErrorToast`: mensaje + código backend; si hay `position`, resaltar carácter en la expresión.
- `AbortController`: cada `calculate()` cancela la petición anterior.
- `Display`: `es-ES` vía `Intl.NumberFormat`; `requestId` visible en modo detalle (debug).

## Pruebas y entrega

- **Unitarias (Vitest + Testing Library):** hooks (`useCalculation`, `useExpressionValidation`, `useKeypad`, `useHistory`), `formatNumber`/`formatResult`, `expressionHelpers`, `Key`/`Display`.
- **Integración (MSW):** `calculationApi` + `useCalculation` (success, 400, 422, timeout, network).
- **E2E manual v1:** escribir `sqrt(percent(200, 15)) + 4 ^ 2` → `=` → resultado → copiar → historial en memoria. (Playwright opcional, fuera de v1.)
- **CI frontend:** `lint`, `typecheck`, `test`, `build` verdes antes de cerrar fase.

## Plan por fases

| Fase                                  | Contenido                                                                                                                                                                                                 | Rama sugerida                   |
| ------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------- |
| **F1 · Fundaciones**           | Estructura, migración a TypeScript (`tsconfig`, `typecheck`, types en props), `designTokens.ts`, `index.css` (tokens + JetBrains Mono CDN), `components/` base, `DeleteIcon` como componente | `feat/frontend-foundations`   |
| **F2 · API + hooks**           | `calculationApi`, `types`, 4 hooks + tests (Vitest + MSW)                                                                                                                                             | `feat/frontend-api-hooks`     |
| **F3 · Componentes feature**   | `Display`, `Key`, `Keypad` (grid 4×6), `CalculatorLayout`, `HistoryPanel`, `ErrorToast`                                                                                                      | `feat/frontend-calculator-ui` |
| **F4 · Integración + polish** | `App.tsx` wiring, responsive, errores 400/422, copy, loading, historial memoria                                                                                                                         | `feat/frontend-integration`   |
| **F5 · Calidad**               | README, lint/typecheck/tests/build, PR                                                                                                                                                                    | `chore/frontend-quality`      |

## Mantenimiento de este documento

Si cambia el layout del keypad, los tokens, el contrato con el backend o las fases, **actualizar este archivo en el mismo PR** para que todos los agentes sigan una sola fuente de verdad.
