package webmentions

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"tiim/go-comment-api/model"
	"time"

	"github.com/google/uuid"
)

type webmentionsStore struct {
	db    *sql.DB
	queue chan *QueuedWebmention
}

type QueuedWebmention struct {
	webmention *Webmention
}

type Webmention struct {
	Id        string
	Source    string
	Target    string
	TsCreated time.Time
	TsUpdated time.Time
}

func NewWebmention(source, target string) (*Webmention, error) {

	sourceUrl, err := url.ParseRequestURI(source)

	if err != nil || !strings.HasPrefix(sourceUrl.Scheme, "http") {
		return nil, fmt.Errorf("invalid source url: %w", err)
	}

	targetUrl, err := url.ParseRequestURI(target)

	if err != nil || !strings.HasPrefix(targetUrl.Scheme, "http") {
		return nil, fmt.Errorf("invalid target url: %w", err)
	}

	if *sourceUrl == *targetUrl {
		return nil, fmt.Errorf("source and target are the same")
	}

	return &Webmention{
		Id:        uuid.New().String(),
		Source:    source,
		Target:    target,
		TsCreated: time.Now(),
		TsUpdated: time.Now(),
	}, nil
}

func (w *Webmention) SourceUrl() *url.URL {
	u, _ := url.Parse(w.Source)
	return u
}

func NewStore(store *model.SQLiteStore) *webmentionsStore {
	s := &webmentionsStore{db: store.GetDBConnection(), queue: make(chan *QueuedWebmention, 20)}
	s.populateQueue()
	return s
}

func (s *webmentionsStore) ScheduleForProcessing(w *Webmention) error {
	_, err := s.db.Exec("INSERT INTO webmentions_queue (id, source, target, timestamp) VALUES (?, ?, ?, ?)", w.Id, w.Source, w.Target, w.TsCreated)

	if err != nil {
		return fmt.Errorf("could not insert webmention into queue: %w", err)
	}

	s.queue <- &QueuedWebmention{webmention: w}
	return nil
}

func (s *webmentionsStore) NextWebmentionFromQueue() (*QueuedWebmention, error) {
	mention := <-s.queue
	return mention, nil
}

func (s *webmentionsStore) MarkInvalid(w *QueuedWebmention, reason string) error {
	tx, err := s.db.Begin()
	defer tx.Rollback()

	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}

	_, err = s.db.Exec("DELETE FROM webmentions_queue WHERE id = ?", w.webmention.Id)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO webmentions_rejected (id, source, target, timestamp, reason) VALUES (?, ?, ?, ?, ?)",
		w.webmention.Id, w.webmention.Source, w.webmention.Target, w.webmention.TsCreated, reason)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *webmentionsStore) MarkSuccess(w *QueuedWebmention) error {
	tx, err := s.db.Begin()
	defer tx.Rollback()

	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	_, err = tx.Exec("INSERT INTO webmentions (id, source, target, ts_created) VALUES (?, ?, ?, ?)",
		w.webmention.Id, w.webmention.Source, w.webmention.Target, w.webmention.TsCreated)
	if err != nil {
		return fmt.Errorf("could not insert queued webmention to webmention list: %w", err)
	}
	_, err = tx.Exec("DELETE FROM webmentions_queue WHERE id = ?", w.webmention.Id)
	if err != nil {
		return fmt.Errorf("could not delete webmention from queue: %w", err)
	}

	return tx.Commit()
}

func (s *webmentionsStore) populateQueue() error {
	rows, err := s.db.Query("SELECT id, source, target, timestamp FROM webmentions_queue ORDER BY TIMESTAMP ASC")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id, source, target, timestamp string
		err := rows.Scan(&id, &source, &target, &timestamp)
		if err != nil {
			return err
		}

		ts, err := time.Parse(time.RFC3339, timestamp)
		if err != nil {
			return err
		}

		s.queue <- &QueuedWebmention{
			webmention: &Webmention{
				Id:        id,
				Source:    source,
				Target:    target,
				TsCreated: ts,
				TsUpdated: time.Now(),
			},
		}
	}
	return nil
}

func (wm *Webmention) String() string {
	return fmt.Sprintf("Webmention: %s -> %s", wm.Source, wm.Target)
}

func (qwm *QueuedWebmention) String() string {
	return fmt.Sprintf("Queued%s", qwm.webmention)
}
