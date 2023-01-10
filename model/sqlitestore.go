package model

import (
	"database/sql"
	"embed"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/go-co-op/gocron"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	db        *sql.DB
	scheduler *gocron.Scheduler
	logger    *log.Logger
}

//go:embed sqlite-migrations/*.sql
var migrationsFs embed.FS

func (c *SQLiteStore) runMigrations() error {
	goose.SetBaseFS(migrationsFs)
	goose.SetLogger(c.logger)
	err := goose.SetDialect("sqlite3")
	if err != nil {
		return fmt.Errorf("error setting goose dialect for migration: %w", err)
	}
	err = goose.Up(c.db, "sqlite-migrations")
	if err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}
	return nil
}

func (c *SQLiteStore) CleanUp() error {
	log.Println("Vacuuming database")
	stmt := "VACUUM;"
	_, err := c.db.Exec(stmt)
	if err != nil {
		return fmt.Errorf("error performing db vaccum (sqlite clean up): %w", err)
	}
	return nil
}

func (c *SQLiteStore) Backup() (io.Reader, error) {
	stmt := "VACUUM INTO 'backup.sqlite'"
	_, err := c.db.Exec(stmt)
	if err != nil {
		return nil, fmt.Errorf("error creating backup: %w", err)
	}
	f, err := os.Open("backup.sqlite")
	if err != nil {
		return nil, fmt.Errorf("error opening backup: %w", err)
	}

	backupReader := BackupReader{*f}
	return backupReader, nil
}

type BackupReader struct {
	f os.File
}

func (b BackupReader) Read(p []byte) (int, error) {
	n, err := b.f.Read(p)
	if err == io.EOF {
		err := b.f.Close()
		if err != nil {
			return n, fmt.Errorf("error closing backup after EOF when reading: %w", err)
		}
		err = os.Remove("backup.sqlite")
		if err != nil {
			return n, fmt.Errorf("error deleting backup file after EOF when reading: %w", err)
		}
	}
	return n, err
}

func (ss *SQLiteStore) GetDBConnection() *sql.DB {
	return ss.db
}

func NewSQLiteStore(scheduler *gocron.Scheduler, logger *log.Logger) (*SQLiteStore, error) {
	path := "./db/comments.sqlite"
	pragma := "_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)&_pragma=synchronous(NORMAL)&_pragma=journal_size_limit(100000000)"
	db, err := sql.Open("sqlite", fmt.Sprintf("%s?%s", path, pragma))
	if err != nil {
		return nil, fmt.Errorf("error opening comments database: %w", err)
	}

	store := &SQLiteStore{db, scheduler, logger}

	return store, nil
}
