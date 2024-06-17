package handlers

import (
	"net/http"
)

func HandleTest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test"))
}