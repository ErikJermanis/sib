package handlers

import (
	"log/slog"
	"net/http"
)

type HTTPHandler func(w http.ResponseWriter, r *http.Request) error

func Make(h HTTPHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		if err := h(w, r); err != nil {
			slog.Error("error handling request", "err", err, "path", r.URL.Path)
		}
	}
}
