package model

import (
	"io"
)

type CleanupStore interface {
	CleanUp() error
}

type BackupStore interface {
	Backup() (io.Reader, error)
}
