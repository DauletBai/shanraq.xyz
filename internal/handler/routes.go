package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (h *Handler) Routes(staticAssetFS http.FileSystem) http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Compress(5))

	fileServer := http.FileServer(staticAssetFS)
	mux.Handle("/assets/*", http.StripPrefix("/assets/", fileServer))

	// Маршруты
	mux.Get("/", h.Home)
	mux.Get("/welcome", h.Welcome)
	mux.Get("/login", h.LoginGet)
	mux.Post("/login", h.LoginPost) // Строка 26
	mux.Post("/register", h.RegisterPost)
	mux.Get("/dashboard", h.Dashboard)
	// mux.Post("/logout", h.LogoutPost)

	mux.NotFound(h.notFound)
	return mux
}