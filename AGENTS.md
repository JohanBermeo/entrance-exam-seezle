# Agent instructions · entrance-exam-seezle

This repository is a fullstack calculator project (Go backend + React frontend) for a technical assessment. **Before changing backend architecture, HTTP contracts, concurrency, or milestone scope, read the canonical plan.**

## Canonical documents (read first)

| Document | Purpose |
| --- | --- |
| [docs/backend-calculator-plan.md](docs/backend-calculator-plan.md) | Backend architecture, DAG execution model, REST surface, milestones, folder layout, numeric rules, and testing expectations |
| [docs/development-protocol.md](docs/development-protocol.md) | Branch naming, commits, milestone gates, human vs agent duties (no PR/merge unless explicitly asked) |

When the plan and existing code disagree, **do not silently drift**: align with the plan, ask the user, or update the plan document explicitly in the same change.

## Workflow rules (mandatory)

- **One milestone at a time.** Do not start the next milestone until the user **explicitly confirms** the current one is approved and asks to continue.
- **PR por hito (agente):** al cerrar el hito, abrir la PR hacia `main` con una descripción adecuada de lo realizado; **no mergear** salvo petición explícita del humano en esa conversación.
- **End of milestone:** run verification; update the matching README (`backend/README.md` or `frontend/front-calculator/README.md`, or both if applicable); open the PR with a proper description, give a short summary and branch name, then **stop and wait** for user approval before any new scope.

Full detail: [docs/development-protocol.md](docs/development-protocol.md).

## Repository layout

- `backend/` — Go HTTP API (`back-calculator` module). See [backend/README.md](backend/README.md) for local run and env vars.
- `frontend/front-calculator/` — React + Vite UI (consumes backend once OpenAPI is stable).
- `.github/workflows/` — CI (backend tests, lint, race detector as configured).

## Backend implementation status (high level)

Track detailed milestones in [docs/backend-calculator-plan.md](docs/backend-calculator-plan.md#plan-por-hitos).

| Milestone | Theme | Typical branch prefix |
| --- | --- | --- |
| 01 | Foundations (config, router, health, CI) | merged via `feature/backend-foundations` |
| 02 | Domain operators + typed errors | `feat/backend-domain` |
| 03 | DAG contract, refs, validation, cycles; expression mode (Pratt → DAG) | `feat/backend-dag` (or similar) |
| 04 | Concurrent scheduler (fan-out / fan-in) | `feat/backend-dag-scheduler` |
| 05 | OpenAPI, integration tests, metrics | `chore/openapi-contract`, etc. |

**v1 operators (API names):** `add`, `subtract`, `multiply`, `divide`, `power`, `sqrt`, `percent` (`percent(value, rate)` → `value × rate / 100`).

## Commands agents should run

Backend (from `backend/`):

```bash
go test ./...
go test -race ./...
go vet ./...
```

Frontend (from `frontend/front-calculator/`): follow that package’s README when touching UI.

## Scope guardrails for v1

- Input model: **explicit operation DAG** in JSON **or** **`expression`** string compiled to the same DAG via Pratt parser (`domain/expression`) — both are v1.
- **No** user accounts, persistent history, or database in v1.
- Numeric model: **`float64`**, reject NaN, infinity, division by zero, and negative square roots.
- Thin HTTP layer → `ExecuteCalculation` (application) → domain operators + graph engine.

## Cursor / local IDE metadata

The `.cursor/` directory is **gitignored** (local rules and canvases). The **versioned** source of truth for all agents is under `docs/` and this file.
