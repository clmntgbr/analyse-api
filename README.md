# Go API Template

A production-ready Go API starter built with **Clean Architecture**, **Fiber**, **PostgreSQL**, and **Clerk** authentication.

It provides a solid foundation for building HTTP APIs with:

- JWT authentication via Clerk JWKS
- Clerk webhook synchronization (Svix signature verification)
- User provisioning on first authenticated request
- Docker-based development with hot reload (Air)
- CLI for database migrations
- CORS, rate limiting, and security middleware

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.25 |
| HTTP framework | [Fiber v3](https://gofiber.io/) |
| ORM | [GORM](https://gorm.io/) |
| Database | PostgreSQL 16 |
| Authentication | [Clerk](https://clerk.com/) + JWT (JWKS) |
| Webhooks | [Svix](https://www.svix.com/) signature verification |
| CLI | [Cobra](https://github.com/spf13/cobra) |
| Dev tooling | Docker, Air, golangci-lint |

## Architecture

The project follows **Clean Architecture** with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────┐
│  cmd/          Entry points (api, cli)                  │
├─────────────────────────────────────────────────────────┤
│  handler/      HTTP handlers + middleware               │
│  presenter/    API response formatting                  │
├─────────────────────────────────────────────────────────┤
│  usecase/      Business logic (application layer)       │
├─────────────────────────────────────────────────────────┤
│  domain/       Entities + repository interfaces         │
├─────────────────────────────────────────────────────────┤
│  repository/   GORM implementations                     │
│  infrastructure/  External services (Clerk, config…)    │
└─────────────────────────────────────────────────────────┘
```

### Request flow (authenticated API)

```
Client
  │
  ▼
Fiber middleware (helmet, cors, rate limit, logger)
  │
  ▼
AuthenticateMiddleware
  ├─ Extract Bearer token
  ├─ ValidateTokenUseCase (JWKS from Clerk)
  ├─ If user missing locally → FetchUserUseCase (Clerk API) → CreateUserUseCase
  └─ Set user in request context
  │
  ▼
Handler → Presenter → JSON response
```

### Webhook flow (Clerk)

```
Clerk Dashboard
  │
  ▼
POST /webhook/clerk
  │
  ▼
ClerkMiddleware (Svix signature verification)
  │
  ▼
ClerkHandler
  ├─ user.created  → CreateUserUseCase
  ├─ user.updated  → UpdateUserUseCase
  └─ user.deleted  → DeleteUserByClerkIDUseCase
```

## Project Structure

```
.
├── cmd/
│   ├── api/                 # HTTP server entry point
│   │   ├── main.go
│   │   ├── routes.go
│   │   └── wire/            # Dependency injection container
│   └── cli/                 # CLI commands (migrations, etc.)
│       ├── main.go
│       └── command/
├── domain/
│   ├── entity/              # Database models
│   └── repository/          # Repository interfaces
├── usecase/
│   ├── auth/                # Token validation
│   ├── clerk/               # Clerk API integration
│   └── user/                # User CRUD
├── handler/
│   ├── context/             # Request context helpers
│   └── middleware/          # Auth + webhook middleware
├── infrastructure/
│   ├── auth/                # JWT DTOs
│   ├── clerk/               # Clerk DTOs + JWKS provider
│   ├── config/              # Environment + database config
│   └── paginate/            # Pagination helpers
├── repository/
│   └── gorm/                # GORM repository implementations
├── presenter/               # API response models
├── compose.yaml             # Docker Compose (dev)
├── Dockerfile               # Multi-stage build (dev + prod)
├── Makefile                 # Dev shortcuts
└── .air.toml                # Hot reload config
```

## Applications

### API server (`cmd/api`)

Main HTTP server. Started automatically in Docker via Air hot reload.

```bash
go run ./cmd/api
```

### CLI (`cmd/cli`)

Command-line tool for operational tasks.

| Command | Description |
|---|---|
| `migrate` | Run GORM auto-migrations on the database |

```bash
# Via Docker (recommended)
make migrate

# Or locally
go run ./cmd/cli migrate
```

## API Endpoints

### Health checks

| Method | Path | Auth |
|---|---|---|
| `GET` | `/livez` | No |
| `GET` | `/readyz` | No |
| `GET` | `/startupz` | No |

### Protected routes

All `/api/*` routes require a valid Clerk JWT in the `Authorization` header.

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/users/me` | Returns the authenticated user profile |

**Example request:**

```bash
curl -H "Authorization: Bearer <clerk_jwt>" http://localhost:4000/api/users/me
```

**Example response:**

```json
{
  "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "clerkId": "user_2abc123",
  "firstName": "John",
  "lastName": "Doe",
  "createdAt": "2026-06-09T18:00:00Z",
  "updatedAt": "2026-06-09T18:00:00Z"
}
```

### Webhooks

| Method | Path | Auth |
|---|---|---|
| `POST` | `/webhook/clerk` | Svix signature |

Supported Clerk events:

- `user.created` — creates a local user
- `user.updated` — updates first name, last name, banned status
- `user.deleted` — removes the local user

## Clerk Setup

### 1. Create a Clerk application

1. Go to [clerk.com](https://clerk.com/) and create an application.
2. Note your **Frontend API URL** (e.g. `https://your-app.clerk.accounts.dev`).
3. Generate a **Secret Key** from the API Keys section.

### 2. Configure environment variables

```env
CLERK_SECRET_KEY=sk_test_...
CLERK_FRONTEND_API=https://your-app.clerk.accounts.dev
CLERK_WEBHOOK_SECRET=whsec_...
```

### 3. Configure the webhook endpoint

1. In the Clerk Dashboard, go to **Webhooks** → **Add Endpoint**.
2. Set the URL to your public API endpoint:
   ```
   https://your-domain.com/webhook/clerk
   ```
3. Subscribe to events: `user.created`, `user.updated`, `user.deleted`.
4. Copy the **Signing Secret** into `CLERK_WEBHOOK_SECRET`.

#### Local development with ngrok

The `compose.yaml` includes an ngrok service to expose your local API publicly for webhook testing.

1. Get an auth token at [ngrok.com](https://ngrok.com/).
2. Set `NGROK_AUTHTOKEN` in your `.env`.
3. Update the ngrok `command` in `compose.yaml` with your reserved domain (or remove `--url` for a random tunnel).
4. Use the ngrok URL as your Clerk webhook endpoint.
5. Monitor requests at `http://localhost:4040` (ngrok inspector).

### 4. How authentication works

1. The client sends `Authorization: Bearer <jwt>`.
2. The API validates the JWT against Clerk's JWKS endpoint (`{CLERK_FRONTEND_API}/.well-known/jwks.json`).
3. The token issuer must match `CLERK_FRONTEND_API`.
4. If the user does not exist locally, the API fetches their profile from Clerk and creates a local record.
5. Banned users are rejected with `401`.

## Environment Variables

Copy the template and fill in your values:

```bash
cp .env.dist .env
```

### Required variables

| Variable | Description | Example |
|---|---|---|
| `PORT` | API port inside the container | `3000` |
| `GO_ENV` | Environment (`development` / `production`) | `development` |
| `DATABASE_URL` | PostgreSQL connection string | `postgresql://user:pass@database:5432/db?sslmode=disable` |
| `CLERK_SECRET_KEY` | Clerk backend secret key | `sk_test_...` |
| `CLERK_FRONTEND_API` | Clerk Frontend API URL (issuer + JWKS) | `https://app.clerk.accounts.dev` |
| `CLERK_WEBHOOK_SECRET` | Clerk webhook signing secret | `whsec_...` |
| `CORS_ALLOWED_ORIGINS` | Comma-separated allowed origins | `http://localhost:3000` |
| `CORS_ALLOW_CREDENTIALS` | Allow credentials in CORS | `true` |
| `CORS_ALLOW_METHODS` | Allowed HTTP methods | `GET,POST,PUT,DELETE,OPTIONS` |
| `CORS_ALLOW_HEADERS` | Allowed request headers | `Origin,Content-Type,Accept,Authorization` |
| `CORS_MAX_AGE` | CORS preflight cache (seconds) | `86400` |
| `RATE_LIMIT_MAX` | Max requests per IP per minute | `100` |

### PostgreSQL (Docker Compose)

| Variable | Description | Default |
|---|---|---|
| `POSTGRES_VERSION` | PostgreSQL image version | `16` |
| `POSTGRES_DB` | Database name | `db` |
| `POSTGRES_USER` | Database user | `random` |
| `POSTGRES_PASSWORD` | Database password | `random` |
| `POSTGRES_NAME` | Docker service hostname | `database` |

### ngrok (optional, local dev)

| Variable | Description |
|---|---|
| `NGROK_AUTHTOKEN` | ngrok authentication token |

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/) and Docker Compose
- [Make](https://www.gnu.org/software/make/)
- A [Clerk](https://clerk.com/) account

### Quick start

```bash
# 1. Clone and configure
cp .env.dist .env
# Edit .env with your Clerk keys and database credentials

# 2. Start the development stack
make dev

# 3. Run database migrations
make migrate

# 4. API is available at
# http://localhost:4000
```

The development stack includes:

| Service | Container | Port (host) | Description |
|---|---|---|---|
| API | `api` | `4000` | Go API with Air hot reload |
| Database | `datastore` | `9543` | PostgreSQL 16 |
| ngrok | `ngrok` | `4040` | Public tunnel for webhooks |

## Makefile Commands

| Command | Description |
|---|---|
| `make dev` | Start all services in background |
| `make migrate` | Build CLI and run database migrations |
| `make lint` | Run golangci-lint with auto-fix inside the API container |

## Docker

### Development

The development image uses [Air](https://github.com/air-verse/air) for hot reload. Source code is mounted as a volume, so changes are picked up automatically.

```bash
make dev                  # Start
docker-compose logs -f    # View logs
docker-compose down       # Stop
docker-compose restart    # Restart
```

### Production build

The Dockerfile supports a multi-stage production build:

```bash
docker build --target production -t go-api:prod .
docker run -p 3000:3000 --env-file .env go-api:prod
```

Production image builds three binaries (`api`). Only `api` is used as the default entrypoint.

## Database

### Schema

The `users` table is managed via GORM auto-migration:

| Column | Type | Notes |
|---|---|---|
| `id` | UUID | Primary key, auto-generated |
| `clerk_id` | string | Unique, Clerk user ID |
| `first_name` | string | Nullable |
| `last_name` | string | Nullable |
| `banned` | boolean | Default `false`, indexed |
| `created_at` | timestamp | Auto-set on create |
| `updated_at` | timestamp | Auto-set on update |

### Migrations

Migrations are handled by the CLI:

```bash
make migrate
```

This runs `AutoMigrate` inside a transaction for all registered entities. Add new entities to `cmd/cli/command/migrate.go` as the project grows.

### Connection pool

The database connection pool is configured in `infrastructure/config/database.go`:

- Max open connections: 25
- Max idle connections: 5
- Connection max lifetime: 5 minutes
- Idle connection max lifetime: 10 minutes

## Pagination

The `infrastructure/paginate` package provides reusable pagination for list endpoints:

```go
type PaginateQuery struct {
    Page    int    `query:"page"`
    Limit   int    `query:"limit"`
    SortBy  string `query:"sortBy"`
    OrderBy string `query:"orderBy"`
    Search  string `query:"search"`
}
```

Use `Normalize()` to apply defaults (page 1, limit 20, max 100) and `repository/gorm/helper.go` `Paginate()` for GORM queries.

## Middleware

The API applies the following middleware stack (in order):

1. **Helmet** — security headers
2. **CORS** — cross-origin configuration
3. **Rate limiter** — per-IP request throttling
4. **Logger** — request logging
5. **Server-Timing** — response duration header

## Extending the Project

### Adding a new entity

1. Create the entity in `domain/entity/`.
2. Define the repository interface in `domain/repository/`.
3. Implement it in `repository/gorm/`.
4. Create use cases in `usecase/`.
5. Add handlers in `handler/` and presenters in `presenter/`.
6. Wire dependencies in `cmd/api/wire/wire.go`.
7. Register routes in `cmd/api/routes.go`.
8. Add the entity to `cmd/cli/command/migrate.go`.

### Adding a new API route

1. Create the handler method.
2. Register it in `cmd/api/routes.go` under the appropriate group.
3. Use `AuthenticateMiddleware.Protected()` for routes that require auth.

## Troubleshooting

| Problem | Solution |
|---|---|
| `package go-api/infrastructure/paginate is not in std` | Ensure `infrastructure/paginate/paginate.go` exists |
| `required environment variable X is not set` | Add the missing variable to `.env` (see [Environment Variables](#environment-variables)) |
| `401 Invalid token` | Verify `CLERK_FRONTEND_API` matches your Clerk app's Frontend API URL |
| `401 invalid signature` on webhook | Verify `CLERK_WEBHOOK_SECRET` matches the Clerk endpoint signing secret |
| Database connection refused | Ensure the `database` container is healthy: `docker-compose ps` |
| JWKS fetch failure on startup | Check network access to `CLERK_FRONTEND_API` from the container |
| Port 4000 already in use | Change the host port mapping in `compose.yaml` |

## License

This project is a template — use it freely for your own applications.
