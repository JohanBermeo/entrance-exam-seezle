# entrance-exam-seezle

Repository for the Seezle technical assessment: **Go calculator API** (DAG execution) + **React** frontend.

## Run everything (Docker)

```powershell
./run.ps1            # backend :8080 + frontend :3000
./run.ps1 -Action down
```

Options: `-Action up|down|restart|logs|build`, `-FrontendPort`, `-BackendPort`, `-ApiUrl`, `-SkipBuild` (reuse images). Requires Docker with a running daemon.

## Documentation

- **[AGENTS.md](AGENTS.md)** — entry point for AI coding agents (OpenCode, Cursor, Copilot, etc.)
- **[docs/backend-calculator-plan.md](docs/backend-calculator-plan.md)** — canonical backend architecture and milestones
- **[docs/frontend-calculator-plan.md](docs/frontend-calculator-plan.md)** — canonical frontend architecture (React expression mode, keypad, tokens, phases F1–F5)
- **[docs/development-protocol.md](docs/development-protocol.md)** — branches, commits, and PR workflow
- **[backend/README.md](backend/README.md)** — run and configure the API locally
- **[frontend/front-calculator/README.md](frontend/front-calculator/README.md)** — run the UI, env vars, API integration, keypad layout
