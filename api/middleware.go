package api

// request_logger.go
import (
	"net/http"
	"strings"
)

func trailingSlashMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/") + "/"
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
