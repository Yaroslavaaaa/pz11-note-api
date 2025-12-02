package httpx

import (
	"net/http"

	"example.com/notes-api/internal/http/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(h *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/notes", func(r chi.Router) {
			r.Post("/", h.CreateNote)
			r.Get("/", h.GetAllNotes)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.GetNote)
				r.Patch("/", h.PatchNote)
				r.Delete("/", h.DeleteNote)

			})
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	return r
}
