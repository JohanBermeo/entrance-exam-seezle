# Diálogo de sesión · Backend hitos 01–05 + Docker + inicio frontend

> **Nota:** este archivo se generó por petición explícita del usuario **sin aplicar
> el protocolo de desarrollo** (sin rama, commits ni PR). Es un registro del diálogo
> de la sesión de trabajo, no un documento canónico del plan.
> Fecha: 17–18 de septiembre de 2026. Idioma de la sesión: español.

---

## 1. Análisis de contexto del proyecto

**Usuario:** «Analiza los archivos AGENTS y toda la documentación relacionada para que
obtengas todo el contexto del estado del proyecto y lo que queda por hacer».

**Asistente:** leyó `AGENTS.md`, `docs/backend-calculator-plan.md`,
`docs/development-protocol.md`, exploró `backend/` y `frontend/front-calculator/`,
ejecutó `go test ./...` y revisó `git log`. Entregó el estado:

**Hitos completados:**

| Hito | Estado | Rama |
|------|--------|------|
| 01 · Fundaciones | ✅ hecho | `feature/backend-foundations` (mergeada) |
| 02 · Dominio | ✅ hecho | `feat/backend-domain` (mergeada) |

**Hitos pendientes:**

| Hito | Tema |
|------|------|
| 03 · DAG y expresión | Contrato de lote, refs, validación y ciclos; modo `expression` (Pratt → DAG) |
| 04 · Concurrencia | Scheduler fan-out/fan-in, pool acotado, cancelación y límites |
| 05 · Calidad | OpenAPI, integración HTTP, race detector y métricas |

Estructura existente entonces: `cmd/api`, `internal/{config,domain/{calculation/errors.go,operators},transport/http}`,
`Dockerfile`, `Makefile`, `go.mod` (módulo `back-calculator`, solo stdlib). Faltaban
`internal/application`, `internal/domain/expression`, `internal/engine`, `api/openapi.yaml`
y `tests/`. Preguntó si se procedía con el hito 03.

---

## 2. Hito 03 · DAG y expresión (`feat/backend-dag`)

**Usuario:** «Teniendo en cuenta el development-protocol procede con la implementación
del hito 3».

**Asistente:** creó la rama `feat/backend-dag` y el todo-list del hito. Implementación
por pasos:

1. **Dominio cálculo** (`internal/domain/calculation/`): `operation.go` (`Input`
   literal/`ref`, `Operation`, `Validate`), `graph.go` (`NewGraph`, detección de
   ciclos por DFS, `Dependencies`/`Dependents`, `ReadyOperations`), `result.go`
   (`Results`, `OutputValues`, `OutputErrors`).
2. **Operadores:** agregó `Arity()` a la interfaz `Operator` y a las 7
   implementaciones, más `ArityMap()` al `Registry` (única fuente de aridades para
   validador y compilador).
3. **Expresión** (`internal/domain/expression/`): `token.go`, `lexer.go`, `ast.go`,
   `parser.go` (Pratt), `compiler.go` (AST → mismo DAG, alias de salida, límites
   100 nodos / profundidad 50) + tests de lexer, parser (precedencia/asociatividad)
   y compilador.
4. **Aplicación:** `internal/application/execute_calculation.go` (modos DAG y
   `expression` mutuamente excluyentes; ejecución secuencial temporal).
5. **HTTP:** `calculation_handler.go` (`POST /v1/calculations`, mapeo 400/422/408/499/500)
   y registro de la ruta en `router.go`.
6. **README** del backend actualizado con uso, operadores, precedencia y errores.

**Bugs reales encontrados y corregidos durante el hito:**

- Bucle infinito en el lexer (faltaba `position++` en operadores de un carácter).
- El compilador modelaba literales como operación `add` de 1 entrada (aridad inválida)
  → rediseño a literales inline + alias identidad `x + 0` bajo el nombre de salida.
- `map[Expr]string` inválido (`CallExpr` contiene slices) → se eliminó el caché por AST.
- Precedencia: `^` debe ser asociativa a la derecha y de mayor precedencia que el
  unario (según el plan) → reorden de constantes + `precedence - 1` en rama derecha.
- `string(rune(i))` para índices → `strconv.Itoa(i)`.
- `ReadyOperations` devolvía punteros a la variable del `for` → índice explícito.
- **Adyacencia que borraba aristas:** `adj[op.ID] = nil` dentro del mismo bucle que
  agregaba aristas ocultaba ciclos de 2 nodos → se separó en dos bucles (detectado
  con test de depuración temporal, luego eliminado).

**Usuario (intermedio):** «No olvides la creación de los commits periódicos».

**Asistente:** commits pequeños y descriptivos (Conventional Commits) en la rama:

- `f880084` feat(backend): operator arity and registry arity map
- `17bff5d` feat(backend): expression lexer, Pratt parser and DAG compiler
- `914108e` feat(backend): validate DAG refs and cycles with typed errors
- `95eb62a` feat(backend): POST /v1/calculations endpoint
- `d5928dd` docs(backend): milestone 03 calculation endpoint

Verificación: `go vet` limpio, `go test ./...` todo OK (`-race` sin gcc local; lo cubre CI).
Resumen entregado y parada según protocolo.

---

## 3. Cambio al protocolo de desarrollo (directo en `main`)

**Usuario:** «Implementación aprobada, procede con el hito 4» y, en el mismo mensaje:
«Realiza un pequeño cambio en el protocolo de desarrollo directamente en main, agrega
que el agente debe crear la pr con una descripción adecuada de lo realizado en el hito
pero que no debe mergearla».

**Asistente:**

- Rama `feat/backend-dag` ya sincronizada en origin (sin `gh` ni token en el entorno,
  imposible abrir la PR programáticamente; se dejó lista para abrirla en la UI).
- Cambio en `main` (commit `f4a91e1`, pusheado directo por petición explícita):
  - `docs/development-protocol.md`: tabla de roles (el agente **sí** abre la PR con
    descripción; el merge sigue siendo del humano), pasos 6–7 (abrir PR y detenerse
    sin mergear) y línea de trunk-based.
  - `AGENTS.md`: bullet «PR por hito (agente)» en vez de «No PR by default».
  - `docs/backend-calculator-plan.md`: línea del esquema de ramificación coherente.

---

## 4. Hito 04 · Concurrencia (`feat/backend-dag-scheduler`)

**Contexto:** rama creada desde `feat/backend-dag` porque `main` aún no contenía el
hito 03 sin mergear (upstream desmarcado para no interferir).

**Asistente:** implementó `internal/engine/` según el layout del plan:

- `dependency_graph.go`: indegree + dependientes congelados desde el grafo validado.
- `executor.go`: interfaz `NodeExecutor` (seam para instrumentar en tests) +
  `RegistryExecutor` + `resolveInputs`.
- `worker_pool.go`: pool por semáforo con `Run(ctx, task)` cancelable.
- `scheduler.go`: fan-out acotado + recolector fan-in único, fail-fast con
  cancelación, cascada de fallos iterativa (sin recursión), `outputs` deterministas,
  constructor `NewScheduler(registry, maxWorkers)`.

Decisiones: ante el primer fallo se cancela el resto y se reporta el error raíz
(HTTp 422/408/499 según plan); `outputs` conserva el orden pedido.

**Tests del engine** (5, todos en verde al primer intento): 25 DAGs aleatorios vs
referencia secuencial, cota del pool (≤3 con paralelismo real), fail-fast (el
dependiente nunca ejecuta), cancelación sin deadlock, diamante determinista ×20.

**Cableado:** `Options{Registry, Scheduler, MaxNodes, MaxDepth}` +
`ExecuteWithOptions` (se eliminó `executeGraph` secuencial); config nueva
(`CALC_MAX_WORKERS=8`, `CALC_TIMEOUT=30s`, `CALC_MAX_BODY_BYTES=1MiB`,
`CALC_MAX_NODES=100`, `CALC_MAX_DEPTH=50`); handler con timeout y `MaxBytesReader`
configurables; `NewRouterWithOptions`; `main.go` con valores de config; tests de
application (workers 1/2/8, expresión, límites) y de config; README con sección de
concurrencia y tabla `CALC_*`.

**Commits:**

- `8220713` feat(backend): concurrent DAG scheduler with bounded pool
- `e4624b9` feat(backend): wire concurrent scheduler with configurable limits
- `c985813` docs(backend): milestone 04 scheduler and limits

Rama pusheada a origin. Resumen + enlace de PR sugerido, sin mergear. Aprobado por
el usuario.

---

## 5. Hito 05 · Calidad (`chore/openapi-contract`)

**Usuario:** «Implementación aprobada, procede con el hito 5».

**Contexto:** `main` ya traía mergeadas las PRs #3 (hito 03) y #4 (hito 04); la rama
se creó desde `main` con prefijo `chore/` según la tabla del AGENTS.

**Asistente:**

1. **`backend/api/openapi.yaml`** (OpenAPI 3.0.3): contrato fuente — `POST
   /v1/calculations` (modos DAG/expresión con ejemplos), `/healthz`, `/readyz`,
   `/metrics`, esquemas y códigos 400/422/408/499/500.
2. **Métricas** (`internal/platform/metrics`, stdlib): contadores por ruta, desglose
   por estado, latencia media e in-flight, servidas en `GET /metrics` + tests.
3. **Normalización de calidad:** `DomainError` con tags JSON en minúsculas
   (`code`/`operation`/`message`) según el contrato.
4. **Integración** (`backend/tests/integration/`): DAG y expresión válidos con campos
   del contrato, forma de errores, payload sobredimensionado (400), deadline
   determinista con ejecutor lento (408) y cliente cancelado (499), `/metrics`
   reflejando tráfico. (Un intento inicial con timeout de 1ns perdía la carrera
   contra el cálculo en memoria; se rediseñó determinista.)
5. **CI** (`.github/workflows/backend-ci.yml`): job `contract` que parsea el YAML
   (pyyaml) y exige rutas y esquemas; validado localmente con las mismas aserciones.
6. **README:** fuente OpenAPI, sección `/metrics`, tests de integración y nota de
   race en CI.

**Commits:**

- `a4d8358` feat(backend): OpenAPI contract and request metrics
- `a945996` test(backend): HTTP integration tests
- `7aa6e75` chore(backend): validate OpenAPI contract in CI and document metrics

Verificación final: `go build`, `go vet`, `go test ./...` (9 paquetes OK). Rama
pusheada. Resumen + datos de PR sugerida (`chore(backend): OpenAPI contract, metrics
and integration tests (milestone 05)`), sin mergear. Aprobado por el usuario
(mergeado luego como PR #5, commit `32af3cb`).

---

## 6. Cierre de fase backend y Docker

**Usuario:** «Implementación aprobada, queda algun hito faltante dentro de esta fase
del proyecto?».

**Asistente:** verificó el plan — los 5 hitos están implementados; alcance v1 completo
(sin cuentas, historial persistente ni DB). Solo restaba el merge manual de PRs.

**Usuario:** «Las prs se encuentran ya mergeadas, el docker está completamente
actualizado y funcional?».

**Asistente (modo lectura):** revisó `Dockerfile` (multi-stage
`golang:1.27.1-alpine` → `distroless/static-debian12:nonroot`, `CGO_ENABLED=0`,
`EXPOSE 8080`) y el árbol de `backend/`. Conclusión: funcional; solo faltaba
`.dockerignore` (opcional recomendado) + comandos de prueba local.

**Usuario:** «Crea el dockerignore con las especificaciones necesarias para el backend».

**Asistente:** rama `chore/backend-dockerignore`, archivo `backend/.dockerignore`
(control de versiones, Dockerfiles, tooling/docs, `.env*`, tests, artefactos Go,
IDE/OS; conserva `go.mod`, `cmd/`, `internal/`, `api/`), commit `7cf1494`, push.
Sin Docker local no se corrió el build; verificación lógica documentada. PR sugerida,
sin mergear.

---

## 7. Inicio de fase frontend (Figma)

**Usuario:** «vamos a dar el paso al análisis y planeación del desarrollo frontend,
para esto requiero que vayas a figma y analices el diseño de Calculadora web estética».

**Asistente:** pidió la URL de Figma + decisiones de alcance/stack (Tailwind,
TypeScript, orval, React Query, a11y) antes de redactar el plan.

**Usuario:** «utiliza el mcp de composio para acceder al diseño de figma».

**Asistente:** informó que Composio MCP no está disponible en el entorno (sin skills
ni herramientas MCP configuradas) y reiteró la necesidad de la URL, con tabla de
decisiones preliminares v1. Queda pendiente la respuesta del usuario.

---

## 8. Estado final relevante

- `main` (origin): hitos 01–05 mergeados (PRs #1–#5) + cambio de protocolo (`f4a91e1`).
- Ramas pusheadas pendientes de PR/merge: `chore/backend-dockerignore` (`7cf1494`).
- Regla vigente: el agente abre la PR del hito con descripción y **no mergea**
  (limitación práctica: sin `gh` CLI ni token en este entorno, deja rama pusheada +
  título/cuerpo sugeridos).
- Backend v1 completo y verificado (`go build`/`vet`/`test` en verde; `-race` en CI).
- Frontend: en análisis/planeación, bloqueado a la espera de la URL de Figma.
