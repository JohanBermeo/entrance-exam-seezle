# Plan de backend · Calculadora REST

**Fuente canónica versionada** del plan de arquitectura (equivalente al canvas de diseño en Cursor). Cualquier agente de IA debe tratar este documento como referencia principal para el backend.

Stack: **Go**, **REST**, ejecución **fan-out / fan-in** sobre un **DAG** de operaciones.

## Visión

API pequeña y escalable: cada solicitud se modela como un grafo de operaciones. Los nodos independientes pueden ejecutarse en paralelo; los nodos con dependencias esperan sus referencias. El cliente hace **una sola llamada HTTP**.

### Decisión que habilita la concurrencia

La entrada es un **DAG**, no una lista lineal. Si `sum` y `power` no dependen entre sí, se calculan en paralelo. Una operación `total` que referencia ambas se ejecuta cuando terminan. La concurrencia **no** cambia el orden ni la semántica que ve el cliente.

## Flujo de ejecución

```text
HTTP → Validar (JSON) → [DAG explícito | expression → compilar a DAG] → Plan (nodos listos) → Fan-out → Fan-in → JSON respuesta
```

Entrada HTTP → validación → DAG (directo o compilado desde expresión) → ejecución concurrente limitada → agregación de resultados → respuesta.

## Alcance v1

| Área | Incluido |
| --- | --- |
| Básicas | Suma, resta, multiplicación, división |
| Complejas | Potencia, raíz cuadrada, porcentaje |
| Porcentaje | `percent(value, rate)` devuelve `value × rate / 100` |
| Modelo | Resultados por ID, referencias entre nodos, múltiples salidas |
| Entrada | DAG explícito **o** modo `expression` (lexer + parser Pratt → AST → mismo DAG que ejecuta el scheduler) |

**Excluido de v1:** cuentas de usuario e historial persistente.

### Nombres de operador en la API

| `op` | Aridad | Notas |
| --- | --- | --- |
| `add` | 2 | |
| `subtract` | 2 | |
| `multiply` | 2 | |
| `divide` | 2 | Rechazar divisor cero |
| `power` | 2 | Validar resultado finito |
| `sqrt` | 1 | Rechazar operandos negativos |
| `percent` | 2 | `value × rate / 100` |

## Contrato de cálculo

`POST /v1/calculations` acepta **uno** de los modos de entrada (mutuamente excluyentes en la misma solicitud):

1. **DAG explícito** — campo `operations` (+ `outputs`).
2. **Expresión** — campo `expression` (string); el módulo `domain/expression` la compila al mismo grafo que el modo DAG antes del scheduler.

### Modo DAG (ejemplo)

Solicitud representativa (`total` espera a `sum` y `pow`):

```json
{
  "operations": [
    {"id": "sum", "op": "add", "inputs": [12, 8]},
    {"id": "pow", "op": "power", "inputs": [4, 2]},
    {"id": "total", "op": "multiply", "inputs": [{"ref": "sum"}, {"ref": "pow"}]}
  ],
  "outputs": ["total"]
}
```

**Respuesta (objetivo v1):** valores por ID, lista `outputs`, errores por operación si aplica, duración total y `requestId`.

Referencias: cada input es un número literal o `{"ref": "<operation-id>"}`.

### Modo expression (ejemplo)

```json
{
  "expression": "sqrt(percent(200, 15)) + 4 ^ 2",
  "outputs": ["result"]
}
```

El compilador materializa nodos internos con IDs estables, aplica las mismas reglas de validación que el DAG explícito y reutiliza operadores y scheduler sin rama de ejecución paralela distinta.

## Superficie REST

| Método | Ruta | Descripción |
| --- | --- | --- |
| `POST` | `/v1/calculations` | Ejecuta cálculo vía DAG explícito o `expression` compilada a DAG |
| `GET` | `/healthz` | Salud del proceso (balanceadores) |
| `GET` | `/readyz` | Listo para recibir tráfico |

Versionado por URL. **OpenAPI** (`backend/api/openapi.yaml` cuando exista) es la fuente de contrato para React e integración.

## Arquitectura por capas

| Capa | Responsabilidad |
| --- | --- |
| Transporte | Router (stdlib o Chi); handlers delgados; recovery; request ID; límites de body; CORS; logging estructurado |
| Aplicación | `ExecuteCalculation`: valida, plan de dependencias, scheduler, mapeo a respuesta HTTP |
| Dominio | Operadores puros + registro `Operator`; grafo/operación/resultado; **expression** (token, lexer, AST, parser Pratt, compilador a DAG) |
| Orquestación | Scheduler DAG: cola de nodos listos, fan-in con un recolector |
| Plataforma | Config e instrumentación inyectables; sin DB en v1 |

Principio: añadir un operador no debe obligar a tocar transporte ni concurrencia.

## Distribución de carpetas (objetivo)

```text
backend/
├── cmd/api/main.go
├── internal/
│   ├── config/config.go
│   ├── transport/http/{router,calculation_handler,request,response}.go
│   ├── application/execute_calculation.go
│   ├── domain/calculation/{operation,graph,result,errors}.go
│   ├── domain/operators/{operator,registry,addition,subtraction,multiplication,division,power,square_root,percentage}.go
│   ├── domain/expression/{token,lexer,ast,parser,compiler}.go
│   ├── engine/{scheduler,worker_pool,dependency_graph,executor}.go
│   └── platform/{logger,metrics}/
├── api/openapi.yaml
├── tests/{unit,integration,testdata}/
└── {Dockerfile,Makefile,go.mod,go.sum,.golangci.yml}
```

La estructura actual puede estar parcialmente implementada según el hito; no eliminar capas planificadas sin actualizar este documento.

## Diseño fan-out / fan-in

1. **Validar primero:** IDs únicos, operador conocido, aridad correcta, referencias existentes, sin ciclos.
2. **Fan-out acotado:** pool configurable limita goroutines por solicitud; solo entran nodos con dependencias satisfechas.
3. **Fan-in con un propietario:** un recolector actualiza el DAG, libera dependientes y materializa resultados; workers no comparten mapas mutables.
4. **Cancelación:** `context.Context`, deadline por llamada, cancelación temprana si falla una operación requerida.
5. **Determinismo:** la respuesta respeta el orden declarado en `outputs` aunque la ejecución interna sea paralela.

## Modelo numérico inicial

- Tipo: **`float64`** documentado.
- Rechazar: NaN, infinito en entradas o resultados, división por cero, raíz de negativos.
- Códigos de dominio (implementación en progreso): incluyen p. ej. `unknown_operation`, `invalid_input`, `invalid_arity`, `division_by_zero`, `negative_square_root`, `non_finite_number`.
- Si más adelante se exige precisión decimal financiera: evolucionar a decimales serializados como strings **sin mezclar modelos en silencio**.

## Errores HTTP

| Código | Casos |
| --- | --- |
| `400` | JSON inválido, tipo no permitido, argumento faltante, expresión malformada (posición/token) |
| `422` | Operación desconocida, referencia ausente, ciclo, división por cero, raíz inválida, reglas de dominio, límites de expresión superados |
| `408` / `499` | Deadline agotado o cliente canceló |
| `500` | Fallo inesperado con `requestId`, sin filtrar internals |

Configurar límites: tamaño de payload, número de nodos, profundidad del DAG, workers, timeouts. Observabilidad: latencia, errores por operación, tamaño del grafo, operaciones activas.

## Parser Pratt (v1)

Módulo obligatorio en v1 bajo `internal/domain/expression/`. Pipeline: **lexer** (texto → tokens) → **parser Pratt** (precedencia, asociatividad, paréntesis) → **AST** → **compilador** a nodos del mismo DAG que consume el scheduler.

**Precedencia:** paréntesis → funciones (`sqrt`, `percent`) → potencia `^` (asociativa a la derecha) → signo unario → `* /` → `+ -`.

**Semántica:** `percent(value, rate)` = `value × rate / 100`; las funciones y operadores del AST mapean al registro de operadores existente.

**Errores del parser (v1):** reportar posición, token recibido y tokens esperados cuando aplique; rechazar expresiones que superen longitud máxima o profundidad máxima del AST (protección ante cargas abusivas).

**Pruebas (v1):** unitarias de lexer, parser (precedencia/asociatividad), compilador a DAG y casos de error con posición.

## Plan por hitos

| Hito | Contenido |
| --- | --- |
| **01 · Fundaciones** | Módulo Go, configuración, router, health checks, logging y CI |
| **02 · Dominio** | Operadores puros y errores tipados para las siete operaciones |
| **03 · DAG y expresión** | Contrato de lote, referencias, validación y ciclos; modo `expression` con lexer, parser Pratt y compilación al DAG |
| **04 · Concurrencia** | Scheduler fan-out/fan-in, pool acotado, cancelación y límites |
| **05 · Calidad** | OpenAPI, integración HTTP, race detector y métricas |

## Pruebas y entrega

- **Unitarias:** operadores, reglas de dominio, lexer/parser/compilador de expresiones, validación del grafo, scheduler.
- **Concurrencia:** DAGs aleatorios acíclicos, límites de pool, cancelación, sin deadlocks, `go test -race`.
- **Integración:** contrato OpenAPI, endpoint completo, payloads inválidos, timeout.
- **CI:** formato, lint, tests, race detector, validación del contrato OpenAPI cuando exista.

## Esquema de ramificación

Trunk en `main`. Ramas cortas `feat/`, `fix/`, `chore/`, `docs/`. PR pequeño, revisión y CI verde (PR y merge los gestiona el humano; el agente espera aprobación entre hitos). Detalle en [development-protocol.md](./development-protocol.md).

```text
main ──●────●────●────●── producción
  └─ feat/backend-dag-scheduler ──●─┘
  └─ fix/division-by-zero ──●─┘
  └─ chore/openapi-contract ──●─┘
```

## Mantenimiento de este documento

Si cambias alcance, contrato HTTP o hitos en código o en un canvas de diseño, **actualiza este archivo en el mismo PR** para que Cursor, OpenCode, Copilot y otros agentes sigan una sola fuente de verdad.
