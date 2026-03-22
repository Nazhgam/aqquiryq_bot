package migrations

import (
	"embed"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type migrationLog struct{}

func (m migrationLog) Fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

func (m migrationLog) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

//go:embed database/*.sql
var embedMigrations embed.FS

func Migration(pool *pgxpool.Pool) error {
	goose.SetLogger(&migrationLog{})
	goose.SetBaseFS(embedMigrations)
	log.Println("completed load sql for migration")
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to set dialect: %v", err)
	}
	log.Println("completed set dialect")
	db := stdlib.OpenDBFromPool(pool)
	log.Println("completed open database connection(std)")
	version, err := goose.EnsureDBVersion(db)
	if err != nil {
		return fmt.Errorf("migration: %w", err)
	}

	log.Printf("db version: %d", version)
	if err = goose.Up(db, "database"); err != nil {
		return fmt.Errorf("migration: %w", err)
	}

	if err = db.Close(); err != nil {
		return fmt.Errorf("migration: %w", err)
	}

	return nil
}
