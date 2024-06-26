package api

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func CheckIfAuthenticated(w http.ResponseWriter, r *http.Request) {
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
		RespondWithJSON(w, http.StatusOK, ResponseJson{ "authenticated": true })
	} else {
		RespondWithJSON(w, http.StatusUnauthorized, ResponseJson{ "message": "Access denied." })
		return
	}
}
