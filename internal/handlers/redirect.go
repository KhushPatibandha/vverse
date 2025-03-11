package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"

	"github.com/KhushPatibandha/vverse/api"
	db "github.com/KhushPatibandha/vverse/internal/DB"
)

//	@Summary		Redirect to the uploaded video
//	@Description	Given a valid temporary link, redirects to the video file
//	@Tags			video
//	@Produce		octet-stream
//	@Param			link	path	string	true	"Temporary video link"
//	@Success		200		"Video file served"
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/uploads/{link} [get]
func RedirectVid(w http.ResponseWriter, r *http.Request) {
	link := chi.URLParam(r, "link")
	if link == "" {
		err := errors.New("No link provided")
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	var vId int
	var expiry time.Time
	query := "SELECT video_id, expiry FROM links WHERE link = ?"
	err := db.QueryRow(query, link).Scan(&vId, &expiry)
	if err != nil {
		err = errors.New("Failed to query DB for link: " + err.Error())
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	if time.Now().After(expiry) {
		err := errors.New("Link expired")
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	var vName string
	query = "SELECT name FROM videos WHERE id = ?"
	err = db.QueryRow(query, vId).Scan(&vName)
	if err != nil {
		err = errors.New("Failed to query DB for video name")
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	vPath := "./uploads/" + vName
	http.ServeFile(w, r, vPath)
}
