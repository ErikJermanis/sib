package handlers

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/ErikJermanis/sib-web/db"
	"github.com/ErikJermanis/sib-web/views/auth"
	"github.com/golang-jwt/jwt/v5"
)

func generateJWT() (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 365 * 100)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("SIB_API_JWT_SECRET")))
	return tokenString, err
}

func HandleRenderAuth(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("SIB_AUTH_TOKEN")
	if err == nil && cookie != nil {
		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(os.Getenv("SIB_API_JWT_SECRET")), nil
		})
		if err == nil {
			if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				GoToRoute(w, r, "/wishlist")
				return nil
			}
		}
	}

	return auth.Index().Render(r.Context(), w)
}

func HandleAuthenticate(w http.ResponseWriter, r *http.Request) error {
	otp := r.FormValue("otp")

	var used bool
	var expiresAt time.Time
	err := db.Db.QueryRow("SELECT used, expiresat FROM otps WHERE otp = $1", otp).Scan(&used, &expiresAt)
	if err != nil || used || time.Now().After(expiresAt) {
		return auth.PinInput().Render(r.Context(), w)
	}

	if _, err = db.Db.Exec("UPDATE otps SET used = true WHERE otp = $1", otp); err != nil {
		return err
	}

	token, err := generateJWT()
	if err != nil {
		return err
	}

	cookie := http.Cookie{
		Name: "SIB_AUTH_TOKEN",
		Value: token,
		Expires: time.Now().Add(time.Hour * 24 * 365 * 100),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)

	w.Header().Set("HX-Redirect", "/wishlist")

	return nil
}