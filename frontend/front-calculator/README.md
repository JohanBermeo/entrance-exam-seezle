# front-calculator · React + TypeScript + Vite

UI de calculadora (modo **expresión**) que consume `POST /v1/calculations` del backend.
Plan canónico: [docs/frontend-calculator-plan.md](../../docs/frontend-calculator-plan.md). Contrato backend: [docs/backend-calculator-plan.md](../../docs/backend-calculator-plan.md).

## Requisitos

- Node 20+, `pnpm` 10+ (CI usa pnpm 11 + Node 22).
- Backend corriendo (por defecto `http://localhost:8080`).

## Variables de entorno

```env
# frontend/front-calculator/.env
VITE_API_BASE_URL=http://localhost:8080
```

## Scripts

```bash
pnpm install
pnpm dev        # desarrollo (http://localhost:5173)
pnpm build      # typecheck + build producción
pnpm preview    # previsualizar build
pnpm lint       # eslint
pnpm typecheck  # tsc --noEmit
```

```bash
pnpm test       # vitest run (jsdom + Testing Library + MSW)
```

## Docker

Imagen multi-stage (`node:22` → `nginx:1.27-alpine` con fallback SPA). Contexto: `frontend/front-calculator/`.

```bash
docker build -t front-calculator ./frontend/front-calculator
docker run --rm -p 8080:80 front-calculator
# http://localhost:8080
```

`VITE_API_BASE_URL` se incrusta en build (Vite). Para otro backend:

```bash
docker build --build-arg VITE_API_BASE_URL=http://api:8080 -t front-calculator ./frontend/front-calculator
```

El navegador llama a esa URL directamente: debe ser alcanzable desde el navegador, no solo entre contenedores.

## Arquitectura

```text
src/
├── features/calculator/
│   ├── api/         # calculationApi.ts (fetch), types.ts (request/response/error)
│   ├── components/  # CalculatorLayout, Display, Keypad, Key, HistoryPanel, ErrorToast
│   ├── hooks/       # useCalculation, useExpressionValidation, useKeypad, useHistory
│   └── utils/       # keypadLayout, expression (toggleLastNumberSign)
├── components/    # Button, Card, Icon, Toast (base reutilizable)
├── shared/
│   ├── utils/       # cn, formatNumber
│   └── constants/   # designTokens.ts
├── assets/icons/    # DeleteIcon.tsx (único icono; teclas de función con texto)
├── pages/App.tsx / main.tsx / index.css
```

Reglas: componentes presentacionales + hooks con lógica; solo `api/` conoce la URL y el mapeo de errores HTTP.

## Integración con la API

- `POST {VITE_API_BASE_URL}/v1/calculations` con `{ expression, outputs: ["result"] }`.
- Mapeo: `400→validation`, `422→domain`, `408/499→timeout`, fallo red→`network`.
- `useCalculation` usa `AbortController` (cada cálculo cancela el anterior) y guarda historial en memoria (máx. 20, sin `localStorage` en v1).
- Formato de números: locale `es-ES` (`Intl.NumberFormat('es-ES')`).

## Keypad (layout acordado)

```text
AC  +/-  (   ) | √  xʸ  ⌫  % | 7 8 9 ÷ | 4 5 6 × | 1 2 3 − | 0  .  =  +
```

- Números/`.` → fondo `#FFF`, texto `#1A1A1A`.
- Operaciones/funciones → fondo `rgba(131,58,237,0.12)`, símbolo `#833AED`.
- `=` → fondo `#6B21A8` (hover `#581C87`), texto blanco.
- Grid 4 cols, gap `8px`; botón `64px` (56px móvil); módulo `max-width 360px`, padding `24px`, radius `20px`, sombra morada.
- Fuente `JetBrains Mono` vía CDN (`@import` en `index.css`); display `48px` (36px móvil).

## Flujo de uso

1. Construir expresión con keypad o teclado físico (`useKeypad`).
2. Validación cliente liviana (caracteres, paréntesis, funciones `sqrt`/`percent`).
3. Pulsar `=` → `POST /v1/calculations` → resultado en `Display` o `ErrorToast` (con `position` si el backend la devuelve).
4. Historial en memoria: clic para reutilizar expresión.

## CI

`.github/workflows/frontend-ci.yml` corre en PR y push a `main` con cambios en `frontend/**`: `install --frozen-lockfile`, `lint`, `typecheck`, `test` y `build`.

## Verificación por fase

| Fase | Verificar con |
| --- | --- |
| F1 Fundaciones | `pnpm lint`, `pnpm typecheck`, `pnpm build`, check visual de `components/` |
| F2 API + hooks | `pnpm test` (Vitest + MSW: success, 400, 422, timeout, network) |
| F3 Componentes | tests de `Key`/`Display`/`Keypad`, grid 4×6 responsive |
| F4 Integración | E2E manual: `sqrt(percent(200, 15)) + 4 ^ 2` → `=` → resultado ES → copiar → historial |
| F5 Calidad | lint + typecheck + tests + build verdes, PR con descripción |

## Estado de implementación

- [x] **F1 Fundaciones** (rama `feat/frontend-foundations`): stack React + TypeScript (`tsconfig`, `typecheck`, types en props), tokens (`designTokens.ts` + `index.css` con JetBrains Mono CDN), `components/` (`Button`, `Card`, `Icon`, `Toast`), `shared/utils` (`cn`, `formatNumber` es-ES), `DeleteIcon` como componente (único icono; teclas de función con texto), shell visual en `pages/App.tsx`, fix del import en `main.tsx`. Verificado: `pnpm lint` ✅, `pnpm typecheck` ✅, `pnpm build` ✅.
- [x] **F2 API + hooks** (rama `feat/frontend-api-hooks`): `api/types` + `calculationApi` (mapeo 400/422/408-499/red/abort), hooks `useCalculation` (abort + historial memoria), `useExpressionValidation`, `useKeypad` (virtual + físico), `useHistory` (límite 20), `utils/expression` (`toggleLastNumberSign`), suite Vitest + MSW (34 tests). Verificado: `pnpm lint` ✅, `pnpm typecheck` ✅, `pnpm test` ✅, `pnpm build` ✅.
- [x] **F3 Componentes feature** (rama `feat/frontend-calculator-ui`): `utils/keypadLayout` (grid 4×6 definitivo, `=` solo fila 5, `0` span 2 + celda vacía), `Key` (variantes + icono delete + span), `Keypad` (22 teclas, con `disabled` solo AC habilitado), `Display` (expresión + resultado es-ES + loading), `CalculatorLayout`, `HistoryPanel` (colapsable, reutilizar, limpiar), `ErrorToast` (mapeo `CalculationError`→Toast) + estilos en `index.css`. Suite: 58/58 tests. Verificado: `pnpm lint` ✅, `pnpm typecheck` ✅, `pnpm test` ✅, `pnpm build` ✅.
- [x] **F4 Integración + polish** (rama `feat/frontend-integration`): `App.tsx` cablea keypad + validación cliente + `calculate` + historial (reutilizar/limpiar), `ErrorToast` (cliente 400-con-posición / servidor 422 / red), `Keypad` deshabilitado en carga (salvo AC), copiar resultado, meta `requestId · ms`, responsive 480px. Aditivo en hooks: `clearError` (`useCalculation`), `onEdit` (`useKeypad`). Test integración `App` con MSW (4 casos). E2E manual contra backend real: `sqrt(percent(200, 15)) + 4 ^ 2` → `21.477…`, `5/0` → `division_by_zero`, `1 +` → `invalid_input`. Suite: 63/63 tests. Verificado: `pnpm lint` ✅, `pnpm typecheck` ✅, `pnpm test` ✅, `pnpm build` ✅.
- [x] **F5 Calidad** (rama `chore/frontend-quality`): CI con paso `pnpm test`, `index.html` (`lang="es"`, título `Calculadora`), `.env` ignorado, README final. Verificado: `pnpm lint` ✅, `pnpm typecheck` ✅, `pnpm test` ✅, `pnpm build` ✅.
