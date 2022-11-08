package webmentions

import (
	"database/sql"
	"fmt"
	"log"
	"tiim/go-comment-api/event"
	"tiim/go-comment-api/model"
	"time"
)

type webmentionsStore struct {
	db           *sql.DB
	queue        chan *QueuedWebmention
	eventHandler event.Handler
}

type QueuedWebmention struct {
	webmention *Webmention
}

func NewStore(store *model.SQLiteStore) *webmentionsStore {
	s := &webmentionsStore{
		db:    store.GetDBConnection(),
		queue: make(chan *QueuedWebmention, 20),
	}
	s.populateQueue()
	return s
}

func (s *webmentionsStore) SetEventHandler(handler event.Handler) {
	s.eventHandler = handler
}

func (s *webmentionsStore) GetWebmentions() ([]*Webmention, error) {
	rows, err := s.db.Query("SELECT id, source, target, ts_created, ts_updated FROM webmentions WHERE NOT deleted ORDER BY ts_created DESC")
	if err != nil {
		return nil, fmt.Errorf("unable to query webmentions: %w", err)
	}
	defer rows.Close()

	var webmentions []*Webmention

	for rows.Next() {
		var webmention Webmention
		err := rows.Scan(&webmention.Id, &webmention.Source, &webmention.Target, &webmention.TsCreated, &webmention.TsUpdated)
		if err != nil {
			return nil, fmt.Errorf("unable to scan webmention: %w", err)
		}
		webmentions = append(webmentions, &webmention)
	}

	return webmentions, nil
}

func (s *webmentionsStore) GetWebmention(id string, tx *sql.Tx) (*Webmention, error) {

	var webmention Webmention
	query := "SELECT id, source, target, ts_created, ts_updated FROM webmentions WHERE id = ? AND NOT deleted"
	var row *sql.Row
	if tx != nil {
		row = s.db.QueryRow(query, id)
	} else {
		row = tx.QueryRow(query, id)
	}
	err := row.Scan(&webmention.Id, &webmention.Source, &webmention.Target, &webmention.TsCreated, &webmention.TsUpdated)
	if err != nil {
		return nil, fmt.Errorf("unable to query webmention: %w", err)
	}
	return &webmention, nil
}

func (s *webmentionsStore) DeleteWebmention(id string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	mention, err := s.GetWebmention(id, tx)
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE webmentions SET deleted = true WHERE id = ?", id)

	if err != nil {
		return fmt.Errorf("could not delete webmention: %w", err)
	}

	genericComment := mention.ToGenericComment()
	ok, err := s.eventHandler.OnDeleteComment(&genericComment)

	if err != nil {
		return fmt.Errorf("could not handle delete event: %w", err)
	} else if !ok {
		log.Printf("Delete rejected by event handler")
		return nil
	}

	return tx.Commit()
}

func (s *webmentionsStore) DenyListDomain(domain string) error {
	_, err := s.db.Exec("INSERT INTO domain_deny_list (domain) VALUES (?)", domain)
	if err != nil {
		return fmt.Errorf("could not insert domain to deny list: %w", err)
	}
	wm, err := s.GetWebmentions()
	if err != nil {
		return fmt.Errorf("could not get webmentions: %w", err)
	}
	for _, w := range wm {
		if w.SourceUrl().Hostname() == domain {
			err := s.DeleteWebmention(w.Id)
			if err != nil {
				return fmt.Errorf("could not delete webmention: %w", err)
			}
		}
	}

	return nil
}

func (s *webmentionsStore) GetDomainDenyList() ([]string, error) {
	rows, err := s.db.Query("SELECT domain FROM domain_deny_list")
	if err != nil {
		return nil, fmt.Errorf("unable to query domain deny list: %w", err)
	}
	defer rows.Close()

	var domains []string

	for rows.Next() {
		var domain string
		err := rows.Scan(&domain)
		if err != nil {
			return nil, fmt.Errorf("unable to scan domain: %w", err)
		}
		domains = append(domains, domain)
	}

	return domains, nil
}

func (s *webmentionsStore) DeleteDomainFromDenyList(domain string) error {
	_, err := s.db.Exec("DELETE FROM domain_deny_list WHERE domain = ?", domain)
	return err
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

	genericComment := w.webmention.ToGenericComment()
	ok, err := s.eventHandler.OnNewComment(&genericComment)
	if err != nil {
		return fmt.Errorf("error handling event: %w", err)
	} else if !ok {
		log.Println("Webmention rejected by event handler")
		return nil
	}

	return tx.Commit()
}

func (s *webmentionsStore) GetAllGenericComments(since time.Time) ([]*model.GenericComment, error) {
	rows, err := s.db.Query("SELECT id, source, target, ts_created, ts_updated FROM webmentions WHERE deleted = false AND ts_updated > ?", since)
	if err != nil {
		return nil, fmt.Errorf("unable to query webmentions: %w", err)
	}
	defer rows.Close()

	var comments []*model.GenericComment

	for rows.Next() {
		var comment Webmention
		err := rows.Scan(&comment.Id, &comment.Source, &comment.Target, &comment.TsCreated, &comment.TsUpdated)
		if err != nil {
			return nil, fmt.Errorf("unable to scan webmention: %w", err)
		}
		genericComment := comment.ToGenericComment()
		comments = append(comments, &genericComment)
	}

	return comments, nil
}

func (s *webmentionsStore) GetGenericCommentsForPage(page string, since time.Time) ([]*model.GenericComment, error) {
	rows, err := s.db.Query("SELECT id, source, target, ts_created, ts_updated FROM webmentions WHERE deleted = false AND target = ? AND ts_updated > ?", page, since)
	if err != nil {
		return nil, fmt.Errorf("unable to query webmentions: %w", err)
	}
	defer rows.Close()

	var comments []*model.GenericComment

	for rows.Next() {
		var comment Webmention
		err := rows.Scan(&comment.Id, &comment.Source, &comment.Target, &comment.TsCreated, &comment.TsUpdated)
		if err != nil {
			return nil, fmt.Errorf("unable to scan webmention: %w", err)
		}
		genericComment := comment.ToGenericComment()
		comments = append(comments, &genericComment)
	}

	return comments, nil
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
