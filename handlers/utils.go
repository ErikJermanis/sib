package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
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

func GoToRoute(w http.ResponseWriter, r *http.Request, route string) {
	if r.Header.Get("Hx-Request") == "true" {
		w.Header().Set("HX-Redirect", route)
	} else {
		http.Redirect(w, r, route, http.StatusSeeOther)
	}
}

func Protect(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("SIB_AUTH_TOKEN")

		if err != nil {
			if !errors.Is(err, http.ErrNoCookie) {
				slog.Error("error retrieving a cookie", "err", err)
			}
			GoToRoute(w, r, "/authenticate")
			return
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(os.Getenv("SIB_API_JWT_SECRET")), nil
		})

		if err != nil {
			GoToRoute(w, r, "/authenticate")
			return
		}

		if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
			GoToRoute(w, r, "/authenticate")
			return
		}

		next.ServeHTTP(w, r)
	}
}
