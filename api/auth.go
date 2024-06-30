package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ErikJermanis/sib-web/db"
	"github.com/ErikJermanis/sib-web/handlers"
	"github.com/golang-jwt/jwt/v5"
)

type AuthenticateBody struct {
	Otp string `json:"otp"`
}

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

func Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestBody AuthenticateBody
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		RespondWithJSON(w, http.StatusBadRequest, ResponseJson{ "message": err.Error() })
		return
	}

	if requestBody.Otp == "" {
		RespondWithJSON(w, http.StatusBadRequest, ResponseJson{ "message": "OTP is required." })
		return
	}

	otpDetails, err := db.FetchOtpDetails(requestBody.Otp)
	if err != nil || otpDetails.Used || otpDetails.ExpiresAt.Before(time.Now()) {
		RespondWithJSON(w, http.StatusUnauthorized, ResponseJson{ "message": "Access denied." })
		return
	}

	if err = db.InvalidateOtp(requestBody.Otp); err != nil {
		RespondWithJSON(w, http.StatusInternalServerError, ResponseJson{ "message": err.Error() })
		return
	}

	token, err := handlers.GenerateJWT()
	if err != nil {
		RespondWithJSON(w, http.StatusInternalServerError, ResponseJson{ "message": err.Error() })
		return
	}

	RespondWithJSON(w, http.StatusOK, ResponseJson{ "token": token })
}
