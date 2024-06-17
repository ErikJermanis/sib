package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/ErikJermanis/sib-web/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err);
	}
	router := chi.NewRouter()

	router.Get("/test", handlers.Make(handlers.HandleTest))
	
	port := os.Getenv("HTTP_PORT")
	slog.Info(fmt.Sprintf("Server is running on port %s", port))
	http.ListenAndServe(port, router)
}
