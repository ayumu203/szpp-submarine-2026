# Repository Guidelines

## Project Structure & Module Organization
- `backend/`: Go backend (layered architecture scaffold).
- `backend/cmd/api/main.go`: API entrypoint.
- `backend/application`, `backend/domain`, `backend/infrastructure`, `backend/presentation`: layer directories for use cases, domain models, external adapters, and handlers.
- `frontend/`: static client files (`index.html`, `main.js`, `index.css`).
- `materials/`: team rules and process docs.
- `documents/`: design artifacts (for example Mermaid diagrams).

## Build, Test, and Development Commands
- `cd backend && go run ./cmd/api`: run backend locally.
- `cd backend && go test ./...`: run all Go tests.
- `cd backend && go fmt ./...`: format Go code before commit.
- `cd frontend && python3 -m http.server 8080`: serve frontend locally at `http://localhost:8080`.

If `main.go` is still a placeholder, implement the package/main function first, then run the commands above.

## Coding Style & Naming Conventions
- Follow Go defaults: `gofmt` formatting, tabs/standard spacing, idiomatic package names.
- JavaScript/HTML/CSS: keep files small and focused; split logic by feature when `main.js` grows.
- Naming (from team rules):
  - File/directory names: camelCase (example: `pokemonUseCase`).
  - Class/type files: UpperCamelCase when a file maps to a single class/type (example: `PokemonSpecies.go`).
- Keep architecture boundaries clear: avoid domain logic in presentation/infrastructure layers.

## Testing Guidelines
- Primary framework: Go `testing` package.
- Place tests next to implementation files using `*_test.go`.
- Name tests `Test<Behavior>` (example: `TestPlayerMoveWithinRange`).
- Add table-driven tests for game rules (attack range, movement limits, sunk-state handling).

## Commit & Pull Request Guidelines
- Branch naming: `prefix/verb-noun` (example: `feat/add-cpu-decision-method`).
- Commit format: `prefix: message` (examples: `feat: add move validation.`, `docs: update architecture notes.`).
- Allowed prefixes: `docs`, `feat`, `fix`, `refactor`, `file` (`chore` only for very minor work).
- Use `.github/pull_request_template.md`:
  - link the issue (`close #<id>`),
  - summarize implementation,
  - complete checklist,
  - assign Copilot review and then a human reviewer.
