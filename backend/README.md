# Backend · Calculator API

Servicio HTTP en Go para la calculadora. Esta primera fase incluye la configuración de ejecución, logging estructurado y endpoints de salud.

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
