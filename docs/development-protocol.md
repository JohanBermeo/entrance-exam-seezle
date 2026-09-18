# Protocolo de desarrollo

Documento versionado para humanos y agentes de IA. Complementa [backend-calculator-plan.md](./backend-calculator-plan.md).

## Roles: humano vs agente

| Responsabilidad | Humano (Johan) | Agente de IA |
| --- | --- | --- |
| Aprobar cierre de un hito | Sí | No |
| Entregar mensaje-PR del hito (título + descripción lista para pegar) | No | **Sí**, al cerrar el hito (ver paso 6) |
| Revisar, aprobar y hacer merge del PR | Sí | **No** |
| Crear rama, implementar, commits locales | Puede | Sí (cuando se pida) |
| Ejecutar pruebas / verificación técnica | Puede | Sí, antes de dar el hito por implementado |
| Actualizar README del frente tocado | Puede | Sí, al cerrar el hito (ver paso 5) |

El agente **debe** publicar la rama (`git push`) y entregar el **mensaje-PR** del hito: título sugerido + descripción adecuada de lo realizado (qué incluye, cómo verificar, qué cambió en el contrato si aplica) + enlace directo de creación (`.../pull/new/<rama>`). La PR la crea el humano; el agente **no** debe crearla (ni con `gh`) ni mergear salvo que el humano lo pida **explícitamente** en ese momento.

## Flujo de trabajo

1. **Rama:** crear una rama nueva desde `main` con prefijo acorde al trabajo: `feat/`, `fix/`, `chore/` o `docs/`.
2. **Implementación:** seguir el plan aprobado en [backend-calculator-plan.md](./backend-calculator-plan.md), **solo el hito en curso**.
3. **Commits:** commits pequeños, periódicos y descriptivos (trazabilidad). Preferir [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:`, `test:`, `docs:`, `chore:`).
4. **Verificación:** ejecutar pruebas y comprobaciones necesarias antes de considerar implementado el hito (p. ej. `go test ./...` en `backend/`).
5. **README del frente:** actualizar el README correspondiente al área trabajada, en el mismo hito y antes de darlo por cerrado:
   - Backend → [backend/README.md](../backend/README.md) (cómo ejecutar, endpoints nuevos, variables, verificación).
   - Frontend → [frontend/front-calculator/README.md](../frontend/front-calculator/README.md) (scripts, integración con API, flujo de uso).
   Si el hito toca ambos frentes, actualizar ambos. Cambios solo de documentación global pueden ir en el [README raíz](../README.md).
6. **Entrega del hito (agente):** resumir qué quedó hecho, cómo verificarlo, qué README se actualizó y en qué rama está el trabajo; **entregar el mensaje-PR** (título + cuerpo listos para pegar y enlace de creación de la PR hacia `main`); **detenerse y esperar sin mergear**.
7. **Aprobación del hito (humano):** cuando el humano da por aprobado el hito, revisa la PR, decide cuándo mergea y si continúa al siguiente hito.
8. **Siguiente hito:** el agente **solo** inicia el hito siguiente tras **confirmación explícita** del humano (p. ej. “aprobado, sigue con el hito 03”). Sin esa confirmación, no avanzar de hito aunque el código esté listo.

## Trunk-based development

- `main` permanece integrable y protegido.
- Ramas cortas que nacen de `main` y vuelven con PR pequeño (el agente entrega el mensaje-PR; el humano crea la PR, revisa y mergea).
- Etiquetas semánticas desde `main`. Usar `release/x.y` solo si hace falta estabilizar una versión mientras sigue el desarrollo en paralelo.

## Ejemplos de ramas

| Prefijo | Uso |
| --- | --- |
| `feat/backend-domain` | Hito dominio (operadores) |
| `feat/backend-dag-scheduler` | Scheduler concurrente |
| `fix/division-by-zero` | Corrección acotada |
| `chore/openapi-contract` | Contrato OpenAPI, tooling |
| `docs/backend-plan` | Solo documentación |

## Salida esperada de la planeación

Contrato OpenAPI acordado, esqueleto Go ejecutable y pruebas del motor de DAG **antes** de acoplar React. Frontend y backend pueden avanzar en paralelo sobre una interfaz estable.
