package webmentions

import (
	"database/sql"
	"fmt"
	"math"
	"net/url"
	"strings"
	"tiim/go-comment-api/model"
	"time"

	"github.com/google/uuid"
)

type webmentionsStore struct {
	db *sql.DB
}

type QueuedWebmention struct {
	webmention *Webmention
	tries      int
	nextTry    time.Time
}

type Webmention struct {
	Id     string
	Source string
	Target string
}

func NewStore(store *model.SQLiteStore) *webmentionsStore {
	s := &webmentionsStore{db: store.GetDBConnection()}

	return s
}

func (s *webmentionsStore) scheduleForProcessing(w *Webmention) error {
	_, err := s.db.Exec("INSERT INTO webmentions_queue (id, source, target) VALUES (?, ?, ?)", w.Id, w.Source, w.Target)
	return err
}

func newWebmention(source, target string) (*Webmention, error) {

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
		Id:     uuid.New().String(),
		Source: source,
		Target: target,
	}, nil
}

func (s *webmentionsStore) getNextWebmentionFromQueue() (*QueuedWebmention, error) {
	row := s.db.QueryRow(
		"SELECT id, source, target, tries, next_try FROM webmentions_queue WHERE next_try <= CURRENT_TIMESTAMP ORDER BY timestamp LIMIT 1",
		time.Now().Format(time.RFC3339),
	)

	w := &Webmention{}
	q := &QueuedWebmention{
		webmention: w,
	}
	var nextTry string
	err := row.Scan(&w.Id, &w.Source, &w.Target, &q.tries, &nextTry)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("could not get next webmention from queue: %w", err)
	}
	q.nextTry, err = time.Parse(time.RFC3339, nextTry)
	if err != nil {
		return nil, fmt.Errorf("could not parse next try time: %w", err)
	}
	return q, nil
}

func (s *webmentionsStore) deleteFromQueue(w *QueuedWebmention) error {
	_, err := s.db.Exec("DELETE FROM webmentions_queue WHERE id = ?", w.webmention.Id)
	return err
}

func (s *webmentionsStore) updateNextTry(w *QueuedWebmention) error {
	seconds := math.Exp2(float64(w.tries)) * 10
	w.nextTry = time.Now().Add(time.Duration(seconds) * time.Second)

	w.tries++

	if w.tries > 10 {
		return s.deleteFromQueue(w)
	}

	_, err := s.db.Exec("UPDATE webmentions_queue SET next_try = ?, tries = ? WHERE id = ?",
		w.nextTry.Format(time.RFC3339), w.tries, w.webmention.Id)
	return err
}

func (s *webmentionsStore) moveWebmentionFromQueueToProcessed(w *QueuedWebmention) error {
	_, err := s.db.Exec("INSERT INTO webmentions_processed (id, source, target) VALUES (?, ?, ?)",
		w.webmention.Id, w.webmention.Source, w.webmention.Target)
	if err != nil {
		return err
	}
	return s.deleteFromQueue(w)
}

func (wm *Webmention) String() string {
	return fmt.Sprintf("Webmention: %s -> %s", wm.Source, wm.Target)
}

func (qwm *QueuedWebmention) String() string {
	return fmt.Sprintf("Queued%s (tries: %d, next try: %s)", qwm.webmention, qwm.tries, qwm.nextTry)
}
