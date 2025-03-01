package handlers

import (
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"

	"github.com/KhushPatibandha/vverse/internal/middleware"
)

func Handler(r *chi.Mux) {
	r.Use(chiMiddleware.StripSlashes)

	// Route to let user upload a video file
	r.Route("/video", func(router chi.Router) {
		router.Use(middleware.Auth)
		router.Post("/", UploadVideo)
	})

	// Router to let user trim a video file
	r.Route("/trim", func(router chi.Router) {
		router.Use(middleware.Auth)
		router.Put("/", TrimVideo)
	})

	// Router to let user merge two video files
	r.Route("/merge", func(router chi.Router) {
		router.Use(middleware.Auth)
		router.Post("/", MergeVideo)
	})

	// Router to let user get a link with time-based expiry
	r.Route("/link", func(router chi.Router) {
		router.Use(middleware.Auth)
		router.Get("/", GetLink)
	})
}
