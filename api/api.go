package api

import (
	"encoding/json"
	"fmt"
	stdlog "log"
	"net/http"
	"os"
	"strings"
	"tiim/go-comment-api/model"

	admissioncontrol "github.com/elithrar/admission-control"
	log "github.com/go-kit/kit/log"
)

type commentServer struct {
	store *model.CommentStore
}

func NewCommentServer(store *model.CommentStore) *commentServer {
	return &commentServer{store: store}
}

func (cs *commentServer) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/comment/", cs.handleComment)
	fmt.Println("Listening on port 8080")
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	stdlog.SetOutput(log.NewStdlibAdapter(logger))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "loc", log.DefaultCaller)
	loggingMiddleware := admissioncontrol.LoggingMiddleware(logger)
	loggedRouter := loggingMiddleware(trailingSlashMiddleware(mux))

	http.ListenAndServe(":8080", loggedRouter)
}

func (cs *commentServer) handleComment(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		cs.handlePostComment(w, r)
	case http.MethodGet:
		cs.handleGetComments(w, r)
	}
}

func (cs *commentServer) handlePostComment(w http.ResponseWriter, r *http.Request) {
	var comment model.Comment

	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Error parsing json: "+err.Error(), http.StatusBadRequest)
		return
	}

	if comment.Content == "" || comment.Page == "" {
		http.Error(w, "Invalid comment", http.StatusBadRequest)
		return
	}

	if len(comment.Content) > 1024 || len(comment.Page) > 50 || len(comment.Name) > 70 || len(comment.Email) > 60 || len(comment.ReplyTo) > 40 {
		http.Error(w, "At least one field is too long", http.StatusBadRequest)
		return
	}
	err = cs.store.NewComment(&comment)
	js, err := json.Marshal(comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (cs *commentServer) handleGetComments(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")
	var commentsList []model.Comment
	if len(parts) == 1 {
		comments, err := cs.store.GetAllComments()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		commentsList = comments
	} else if len(parts) == 2 {
		comments, err := cs.store.GetCommentsForPost(parts[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		commentsList = comments
	} else {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	js, err := json.Marshal(commentsList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
