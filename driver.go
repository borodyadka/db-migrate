package migrate

import (
	"context"
)

type Driver interface {
	EnsureMigrationsTable(ctx context.Context, table *Table) error
	ApplyMigration(ctx context.Context, history *Table, migration *Entry) error
	RevertMigration(ctx context.Context, history *Table, migration *Entry) error
	MigrationsHistory(ctx context.Context, history *Table) ([]Entry, error)
}
