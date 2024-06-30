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

func GenerateJWT() (string, error) {
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

	otpDetails, err := db.FetchOtpDetails(otp)
	if err != nil || otpDetails.Used || otpDetails.ExpiresAt.Before(time.Now()) {
		return auth.PinInput().Render(r.Context(), w)
	}

	if err = db.InvalidateOtp(otp); err != nil {
		return err
	}

	token, err := GenerateJWT()
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