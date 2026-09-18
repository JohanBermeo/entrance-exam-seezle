# Backend · Calculator API

Servicio HTTP en Go para la calculadora. Incluye configuración de ejecución, logging estructurado, endpoints de salud y el endpoint de cálculo `POST /v1/calculations` (hito 03: DAG explícito o `expression` compilada al mismo DAG; ejecución secuencial — el scheduler concurrente llega en el hito 04).

**Plan de arquitectura y hitos (canónico):** [../docs/backend-calculator-plan.md](../docs/backend-calculator-plan.md)  
**Instrucciones para agentes de IA:** [../AGENTS.md](../AGENTS.md)

## Requisitos

- Go 1.27.1 o posterior compatible con el módulo.

## Ejecutar localmente

```bash
go run ./cmd/api
```

El servicio escucha en `http://localhost:8080` de forma predeterminada.

```bash
curl http://localhost:8080/healthz
# {"status":"ok"}

curl http://localhost:8080/readyz
# {"status":"ready"}
```

## Calcular (hito 03)

Un solo endpoint acepta **uno** de los dos modos (mutuamente excluyentes):

```bash
# DAG explícito
curl -X POST http://localhost:8080/v1/calculations \
  -H 'Content-Type: application/json' \
  -d '{
    "operations": [
      {"id": "sum", "op": "add", "inputs": [{"value": 12}, {"value": 8}]},
      {"id": "pow", "op": "power", "inputs": [{"value": 4}, {"value": 2}]},
      {"id": "total", "op": "multiply", "inputs": [{"ref": "sum"}, {"ref": "pow"}]}
    ],
    "outputs": ["total"]
  }'
# {"requestId":"...","results":{...},"outputs":["total"],"outputValues":[320],"durationMs":...}

# Modo expresión (una sola salida en `outputs`)
curl -X POST http://localhost:8080/v1/calculations \
  -H 'Content-Type: application/json' \
  -d '{"expression": "sqrt(percent(200, 15)) + 4 ^ 2", "outputs": ["result"]}'
# outputValues ≈ [21.477]
```

Operadores: `add`, `subtract`, `multiply`, `divide`, `power`, `sqrt`, `percent(value, rate)`.
Precedencia de la expresión: paréntesis → funciones (`sqrt`, `percent`) → potencia `^` (derecha) → signo unario → `* /` → `+ -`. Límites: 100 nodos y profundidad 50 por expresión.

Errores: `400` (JSON inválido, modos mezclados, expresión malformada), `422` (operación desconocida, referencia ausente, ciclo, división por cero, raíz negativa, resultado no finito, aridad), `408`/`499` (deadline/cancelación).

## Configuración

| Variable | Valor predeterminado | Descripción |
| --- | --- | --- |
| `PORT` | `8080` | Puerto TCP del servidor. |
| `APP_ENV` | `development` | Entorno declarado en los logs. |
| `LOG_LEVEL` | `info` | Nivel: `debug`, `info`, `warn` o `error`. |
| `HTTP_READ_HEADER_TIMEOUT` | `5s` | Tiempo máximo para leer cabeceras. |
| `HTTP_READ_TIMEOUT` | `15s` | Tiempo máximo para leer la solicitud. |
| `HTTP_WRITE_TIMEOUT` | `15s` | Tiempo máximo para escribir la respuesta. |
| `HTTP_IDLE_TIMEOUT` | `60s` | Tiempo máximo de una conexión inactiva. |
| `SHUTDOWN_TIMEOUT` | `10s` | Límite para apagado ordenado. |

## Verificación

```bash
go test ./...
go test -race ./...
go vet ./...
```
