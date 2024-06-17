package handlers

import (
	"net/http"

	"github.com/ErikJermanis/sib-web/views/test"
)

func HandleTest(w http.ResponseWriter, r *http.Request) error {
	return test.Hello().Render(r.Context(), w)
}