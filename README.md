# Weekly Wrapped — Backend

The REST API powering **Weekly Wrapped**, a "Spotify Wrapped"-style tracker that turns logged daily activities into a weekly analytics dashboard and a shareable recap card. Built with **Go**, **Gin**, and **PostgreSQL** (via `pgx`, no ORM).

Pairs with the [Weekly Wrapped frontend](https://github.com/admalfrizi/weekly-wrapped-fe) (Next.js), which talks to this API through a BFF proxy.

## Tech Stack

| Layer | Choice |
|---|---|
| Language / framework | Go 1.26, Gin |
| Database | PostgreSQL, accessed via `pgx/v5` + `pgxpool` (raw SQL, no ORM) |
| Migrations | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Auth | JWT (`golang-jwt/jwt/v5`, HS256) + `bcrypt` password hashing |
| CORS | `gin-contrib/cors` |
| Hot reload (dev) | [Air](https://github.com/air-verse/air) |
| Containerization | Docker (multi-stage: `builder` → `development` / `production`), Docker Compose for local dev |

## Architecture

The service follows a layered **Controller → Service → Repository** pattern with dependencies wired by hand in `cmd/main.go` (no DI framework/container):

```
internal/
├── config/       # Env-based config loading + pgxpool database initialization
├── controller/   # Gin handlers — request binding, calls service, shapes HTTP response
├── service/      # Business logic (aggregation, token generation, recap building, etc.)
├── repository/   # Data access — raw SQL via pgx, one repository per domain
├── query/        # Reusable raw SQL query strings used by repositories
├── model/        # Core domain structs (User, Activity, Category, WeeklyRecap)
├── dto/          # Request payload structs with Gin `binding` validation tags
├── response/     # Response DTOs + mappers from model → response, and the shared API envelope
├── middleware/   # JWT auth middleware
└── router/       # Route registration, grouped by domain

db/
├── migrations/   # SQL up/down migrations (golang-migrate)
└── seeders/      # Standalone Go program to seed default categories

cmd/
└── main.go       # Wires config → repositories → services → controllers → router, starts the server
```

Some notable patterns:

- **Raw SQL, no ORM.** Repositories issue SQL directly against a shared `pgxpool.Pool` (wrapped as a `DBTX` interface in `BaseRepository`), with reusable queries factored into `internal/query`.
- **Consistent response envelope.** Every endpoint returns `{ status, message, data?, meta?, errors? }` via helpers in `internal/response/api_response.go` (`Success`, `SuccessWithPagination`, `Error`), so pagination (`page`, `limit`, `total_items`, `total_pages`) is always shaped the same way.
- **JWT auth.** `JWTMiddleware` validates the `Authorization: Bearer <token>` header, parses the `sub` claim as the user ID, and stores it in the Gin context for handlers to read. Both access and refresh tokens are currently issued with the same expiry window — check `generateTokens` in `auth_service.go` if you need to change that.
- **CORS** is currently configured to allow only `http://localhost:3000` (see `internal/router/router.go`) — update this for staging/production origins.

## Data Model

Defined in `db/migrations/000001_init_schema.up.sql`:

- **users** — account/profile info, `password_hash` (bcrypt).
- **categories** — activity categories (seeded: Coding, Reading, Exercise, Gaming, Working).
- **activities** — a logged entry: `user_id`, `category_id`, `value`, optional `note`, `occurred_at`. Indexed on `(user_id, occurred_at)` and `category_id`.
- **weekly_recaps** — a frozen snapshot of a user's week: `week_start`/`week_end`, a unique `slug`, the full dashboard payload as `stats_snapshot JSONB`, and a generated `narrative`. Constrained to **one recap per user per `week_start`** (`unique_user_week`), so regenerating a recap for a week you've already generated overwrites it in place (upsert) rather than creating a duplicate.

## How the core features work

- **Weekly dashboard** (`dashboard_service.go`) — for a given week, it sums activity values per category, compares totals against the previous week to compute a trend %, builds a per-day chart series, computes each category's share of the week's total, and picks the top category to generate a short Indonesian-language insight sentence.
- **Recap generation** (`recap_service.go`) — reuses the dashboard aggregation above, serializes it into `stats_snapshot`, builds a slug as `wk<week_number>-<6-char hex>`, and wraps the insight into a short narrative. Recap lookup by slug (`GET /recaps/:slug`) is public — no auth required — so it can power a shareable page.

## API Reference

Base path: `/api/v1`. Routes marked 🔒 require a valid `Authorization: Bearer <accessToken>` header.

| Method | Path | Description |
|---|---|---|
| GET | `/ping` | Health check |
| POST | `/auth/register` | Create an account |
| POST | `/auth/login` | Log in, returns access + refresh tokens |
| POST | `/auth/refresh` | Exchange a refresh token for a new token pair |
| GET 🔒 | `/users/me` | Get the current user's profile |
| POST 🔒 | `/users/me` | Update the current user's profile |
| GET 🔒 | `/activity/` | List activities (paginated: `page`, `limit`) |
| POST 🔒 | `/activity/` | Create an activity |
| GET 🔒 | `/activity/categories` | List activity categories |
| GET 🔒 | `/activity/:id` | Get a single activity |
| PUT 🔒 | `/activity/:id` | Update an activity |
| DELETE 🔒 | `/activity/:id` | Delete an activity |
| GET 🔒 | `/dashboard/weekly` | Weekly dashboard aggregation (`?start_date=YYYY-MM-DD`, defaults to current week) |
| POST 🔒 | `/recaps/generate` | Generate (or regenerate) the recap for a week (`?start_date=YYYY-MM-DD`) |
| GET | `/recaps/:slug` | Get a recap by its public slug (no auth) |

## Getting Started

### Prerequisites

- Go 1.26+
- PostgreSQL (or use the provided Docker Compose setup)
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI if running migrations outside Docker

### Environment variables

No `.env.example` is currently committed — create a `.env` file in the project root with:

| Variable | Used by | Default |
|---|---|---|
| `DB_HOST` | app | `localhost` |
| `DB_PORT` | app | `5432` |
| `DB_USER` | app | `root` |
| `DB_PASSWORD` | app | `lokal26` |
| `DB_NAME` | app | `weekly-wrapped-db` |
| `DB_SSLMODE` | app | `disable` |
| `JWT_SECRET` | app | **required, no default** — the app fails to start without it |
| `DB_URL` | category seeder only (`db/seeders/categories.go`) | **required for seeding** — a full `postgres://` connection string, separate from the `DB_*` vars above |

### Run locally

```bash
go run ./cmd
# or, with hot reload (uses the committed .air.toml):
air
```

The server starts on `:8080`.

### Run migrations

```bash
migrate -path=db/migrations -database "postgres://root:lokal26@localhost:5432/weekly-wrapped-db?sslmode=disable" up
```

### Seed categories

```bash
# requires DB_URL to be set
go run ./db/seeders/categories.go
```

### Docker Compose (recommended for local dev)

```bash
docker compose up
```

This spins up three services:
- `db-dev` — Postgres, with a healthcheck
- `migrator-dev` — runs `db/migrations` against `db-dev` once, then exits
- `backend-dev` — builds the `development` Dockerfile stage and runs the app under Air for hot reload, mounted against your local source

The API is then available at `http://localhost:8080`.

### Production build

The `Dockerfile`'s `production` stage compiles both the server and the seeder binary, then runs the seeder before starting the server:

```bash
docker build --target production -t weekly-wrapped-be .
docker run -p 8080:8080 --env-file .env weekly-wrapped-be
```

## License

No license file is currently published in this repository — all rights reserved by the author unless a license is added.
