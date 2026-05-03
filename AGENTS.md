# AGENTS.md

## Purpose

This file documents how coding agents should work in the **safa-went** project.

> For a quick introduction, see [CLAUDE.md](./CLAUDE.md).
├── main.go                        # Entry point, wires router and starts server
├── go.mod
├── .env.example                   # Environment variable reference
├── .gitignore
├── database/
│   ├── migrations/                # SQL migration files (*.up.sql / *.down.sql)
│   ├── models/                    # GORM model structs
│   └── seeders/                   # Database seeders
├── http/
│   ├── controllers/               # Request handlers (one file per resource)
│   ├── middlewares/               # auth.go, cors.go
│   ├── requests/                  # Request validation structs
│   └── resources/                 # Response transformation structs
├── internal/
│   ├── config/config.go           # Env-based config loader
│   ├── responses/error_response.go # Shared API error response helpers
│   └── providers/database_provider.go  # GORM DB initializer
└── routes/
    ├── main_router.go             # Gin engine setup + middleware mounting
    └── *_router.go                # Per-resource route groups
```

## Conventions

- **Models** live in `database/models/`. Embed `gorm.Model` for standard fields.
- **Controllers** live in `http/controllers/`. One file per resource, named `<resource>_controller.go`.
- Shared controller helpers live in `http/controllers/helpers.go` (e.g. `ParseID`).
- **Requests** live in `http/requests/`. One file per resource, named `<resource>_request.go`.
- **Resources** live in `http/resources/`. One file per resource, named `<resource>_resource.go`.
- **Routers** live in `routes/`. One file per resource, named `<resource>_router.go`.
- **Migrations** live in `database/migrations/` as numbered pairs: `NNNNNN_<desc>.up.sql` / `NNNNNN_<desc>.down.sql`.
- Keep file names lowercase.
- Commands other than `create` must be run in a directory that contains `wentconfig.json`.
- Do not declare payload/response DTOs in controllers; use `http/requests` and `http/resources` types.
- In default generation mode, fill skeleton implementations manually.

## Generation Modes

- Default `gen:* <Name>` commands generate related files in skeleton form.
- `gen:* <Name> --all` generates related files in full form.
- `gen:model -m` is preserved for compatibility; skeleton related generation already includes migration files.

## Environment Variables

See `.env.example` for all required variables. Copy it to `.env` before running.
If no env file is found, configuration defaults still fall back to SQLite.

## Key Commands

| Command | Description |
|---------|-------------|
| `go run main.go` | Start the development server |
| `swag init` | Regenerate Swagger docs |
| `went2 gen:model <Name>` | Scaffold a new model |
| `went2 gen:resource <Name>` | Scaffold a new resource |
| `went2 gen:controller <Name>` | Scaffold a new controller |
| `went2 gen:router <Name>` | Scaffold a new router |
| `went2 gen:migration <Name>` | Scaffold a new migration |
| `went2 migrate` | Run pending migrations |
| `went2 migrate:rollback` | Roll back the last migration |
| `went2 migrate:fresh` | Drop all tables and re-run migrations |

## Editing Guidance

- Register new routes in the corresponding `routes/*_router.go` file and mount the group in `routes/main_router.go`.
- After adding or changing model structs, create a new migration pair and run `went2 migrate`.
- Keep read/query logic in `http/resources/*_resource.go` via `XQuery` methods; keep controllers thin.
- Use `Paginate(page, perPage)` in resources for Laravel-style `{data, meta}` responses.
- Use `internal/responses` helpers in controllers for error output (`JSONError` / `JSONErrorWithDetails`).
- Keep Swagger failure schemas aligned with `responses.ErrorBody`.
- After changing any controller annotations, run `swag init` to update docs.
- Do not modify generated migration files that have already been applied.
