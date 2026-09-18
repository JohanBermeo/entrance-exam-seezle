# entrance-exam-seezle

Calculadora fullstack para la prueba técnica de Seezle: **API en Go** con
ejecución concurrente sobre un **DAG de operaciones** + **UI en React** que
consume la API en modo expresión.

```text
Keypad / teclado físico → expression string → POST /v1/calculations
→ backend compila a DAG → fan-out / fan-in → resultado
```

El frontend nunca ejecuta operadores: construye la expresión y hace **una sola
llamada HTTP** por cálculo.

## Documentación

- **[AGENTS.md](AGENTS.md)** — punto de entrada para agentes de IA (OpenCode, Cursor, Copilot, etc.).
- **[docs/backend-calculator-plan.md](docs/backend-calculator-plan.md)** — arquitectura canónica del backend (DAG, REST, hitos).
- **[docs/frontend-calculator-plan.md](docs/frontend-calculator-plan.md)** — arquitectura canónica del frontend (modo expresión, keypad, tokens, fases F1–F5).
- **[docs/development-protocol.md](docs/development-protocol.md)** — ramas, commits y flujo de PRs.
- **[docs/uso-de-ia.md](docs/uso-de-ia.md)** — cómo se usó IA (generativa y agéntica) en el proyecto.
- **[backend/api/openapi.yaml](backend/api/openapi.yaml)** — contrato OpenAPI (fuente para React e integración).

## Backend — resumen

Servicio HTTP en Go (módulo `back-calculator`, Go 1.27.1+). Cada solicitud se
modela como un grafo de operaciones: los nodos independientes se ejecutan en
paralelo (fan-out acotado por `CALC_MAX_WORKERS`) y un recolector agrega
resultados (fan-in); ante el primer fallo se cancela el resto.

| Método | Ruta | Descripción |
| --- | --- | --- |
| `POST` | `/v1/calculations` | Calcula vía DAG explícito **o** `expression` compilada al mismo DAG |
| `GET` | `/healthz`, `/readyz` | Salud del proceso |
| `GET` | `/metrics` | Contadores por ruta, latencia media, solicitudes en curso |

Operadores (`op`): `add`, `subtract`, `multiply`, `divide`, `power`, `sqrt`,
`percent(value, rate)` (`value × rate / 100`). Modelo numérico `float64`:
rechaza NaN, infinitos, división por cero y raíz de negativos. Errores:
`400` (JSON/gramática), `422` (dominio), `408`/`499` (deadline/cancelación).

```bash
# DAG explícito
curl -X POST http://localhost:8080/v1/calculations \
  -H 'Content-Type: application/json' \
  -d '{"operations": [{"id": "sum", "op": "add", "inputs": [{"value": 12}, {"value": 8}]}], "outputs": ["sum"]}'

# Modo expresión
curl -X POST http://localhost:8080/v1/calculations \
  -H 'Content-Type: application/json' \
  -d '{"expression": "sqrt(percent(200, 15)) + 4 ^ 2", "outputs": ["result"]}'
# outputValues ≈ [21.477]
```

Configuración principal (`PORT` por defecto `8080`; ver tabla completa de
`CALC_*` y timeouts HTTP en el plan canónico). Verificación desde `backend/`:

```bash
go test ./...
go test -race ./...
go vet ./...
```

## Frontend — resumen

UI en React 19 + TypeScript + Vite, en modo **expresión** (`{ expression,
outputs: ["result"] }`), con `fetch` nativo + `AbortController` e historial
solo en memoria. Componentes por feature (`Display`, `Keypad`, `Key`,
`HistoryPanel`, `ErrorToast`) + base reutilizable en `src/components/`.

Keypad 6×4 (24 teclas, sin spans ni vacíos):

```text
AC  +/-  (   ) | √  xʸ  ⌫  % | 7 8 9 ÷ | 4 5 6 × | 1 2 3 − | 0  .  =  +
```

Tokens: acento `#833AED`, `=` en `#6B21A8`, fuente JetBrains Mono (CDN),
números en formato `es-ES`. Comportamiento: ante `sqrt(`, `percent(` o `(`
tras dígito/cierre se inserta `×` explícito; tras `=`, el renglón queda vacío
(una operación encadena, número/paréntesis/punto empiezan de cero); `+/-`
niega con paréntesis. Variable: `VITE_API_BASE_URL` (por defecto
`http://localhost:8080`, se incrusta en build). Verificación desde
`frontend/front-calculator/`:

```bash
pnpm install && pnpm lint && pnpm typecheck && pnpm test && pnpm build
```

## Ejecución del proyecto

### Opción recomendada: todo con Docker

```powershell
./run.ps1            # backend :8080 + frontend :3000
./run.ps1 -Action down
```

Opciones: `-Action up|down|restart|logs|build`, `-FrontendPort`,
`-BackendPort`, `-ApiUrl`, `-SkipBuild` (reutiliza imágenes). Requiere Docker
con el daemon en marcha.

### Manual (desarrollo)

```bash
# Terminal 1 — backend (http://localhost:8080)
go run ./cmd/api          # desde backend/

# Terminal 2 — frontend (http://localhost:5173)
pnpm install              # desde frontend/front-calculator/
pnpm dev
```

### Docker por servicio

```bash
docker build -t calc-api ./backend
docker run --rm -p 8080:8080 calc-api

docker build -t front-calculator ./frontend/front-calculator
docker run --rm -p 3000:80 front-calculator
```

### Smoke test

```bash
curl http://localhost:8080/healthz
curl -X POST http://localhost:8080/v1/calculations \
  -H 'Content-Type: application/json' \
  -d '{"expression": "sqrt(percent(200, 15)) + 4 ^ 2", "outputs": ["result"]}'
```

## Estructura del repositorio

```text
├── backend/                 # API Go (cmd/api, internal/{config,transport,application,domain,engine,platform})
├── backend/api/openapi.yaml # Contrato OpenAPI
├── frontend/front-calculator/ # UI React (src/{features/calculator,components,shared,assets,pages})
├── docs/                    # Planes canónicos, protocolo, uso de IA, diálogos de sesión
├── run.ps1                  # Levanta backend + frontend con Docker
└── .github/workflows/       # CI backend (tests, race, vet, OpenAPI) y frontend (lint, typecheck, test, build)
```
