# DB Migrate

Supported databases:

* PostgreSQL;

## Usage

Docker:

```
docker run --rm -v/path/to/migrations:/migrations borodyadka/db-migrate up -d postgres://user:pass@host:5432/database
```

CLI:

```
migrate up -d postgres://user:pass@host:5432/database -s /path/to/migrations
```

Commands:

* `up [n]` — apply migrations, default `n` is `MaxInt32`
* `down [n]` — revert migrations, default `n` is `1`
* `create <name>` — create new migration
* `test` — tests connection

## Configuration

* `DATABSE_URL` or `--db-url`/`-d` in format `postgres://user:pass@host:5432/database?sslmode=disable`
* `HISTORY_TABLE` or `--history-table`/`-t`, default is `migrations_history`
* `LOG_LEVEL` or `--verbosity`/`-v` one of `debug`, `info`, `warning`, `error`
* `SOURCE_DIR` or `--source`/`-s` path to migrations dir, default is pwd

## TODO

* [ ] tests

## License

[MIT](LICENSE)
