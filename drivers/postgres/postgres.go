package postgres

import (
	"context"
	"fmt"
	"github.com/borodyadka/db-migrate"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
	"net/url"
	"time"
)

const Dialect migrate.Dialect = "postgres"

type Driver struct {
	conn    *pgx.Conn
	datadir string
}

func (d *Driver) EnsureSchema() error {
	return nil
}

func (d *Driver) EnsureMigrationsTable(ctx context.Context, table *migrate.Table) error {
	log.Debug("ensure schema")
	_, err := d.conn.Exec(ctx, fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s`, table.String()))
	if err != nil {
		return err
	}

	log.Debug("ensure migrations table")
	_, err = d.conn.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			"ts" BIGINT NOT NULL,
			"name" TEXT NOT NULL,
			UNIQUE ("ts", "name")
		)
	`, table.String()))
	return err
}

func (d *Driver) ApplyMigration(ctx context.Context, history *migrate.Table, migration *migrate.Entry) error {
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	log.Debug("update history table")
	if _, err := tx.Exec(
		ctx,
		fmt.Sprintf(`INSERT INTO %s ("ts", "name") VALUES ($1, $2)`, history.String()),
		migration.Timestamp, migration.Name,
	); err != nil {
		log.Error(err)
		return err
	}

	log.Debug("read migration")
	m, err := migration.ReadUp(d.datadir)
	if err != nil {
		return err
	}
	log.Debug("apply migration")
	if _, err := d.conn.Exec(ctx, m); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (d *Driver) RevertMigration(ctx context.Context, history *migrate.Table, migration *migrate.Entry) error {
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	log.Debug("update history table")
	if _, err := tx.Exec(
		ctx,
		fmt.Sprintf(`DELETE FROM %s WHERE ts = $1 AND name = $2`, history.String()),
		migration.Timestamp, migration.Name,
	); err != nil {
		log.Error(err)
		return err
	}

	log.Debug("read migration")
	m, err := migration.ReadDown(d.datadir)
	if err != nil {
		return err
	}
	log.Debug("revert migration")
	if _, err := d.conn.Exec(ctx, m); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (d *Driver) MigrationsHistory(ctx context.Context, history *migrate.Table) ([]migrate.Entry, error) {
	rows, err := d.conn.Query(
		ctx,
		fmt.Sprintf(`SELECT "ts", "name" FROM %s ORDER BY ts ASC, name ASC`, history.String()),
	)
	if err != nil {
		return nil, err
	}
	var result []migrate.Entry
	for rows.Next() {
		e := migrate.Entry{}
		if err := rows.Scan(&e.Timestamp, &e.Name); err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, nil
}

func (d *Driver) Ping(ctx context.Context) (bool, error) {
	err := d.conn.Ping(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

func NewDriver(ctx context.Context, datadir string, url *url.URL) (*Driver, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, url.String())
	if err != nil {
		return nil, err
	}
	return &Driver{
		conn:    conn,
		datadir: datadir,
	}, nil
}
