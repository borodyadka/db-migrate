# DB Migrate

Supported databases:

* PostgreSQL`;

## Configuration

* `DATABSE_URL` or `--db-url`/`-d` in format `postgres://user:pass@host:5432/database?sslmode=disable`
* `HISTORY_TABLE` or `--history-table`/`-t`, default is `migrations_history`
* `LOG_LEVEL` or `--verbosity`/`-v` one of `debug`, `info`, `warning`, `error`
* `SOURCE_DIR` or `--source`/`-s` path to migrations dir
