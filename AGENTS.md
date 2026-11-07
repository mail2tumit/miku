# AGENTS.md

A small news site: a Go/Gin API backed by PostgreSQL, and a React + Vite frontend that reads from it. There is one commit ("Initial commit") and no CI, tests, or migrations yet.

## Layout

```
backend/main.go                  entire API: DB connection, CORS, handlers
frontend/src/App.jsx             page: fetches news, renders loading/error/grid
frontend/src/services/newsService.js  getNews() -> GET /api/news
frontend/src/App.css             app styles; index.css holds template defaults
docker-compose.yml               postgres:18 only
```

## Stack

- Backend: Go 1.27.3, module `example.com/news-website/backend`; Gin, gin-contrib/cors, pgx v5 (via `stdlib` driver), godotenv.
- Frontend: React 19, Vite 8, Oxlint. No TypeScript, no router, no state library.
- Database: PostgreSQL 18 from `docker-compose.yml` (user `newsuser`, password `news_password`, db `news_db`, host port 5432).

## Running

Backend must be started from `backend/` because `godotenv.Load()` reads `.env` from the working directory:

```bash
docker compose up -d          # start Postgres
cd backend && go run .        # serves on :8080
cd frontend && npm install && npm run dev   # Vite dev server on :5173
```

`DB_URL` must be set (`.env` in `backend/`, git-ignored) or the process exits. The expected DSN format is not documented in the repo. The `news` table (`id`, `title`, `content`, `image_url`, `created_at`) is never created by the code — there are no schema or migration files, so it must exist before `/api/news` works.

## Commands

| Purpose | Command | Notes |
| --- | --- | --- |
| Backend build | `cd backend && go build ./...` | verified clean |
| Backend vet | `cd backend && go vet ./...` | verified clean; no test files exist |
| Frontend build | `cd frontend && npm run build` | |
| Frontend lint | `cd frontend && npm run lint` | Oxlint, config in `.oxlintrc.json` |
| Frontend preview | `cd frontend && npm run preview` | |

There are no tests and no test runner configured in either package.

## API surface

- `GET /api/health` — liveness message.
- `GET /api/news` — all rows ordered by `created_at DESC`; returns `[]` rather than `null` when empty.

Keep routes under the `/api` prefix and answer with `gin.H{"error": ...}` plus a 5xx status on failure, matching `getNews`.

## Conventions

- Go: gofmt/tabs, package-level doc comments on exported types and functions, handlers written as closures returning `gin.HandlerFunc` and closing over `*sql.DB`, SQL as raw backtick strings. Errors are logged with `log.Println` and returned to the client as a generic message — never leak driver errors.
- CORS is allowlisted to `http://localhost:5173` for `GET` only. A new frontend origin or method needs an explicit change in `main.go`.
- The frontend API base URL is hardcoded in `newsService.js` (`http://localhost:8080/api`); no env vars are used.
- JS/JSX style is inconsistent: `App.jsx` and `newsService.js` use semicolons and double quotes; `main.jsx` and `vite.config.js` keep the Vite template's no-semicolon, single-quote style. Match the file you are editing.
- CSS indentation differs too: `App.css` uses 4 spaces, `index.css` uses 2.

## Not verified / gaps

- `.env` values, any real database contents, and the correct `DB_URL` format are unknown; no `.env.example` exists.
- Frontend build and lint were not executed (no `node_modules`; requires `npm install`).
- `frontend/README.md` is the unmodified Vite template and describes nothing project-specific.
