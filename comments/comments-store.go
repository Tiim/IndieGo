package comments

import (
	"database/sql"
	"fmt"
	"log"
	"tiim/go-comment-api/event"
	"tiim/go-comment-api/model"
	"time"

	"github.com/google/uuid"
)

type CommentPageToUrlMapper func(page string, id string) string

type commentStore struct {
	db              *sql.DB
	eventHandler    event.Handler
	pageToUrlMapper CommentPageToUrlMapper
}

func NewCommentStore(db *sql.DB, pageToUrl CommentPageToUrlMapper) *commentStore {
	return &commentStore{db: db, pageToUrlMapper: pageToUrl}
}

func (cs *commentStore) SetEventHandler(h event.Handler) {
	cs.eventHandler = h
}

func (cs *commentStore) NewComment(c *comment) error {

	c.Id = uuid.New().String()
	c.Timestamp = time.Now().UTC().Format(time.RFC3339)
	stmt := "INSERT INTO comments (id, reply_to, timestamp, page, content, name, email, notify) VALUES (?, ?, ?, ?, ?, ?, ?, ?);"

	replyTo := &c.ReplyTo
	if *replyTo == "" {
		replyTo = nil
	}

	tx, err := cs.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	_, err = tx.Exec(stmt, c.Id, replyTo, c.Timestamp, c.Page, c.Content, c.Name, c.Email, c.Notify)
	if err != nil {
		return fmt.Errorf("error inserting comment: %w", err)
	}

	genericComment := c.ToGenericComment()
	ok, err := cs.eventHandler.OnNewComment(&genericComment)
	if err != nil {
		return fmt.Errorf("error handling event: %w", err)
	} else if !ok {
		log.Println("Comment rejected by event handler")
		return nil
	}

	return tx.Commit()
}

func (cs *commentStore) GetAllComments(since time.Time) ([]comment, error) {
	stmt := "SELECT * FROM comments WHERE timestamp > ? ORDER BY timestamp DESC;"
	rows, err := cs.db.Query(stmt, since.UTC().Format(time.RFC3339))
	if err != nil {
		return nil, fmt.Errorf("error querying comments: %w", err)
	}
	defer rows.Close()

	comments := make([]comment, 0)
	for rows.Next() {
		comment, err := cs.readRow(rows)
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

func (cs *commentStore) DeleteComment(id string) error {
	stmt := "DELETE FROM comments WHERE id = ?;"
	tx, err := cs.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()
	comment, err := cs.GetComment(id, tx)

	if err != nil {
		return fmt.Errorf("error getting comment: %w", err)
	}

	_, err = tx.Exec(stmt, id)
	if err != nil {
		return fmt.Errorf("error deleting comment: %w", err)
	}

	genericComment := comment.ToGenericComment()
	ok, err := cs.eventHandler.OnDeleteComment(&genericComment)
	if err != nil {
		return fmt.Errorf("error handling event: %w", err)
	} else if !ok {
		log.Println("Comment rejected by event handler")
		return nil
	}

	return tx.Commit()
}

func (cs *commentStore) GetComment(id string, tx *sql.Tx) (*comment, error) {
	stmt := "SELECT * FROM comments WHERE id = ?;"
	var rows *sql.Rows
	var err error
	if tx != nil {
		rows, err = tx.Query(stmt, id)
	} else {
		rows, err = cs.db.Query(stmt, id)
	}
	if err != nil {
		return nil, fmt.Errorf("error querying comment with id %s: %w", id, err)
	}
	defer rows.Close()

	var comment *comment

	for rows.Next() {
		comment, err = cs.readRow(rows)
		if err != nil {
			log.Printf("error reading comment: %v", err)
			return nil, fmt.Errorf("error reading comment: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error querying comment with id %s: %w", id, err)
	}
	return comment, nil
}

func (cs *commentStore) Unsubscribe(secret string) (*comment, error) {
	stmt := "UPDATE comments SET notify = FALSE WHERE unsubscribe_secret = ?;"
	_, err := cs.db.Exec(stmt, secret)
	if err != nil {
		return nil, fmt.Errorf("error unsubscribing: %w", err)
	}

	stmt = "SELECT * FROM comments WHERE unsubscribe_secret = ?;"
	rows, err := cs.db.Query(stmt, secret)
	if err != nil {
		return nil, fmt.Errorf("error querying unsubscribed comment: %w", err)
	}
	defer rows.Close()
	var comment *comment
	for rows.Next() {
		comment, err = cs.readRow(rows)
		if err != nil {
			log.Printf("error reading row while unsubscribing: %v", err)
			return nil, fmt.Errorf("error reading row while unsubscribing: %w", err)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error querying unsubscribed comment: %w", err)
	}
	return comment, nil
}

func (cs *commentStore) UnsubscribeAll(email string) ([]comment, error) {
	stmt := "UPDATE comments SET notify = FALSE WHERE email = ?;"
	_, err := cs.db.Exec(stmt, email)
	if err != nil {
		return nil, fmt.Errorf("error unsubscribing: %w", err)
	}

	stmt = "SELECT * FROM comments WHERE email = ?;"
	rows, err := cs.db.Query(stmt, email)
	if err != nil {
		return nil, fmt.Errorf("error querying unsubscribed comments: %w", err)
	}
	defer rows.Close()
	comments := make([]comment, 0)
	for rows.Next() {
		comment, err := cs.readRow(rows)
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

func (cs *commentStore) GetGenericCommentsForPage(page string, since time.Time) ([]*model.GenericComment, error) {
	fmt.Printf("SELECT * FROM comments WHERE page = '%s' AND timestamp > '%s' ORDER BY timestamp DESC;\n", page, since.UTC().Format(time.RFC3339))
	stmt := "SELECT * FROM comments WHERE page = ? AND timestamp > ? ORDER BY timestamp DESC;"
	rows, err := cs.db.Query(stmt, page, since.UTC().Format(time.RFC3339))
	if err != nil {
		return nil, fmt.Errorf("error querying comments: %w", err)
	}
	defer rows.Close()

	comments := make([]*model.GenericComment, 0)
	for rows.Next() {
		comment, err := cs.readRow(rows)
		if err != nil {
			log.Printf("error reading all comments: %v", err)
			return nil, fmt.Errorf("error reading all comments: %w", err)
		}
		genericComment := comment.ToGenericComment()
		comments = append(comments, &genericComment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error listing all comments: %w", err)
	}

	return comments, nil
}
func (cs *commentStore) GetAllGenericComments(since time.Time) ([]*model.GenericComment, error) {
	comments, err := cs.GetAllComments(since)
	if err != nil {
		return nil, err
	}
	genericComments := make([]*model.GenericComment, len(comments))
	for i, comment := range comments {
		genericComment := comment.ToGenericComment()
		genericComments[i] = &genericComment
	}
	return genericComments, nil
}

func (cs *commentStore) readRow(rows *sql.Rows) (*comment, error) {
	c := comment{}
	var replyTo sql.NullString
	err := rows.Scan(
		&c.Id,
		&replyTo,
		&c.Timestamp,
		&c.Page,
		&c.Content,
		&c.Name,
		&c.Email,
		&c.Notify,
		&c.UnsubscribeSecret,
	)
	if err != nil {
		return nil, err
	}

	c.Url = cs.pageToUrlMapper(c.Page, c.Id)
	if replyTo.Valid {
		c.ReplyTo = replyTo.String
	}

	return &c, nil
}
