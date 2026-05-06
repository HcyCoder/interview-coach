# Interview Coach

## Tech Stack
- `frontend/`: Next.js 14 App Router, TypeScript, Tailwind CSS, shadcn/ui
- `backend/`: Go, Hertz, GORM, MySQL
- `ai-service/`: Python, FastAPI

## Commands
- Backend build/test: `cd backend && make build && make test`
- Frontend dev/build: `cd frontend && npm run dev`
- AI service dev: `cd ai-service && fastapi dev app/main.py`

## Conventions
- Keep services isolated; do not cross-import runtime code between `frontend`, `backend`, and `ai-service`.
- Use source-driven development for framework-specific code: verify against official docs before coding.
- Prefer small, testable increments and commit each service bootstrap independently.
- Keep API contracts explicit and stable; use `contracts/` or documented request/response shapes when adding endpoints.

## Boundaries
- Do not commit secrets or local environment files.
- Do not add new dependencies unless they are needed for the current service slice.
- Keep generated files and build artifacts out of version control.
