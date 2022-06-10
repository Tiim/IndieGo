package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type Comment struct {
	Id        string `json:"id"`
	ReplyTo   string `json:"reply_to"`
	Timestamp string `json:"timestamp"`
	Page      string `json:"page"`
	Content   string `json:"content"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Notify    bool   `json:"notify"`
}

type CommentStore struct {
	db *sql.DB
}

func (cs *CommentStore) NewComment(c *Comment) error {
	c.Id = uuid.New().String()
	c.Timestamp = time.Now().Format(time.RFC3339)
	stmt := "INSERT INTO comments (id, reply_to, timestamp, page, content, name, email) VALUES (?, ?, ?, ?, ?, ?, ?);"
	res, err := cs.db.Exec(stmt, c.Id, c.ReplyTo, c.Timestamp, c.Page, c.Content, c.Name, c.Email)
	if err != nil {
		return fmt.Errorf("error inserting comment: %w", err)
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error inserting comment: %w", err)
	}

	fmt.Println("Inserted", ra, "rows")
	return nil
}

func (cs *CommentStore) GetAllComments() ([]Comment, error) {
	stmt := "SELECT id, reply_to, timestamp, page, content, name FROM comments;"
	rows, err := cs.db.Query(stmt)
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

func (cs *CommentStore) GetCommentsForPost(page string) ([]Comment, error) {
	stmt := "SELECT id, reply_to, timestamp, page, content, name FROM comments WHERE page = ?;"
	rows, err := cs.db.Query(stmt, page)
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

func (cs *CommentStore) DeleteComment(id string) error {
	stmt := "DELETE FROM comments WHERE id = ?;"
	_, err := cs.db.Exec(stmt, id)
	if err != nil {
		return fmt.Errorf("error deleting comment: %w", err)
	}
	return nil
}

func initTable(db *sql.DB) error {
	stmt := `CREATE TABLE IF NOT EXISTS comments (
		id TEXT not null primary key, 
		reply_to TEXT,
		timestamp TEXT not null,
		page TEXT not null,
		content TEXT not null,
		name TEXT not null,
		email TEXT,
		notify INTEGER not null default FALSE,
		FOREIGN KEY(reply_to) REFERENCES comments(id)
	);`
	_, err := db.Exec(stmt)
	if err != nil {
		return fmt.Errorf("error creating comments table: %w", err)
	}
	return nil
}

func NewCommentStore() (*CommentStore, error) {
	db, err := sql.Open("sqlite3", "./db/comments.sqlite")
	if err != nil {
		return nil, fmt.Errorf("error opening comments database: %w", err)
	}
	err = initTable(db)
	if err != nil {
		return nil, err
	}
	return &CommentStore{db}, nil
}
