package postgres

import (
	"context"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func (r *PostgresRepository) Bootstrap(ctx context.Context) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(r.db.DB, "migrations"); err != nil {
		return err
	}
	return nil
}
