package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/ErikJermanis/sib-web/api"
	"github.com/ErikJermanis/sib-web/db"
	"github.com/ErikJermanis/sib-web/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("error loading .env", err);
	}

	db.Initialize(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_TABLE"))

	router := chi.NewRouter()

	router.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	router.Group(func(router chi.Router) {
		router.Get("/authenticate", handlers.Make(handlers.HandleRenderAuth))
		router.Post("/authenticate", handlers.Make(handlers.HandleAuthenticate))
	})
	router.Group(func(router chi.Router) {
		router.Use(handlers.Protect)
		router.Get("/wishlist", handlers.Make(handlers.HandleGetWishes))
		router.Post("/wishlist", handlers.Make(handlers.HandleCreateWish))
		router.Put("/wishlist/{id}", handlers.Make(handlers.HandleSelectWish))
		router.Patch("/wishlist/{id}", handlers.Make(handlers.HandleDeselectWish))
		router.Post("/wishlist/{id}", handlers.Make(handlers.HandleCompleteWish))
		router.Delete("/wishlist/{id}", handlers.Make(handlers.HandleDeleteWish))
		router.Post("/wishlist/reset/{id}", handlers.Make(handlers.HandleResetWish))
	})

	apiRouter := chi.NewRouter()

	apiRouter.Use(api.CORS)
	apiRouter.Group(func(router chi.Router) {
		router.Get("/authenticate", api.CheckIfAuthenticated)
		router.Post("/authenticate", api.Authenticate)
	})
	apiRouter.Group(func(router chi.Router) {
		router.Use(api.Protect)
		router.Get("/wishlist", api.GetWishes)
		router.Get("/wishlist/{id}", api.GetWish)
		router.Post("/wishlist", api.CreateWish)
		router.Put("/wishlist/{id}", api.UpdateWish)
		router.Delete("/wishlist/{id}", api.DeleteWish)
		router.Get("/items", api.GetItems);
		router.Get("/items/{id}", api.GetItem);
		router.Post("/items", api.CreateItem);
		router.Put("/items/{id}", api.UpdateItem);
		router.Delete("/items/{id}", api.DeleteItem);
		router.Delete("/items", api.DeleteCompletedItems);
	})

	router.Mount("/api", apiRouter)
	
	port := os.Getenv("HTTP_PORT")
	slog.Info(fmt.Sprintf("Server is running on port %s", port))
	http.ListenAndServe(port, router)
}
