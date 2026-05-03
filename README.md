# safa-went

A REST API project scaffolded with [went](https://github.com/ekinatam/went).

> **For AI Agents:** See [CLAUDE.md](./CLAUDE.md) for agent guidelines, which points to [AGENTS.md](./AGENTS.md) for detailed conventions.

- [went2](https://github.com/ekinatam/went) CLI (optional, for code generation)

## Getting Started

```bash
# 1. Copy environment file and fill in your values
cp .env.example .env

# 2. Install dependencies
go mod tidy

# 3. Run migrations
went2 migrate

# 4. Start the server
go run main.go
```

The API will be available at `http://localhost:8080`.
Swagger UI is served at `http://localhost:8080/swagger/index.html`.

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_PORT` | `8080` | HTTP server port |
| `DB_CONNECTION` | `sqlite` | Database driver (`sqlite` / `mysql` / `postgres`) |
| `DB_HOST` | `127.0.0.1` | Database host |
| `DB_PORT` | `3306` | Database port |
| `DB_NAME` | `database` | Database name |
| `DB_USER` | — | Database user |
| `DB_PASSWORD` | — | Database password |
| `DB_STORAGE` | `./database.sqlite` | SQLite file path |
| `JWT_SECRET` | — | Secret key for JWT tokens |

If no env file is present, default DB settings fall back to SQLite.

## Code Generation

Use the `went2` CLI to scaffold new resources:

```bash
went2 gen:model <Name>              # Skeleton related set
went2 gen:model <Name> -m           # Compatibility flag (skeleton related set already includes migration)
went2 gen:model <Name> -a           # Full related set
went2 gen:controller <Name>         # Skeleton related set
went2 gen:controller <Name> -a      # Full related set
went2 gen:router <Name>             # Skeleton related set
went2 gen:router <Name> -a          # Full related set
went2 gen:resource <Name>           # Skeleton related set
went2 gen:resource <Name> -a        # Full related set
went2 gen:migration <Name>          # Skeleton related set
went2 gen:migration <Name> -a       # Full related set
```

Default mode generates skeleton files so you can fill your domain details manually.
Use `--all` to generate fully populated templates.

Commands other than `create` require `wentconfig.json` in the current directory.
If the file is missing, the CLI stops with an error.

Generated router files follow `routes/<name>_router.go` naming in lowercase.

## Migrations

```bash
went2 migrate                       # Run all pending migrations
went2 migrate:rollback              # Roll back the last migration
went2 migrate:rollback --step 3     # Roll back the last 3 migrations
went2 migrate:fresh                 # Drop everything and re-run all migrations
```

## Project Structure

```
safa-went/
├── main.go
├── go.mod
├── .env.example
├── database/
│   ├── migrations/
│   ├── models/
│   └── seeders/
├── http/
│   ├── controllers/
│   ├── middlewares/
│   ├── requests/
│   └── resources/
├── internal/
│   ├── config/
│   ├── responses/
│   └── providers/
└── routes/
```

## API Documentation

After running `swag init` and starting the server, open:

```
http://localhost:8080/swagger/index.html
```

## Resource Query Layer

Generated resources under `http/resources/` include:

- `XResource`: single-item API payload
- `XCollection`: Laravel-style paginated payload with `data` and `meta`
- `XQuery`: DB read operations (`Paginate`, `Find`) that controllers delegate to

Shared pagination metadata lives in `http/resources/pagination.go`.

## Request DTO Layer

Generated request DTOs under `http/requests/` include:

- `XPayload`: create payload
- `XUpdatePayload`: partial update payload

Controllers import these request types instead of defining payload structs inline.

## Global Error Responses

Generated helper file `internal/responses/error_response.go` provides shared API error primitives:

- `ErrorBody`: standard error payload schema
- `JSONError(...)`: common JSON error response helper
- `JSONErrorWithDetails(...)`: JSON error response with optional details

Controllers are generated to use these helpers and Swagger failure responses reference `responses.ErrorBody`.
