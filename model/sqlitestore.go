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
	stmt := "INSERT INTO comments (id, reply_to, timestamp, page, content, name, email) VALUES (?, ?, ?, ?, ?, ?, ?);"
	res, err := cs.db.Exec(stmt, c.Id, c.ReplyTo, c.Timestamp, c.Page, c.Content, c.Name, c.Email)
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
	stmt := "SELECT id, reply_to, timestamp, page, content, name FROM comments WHERE timestamp > ? ORDER BY timestamp DESC;"
	rows, err := cs.db.Query(stmt, since.UTC().Format(time.RFC3339))
	if err != nil {
		return nil, fmt.Errorf("error querying comments: %w", err)
	}
	defer rows.Close()

	comments := make([]Comment, 0)
	for rows.Next() {

		var id string
		var reply_to string
		var timestamp string
		var page string
		var content string
		var name string

		err := rows.Scan(&id, &reply_to, &timestamp, &page, &content, &name)
		if err != nil {
			return nil, fmt.Errorf("error listing all comments: %w", err)
		}
		comment := Comment{Id: id, ReplyTo: reply_to, Timestamp: timestamp, Page: page, Content: content, Name: name, Email: ""}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error listing all comments: %w", err)
	}

	return comments, nil
}

func (cs *SQLiteStore) GetCommentsForPost(page string, since time.Time) ([]Comment, error) {
	stmt := "SELECT id, reply_to, timestamp, page, content, name FROM comments WHERE page = ? AND timestamp > ? ORDER BY timestamp DESC;"
	rows, err := cs.db.Query(stmt, page, since.UTC().Format(time.RFC3339))
	if err != nil {
		return nil, fmt.Errorf("error querying comments for page %s: %w", page, err)
	}
	defer rows.Close()

	comments := make([]Comment, 0)
	for rows.Next() {
		var id string
		var reply_to string
		var timestamp string
		var page string
		var content string
		var name string
		err := rows.Scan(&id, &reply_to, &timestamp, &page, &content, &name)
		if err != nil {
			return nil, err
		}
		comment := Comment{
			Id:        id,
			ReplyTo:   reply_to,
			Timestamp: timestamp,
			Page:      page,
			Content:   content,
			Name:      name,
			Email:     "",
		}
		comments = append(comments, comment)
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
	stmt := "SELECT id, reply_to, timestamp, page, content, name FROM comments WHERE id = ?;"
	rows, err := cs.db.Query(stmt, id)
	if err != nil {
		return Comment{}, fmt.Errorf("error querying comment with id %s: %w", id, err)
	}
	defer rows.Close()

	var comment Comment

	for rows.Next() {
		var id string
		var reply_to string
		var timestamp string
		var page string
		var content string
		var name string
		err := rows.Scan(&id, &reply_to, &timestamp, &page, &content, &name)
		if err != nil {
			return Comment{}, err
		}
		comment = Comment{
			Id:        id,
			ReplyTo:   reply_to,
			Timestamp: timestamp,
			Page:      page,
			Content:   content,
			Name:      name,
			Email:     "",
		}
	}

	if err := rows.Err(); err != nil {
		return Comment{}, fmt.Errorf("error querying comment with id %s: %w", id, err)
	}
	return comment, nil
}

func (c *SQLiteStore) Unsubscribe(secret string) (Comment, error) {
	stmt := "UPDATE comments SET notify = FALSE WHERE unsubscribe_secret = ?;"
	_, err := c.db.Exec(stmt, secret)
	if err != nil {
		return Comment{}, fmt.Errorf("error unsubscribing: %w", err)
	}

	stmt = "SELECT id, reply_to, timestamp, page, content, name, email FROM comments WHERE unsubscribe_secret = ?;"
	rows, err := c.db.Query(stmt, secret)
	if err != nil {
		return Comment{}, fmt.Errorf("error querying unsubscribed comment: %w", err)
	}
	defer rows.Close()
	var comment Comment
	for rows.Next() {
		var id string
		var reply_to string
		var timestamp string
		var page string
		var content string
		var name string
		var email string
		err := rows.Scan(&id, &reply_to, &timestamp, &page, &content, &name, &email)
		if err != nil {
			return Comment{}, err
		}
		comment = Comment{
			Id:        id,
			ReplyTo:   reply_to,
			Timestamp: timestamp,
			Page:      page,
			Content:   content,
			Name:      name,
			Email:     email,
		}
	}
	if err := rows.Err(); err != nil {
		return Comment{}, fmt.Errorf("error querying unsubscribed comment: %w", err)
	}
	return comment, nil
}

func (c *SQLiteStore) UnsubscribeAll(email string) ([]Comment, error) {
	stmt := "UPDATE comments SET notify = FALSE WHERE email = ?;"
	_, err := c.db.Exec(stmt, email)
	if err != nil {
		return nil, fmt.Errorf("error unsubscribing: %w", err)
	}

	stmt = "SELECT id, reply_to, timestamp, page, content, name, email FROM comments WHERE email = ?;"
	rows, err := c.db.Query(stmt, email)
	if err != nil {
		return nil, fmt.Errorf("error querying unsubscribed comments: %w", err)
	}
	defer rows.Close()
	comments := make([]Comment, 0)
	for rows.Next() {
		var id string
		var reply_to string
		var timestamp string
		var page string
		var content string
		var name string
		var email string
		err := rows.Scan(&id, &reply_to, &timestamp, &page, &content, &name, &email)
		if err != nil {
			return nil, err
		}
		comment := Comment{
			Id:        id,
			ReplyTo:   reply_to,
			Timestamp: timestamp,
			Page:      page,
			Content:   content,
			Name:      name,
			Email:     email,
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error querying unsubscribed comments: %w", err)
	}
	return comments, nil
}

//go:embed sqlite-migrations/*.sql
var migrationsFs embed.FS

func runMigrations(db *sql.DB) error {
	goose.SetBaseFS(migrationsFs)
	err := goose.SetDialect("sqlite3")
	if err != nil {
		return fmt.Errorf("error setting goose dialect for migration: %w", err)
	}
	err = goose.Up(db, "sqlite-migrations")
	if err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}
	return nil
}

func NewSQLiteStore() (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", "./db/comments.sqlite")
	if err != nil {
		return nil, fmt.Errorf("error opening comments database: %w", err)
	}
	err = runMigrations(db)
	if err != nil {
		return nil, err
	}
	return &SQLiteStore{db}, nil
}
