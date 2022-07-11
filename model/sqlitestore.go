package model

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"

	"github.com/pressly/goose/v3"
)

type SQLiteStore struct {
	db *sql.DB
}

func (cs *SQLiteStore) NewComment(c *Comment) error {
	c.Id = uuid.New().String()
	c.Timestamp = time.Now().UTC().Format(time.RFC3339)
	stmt := "INSERT INTO comments (id, reply_to, timestamp, page, content, name, email, notify) VALUES (?, ?, ?, ?, ?, ?, ?, ?);"

	replyTo := &c.ReplyTo
	if *replyTo == "" {
		replyTo = nil
	}

	res, err := cs.db.Exec(stmt, c.Id, replyTo, c.Timestamp, c.Page, c.Content, c.Name, c.Email, c.Notify)
	if err != nil {
		return fmt.Errorf("error inserting comment: %w", err)
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error inserting comment: %w", err)
	}

	log.Println("Inserted", ra, "comments")
	return nil
}

func (cs *SQLiteStore) GetAllComments(since time.Time) ([]Comment, error) {
	stmt := "SELECT * FROM comments WHERE timestamp > ? ORDER BY timestamp DESC;"
	rows, err := cs.db.Query(stmt, since.UTC().Format(time.RFC3339))
	if err != nil {
		return nil, fmt.Errorf("error querying comments: %w", err)
	}
	defer rows.Close()

	comments := make([]Comment, 0)
	for rows.Next() {
		comment, err := readRow(rows)
		if err != nil {
			log.Printf("error reading all comments: %v", err)
			return nil, fmt.Errorf("error reading all comments: %w", err)
		}
		comments = append(comments, *comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error listing all comments: %w", err)
	}

	return comments, nil
}

func (cs *SQLiteStore) GetCommentsForPost(page string, since time.Time) ([]Comment, error) {
	stmt := "SELECT * FROM comments WHERE page = ? AND timestamp > ? ORDER BY timestamp DESC;"
	rows, err := cs.db.Query(stmt, page, since.UTC().Format(time.RFC3339))
	if err != nil {
		return nil, fmt.Errorf("error querying comments for page %s: %w", page, err)
	}
	defer rows.Close()

	comments := make([]Comment, 0)
	for rows.Next() {
		comment, err := readRow(rows)
		if err != nil {
			log.Printf("error reading comments: %v", err)
			return nil, fmt.Errorf("error reading comments for page %s: %w", page, err)
		}
		comments = append(comments, *comment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error querying comments for page %s: %w", page, err)
	}
	return comments, nil
}

func (cs *SQLiteStore) DeleteComment(id string) error {
	stmt := "DELETE FROM comments WHERE id = ?;"
	_, err := cs.db.Exec(stmt, id)
	if err != nil {
		return fmt.Errorf("error deleting comment: %w", err)
	}
	return nil
}

func (cs *SQLiteStore) GetComment(id string) (Comment, error) {
	stmt := "SELECT * FROM comments WHERE id = ?;"
	rows, err := cs.db.Query(stmt, id)
	if err != nil {
		return Comment{}, fmt.Errorf("error querying comment with id %s: %w", id, err)
	}
	defer rows.Close()

	var comment *Comment

	for rows.Next() {
		comment, err = readRow(rows)
		if err != nil {
			log.Printf("error reading comment: %v", err)
			return Comment{}, fmt.Errorf("error reading comment: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return Comment{}, fmt.Errorf("error querying comment with id %s: %w", id, err)
	}
	return *comment, nil
}

func (c *SQLiteStore) Unsubscribe(secret string) (Comment, error) {
	stmt := "UPDATE comments SET notify = FALSE WHERE unsubscribe_secret = ?;"
	_, err := c.db.Exec(stmt, secret)
	if err != nil {
		return Comment{}, fmt.Errorf("error unsubscribing: %w", err)
	}

	stmt = "SELECT * FROM comments WHERE unsubscribe_secret = ?;"
	rows, err := c.db.Query(stmt, secret)
	if err != nil {
		return Comment{}, fmt.Errorf("error querying unsubscribed comment: %w", err)
	}
	defer rows.Close()
	var comment *Comment
	for rows.Next() {
		comment, err = readRow(rows)
		if err != nil {
			log.Printf("error reading row while unsubscribing: %v", err)
			return Comment{}, fmt.Errorf("error reading row while unsubscribing: %w", err)
		}
	}
	if err := rows.Err(); err != nil {
		return Comment{}, fmt.Errorf("error querying unsubscribed comment: %w", err)
	}
	return *comment, nil
}

func (c *SQLiteStore) UnsubscribeAll(email string) ([]Comment, error) {
	stmt := "UPDATE comments SET notify = FALSE WHERE email = ?;"
	_, err := c.db.Exec(stmt, email)
	if err != nil {
		return nil, fmt.Errorf("error unsubscribing: %w", err)
	}

	stmt = "SELECT * FROM comments WHERE email = ?;"
	rows, err := c.db.Query(stmt, email)
	if err != nil {
		return nil, fmt.Errorf("error querying unsubscribed comments: %w", err)
	}
	defer rows.Close()
	comments := make([]Comment, 0)
	for rows.Next() {
		comment, err := readRow(rows)
		if err != nil {
			log.Printf("error reading row while unsubscribing: %s", err)
			return nil, fmt.Errorf("error reading row while unsubscribing: %w", err)
		} else {
			comments = append(comments, *comment)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error querying unsubscribed comments: %w", err)
	}
	return comments, nil
}

//go:embed sqlite-migrations/*.sql
var migrationsFs embed.FS

func (c *SQLiteStore) RunMigrations() error {
	goose.SetBaseFS(migrationsFs)
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

func NewSQLiteStore() (*SQLiteStore, error) {
	path := "./db/comments.sqlite"
	pragma := "_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=busy_timeout(8000)&_pragma=journal_size_limit(100000000)"
	db, err := sql.Open("sqlite", fmt.Sprintf("%s?%s", path, pragma))
	if err != nil {
		return nil, fmt.Errorf("error opening comments database: %w", err)
	}
	return &SQLiteStore{db}, nil
}

func readRow(rows *sql.Rows) (*Comment, error) {
	c := Comment{}
	var replyTo sql.NullString
	var notify bool
	var email string
	err := rows.Scan(
		&c.Id,
		&replyTo,
		&c.Timestamp,
		&c.Page,
		&c.Content,
		&c.Name,
		&email,
		&notify,
		&c.UnsubscribeSecret,
	)
	if err != nil {
		return nil, err
	}

	c.Email = Email(email)
	c.Notify = Notify(notify)
	if replyTo.Valid {
		c.ReplyTo = replyTo.String
	}

	return &c, nil
}
