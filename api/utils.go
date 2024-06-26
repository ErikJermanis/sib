package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ResponseJson map[string]interface{}

func RespondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func Protect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			RespondWithJSON(w, http.StatusUnauthorized, ResponseJson{ "message": "Access denied." })
			return
		}
		
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			RespondWithJSON(w, http.StatusUnauthorized, ResponseJson{ "message": "Access denied." })
			return
		}

		token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SIB_API_JWT_SECRET")), nil
		})

		if err != nil {
			RespondWithJSON(w, http.StatusUnauthorized, ResponseJson{ "message": "Access denied." })
			return
		}

		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			next.ServeHTTP(w, r)
		} else {
			RespondWithJSON(w, http.StatusUnauthorized, ResponseJson{ "message": "Access denied." })
			return
		}
	})
}
