package wmsend

import (
	"database/sql"
	"fmt"
	"time"
)

type WmSendStore interface {
	IsItemUpdated(f FeedItem) (bool, error)
	GetUrlsForFeedItem(f FeedItem) ([]string, error)
	SetUrlsForFeedItem(f FeedItem, urls []string) error
}

type wmSendSqliteStore struct {
	db *sql.DB
}

func newWmSendStore(db *sql.DB) *wmSendSqliteStore {
	return &wmSendSqliteStore{db: db}
}

func (s *wmSendSqliteStore) IsItemUpdated(f FeedItem) (bool, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return false, fmt.Errorf("unable to start transaction: %w", err)
	}
	defer tx.Rollback()

	row := tx.QueryRow("SELECT updated FROM wm_send_feed_item WHERE id = ?", f.uid)

	var updated string
	err = row.Scan(&updated)
	if err != nil {
		if err != sql.ErrNoRows {
			return true, fmt.Errorf("unable to scan updated: %w", err)
		}
	}

	_, err = tx.Exec("INSERT INTO wm_send_feed_item (id, updated) VALUES (?, ?) ON CONFLICT (id) DO UPDATE SET updated = excluded.updated", f.uid, f.updated.Format(time.RFC3339))

	if err != nil {
		return true, fmt.Errorf("unable to insert updated: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return true, fmt.Errorf("unable to commit transaction: %w", err)
	}

	if updated == "" {
		return true, nil
	}
	updatedTime, err := time.Parse(time.RFC3339, updated)

	if err != nil {
		return false, fmt.Errorf("unable to parse updated timestamp: %w", err)
	}

	return updatedTime.Before(*f.updated), nil

}

func (s *wmSendSqliteStore) GetUrlsForFeedItem(f FeedItem) ([]string, error) {
	rows, err := s.db.Query("SELECT url FROM wm_send_feed_item_url WHERE feed_item_id = ?", f.uid)
	if err != nil {
		return nil, fmt.Errorf("unable to query urls: %w", err)
	}

	var urls []string

	for rows.Next() {
		var url string
		err = rows.Scan(&url)
		if err != nil {
			return nil, fmt.Errorf("unable to scan url: %w", err)
		}
		urls = append(urls, url)
	}

	return urls, nil
}

func (s *wmSendSqliteStore) SetUrlsForFeedItem(f FeedItem, urls []string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("unable to start transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM wm_send_feed_item_url WHERE feed_item_id = ?", f.uid)
	if err != nil {
		return fmt.Errorf("unable to delete urls: %w", err)
	}

	for _, url := range urls {
		_, err = tx.Exec("INSERT INTO wm_send_feed_item_url (feed_item_id, url) VALUES (?, ?)", f.uid, url)
		if err != nil {
			return fmt.Errorf("unable to insert url: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("unable to commit transaction: %w", err)
	}

	return nil
}
