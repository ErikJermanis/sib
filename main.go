package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/ErikJermanis/sib-web/db"
	"github.com/ErikJermanis/sib-web/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err);
	}

	db.Initialize(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_TABLE"))

	router := chi.NewRouter()

	router.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// TODO: use grouping to assign middleware only to some routes
	router.Get("/authenticate", handlers.Make(handlers.HandleRenderAuth))
	router.Post("/authenticate", handlers.Make(handlers.HandleAuthenticate))
	router.Get("/wishlist", handlers.Protect(handlers.Make(handlers.HandleGetWishes)))
	router.Post("/wishlist", handlers.Protect(handlers.Make(handlers.HandleCreateWish)))
	router.Put("/wishlist/{id}", handlers.Protect(handlers.Make(handlers.HandleSelectWish)))
	router.Patch("/wishlist/{id}", handlers.Protect(handlers.Make(handlers.HandleDeselectWish)))
	router.Post("/wishlist/{id}", handlers.Protect(handlers.Make(handlers.HandleCompleteWish)))
	router.Delete("/wishlist/{id}", handlers.Protect(handlers.Make(handlers.HandleDeleteWish)))
	router.Post("/wishlist/reset/{id}", handlers.Protect(handlers.Make(handlers.HandleResetWish)))

	
	port := os.Getenv("HTTP_PORT")
	slog.Info(fmt.Sprintf("Server is running on port %s", port))
	http.ListenAndServe(port, router)
}
