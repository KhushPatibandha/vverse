package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/KhushPatibandha/vverse/api"
	"github.com/KhushPatibandha/vverse/consts"
	db "github.com/KhushPatibandha/vverse/internal/DB"
)

// @Summary		Generate a temporary link for a video
// @Description	Generates a time-limited access link for a video using its ID
// @Tags			video
// @Produce		json
// @Param			id	query	int	true	"Video ID"
// @Security		ApiKeyAuth
// @Success		200	{object}	api.Response
// @Failure		400	{object}	api.Response
// @Failure		500	{object}	api.Response
// @Router			/link [get]
func GetLink(w http.ResponseWriter, r *http.Request) {
	vIdStr := r.URL.Query().Get("id")
	vId, err := strconv.Atoi(vIdStr)
	if err != nil {
		err := errors.New("Invalid video ID: " + err.Error())
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	var vName string
	query := `SELECT name FROM videos WHERE id = ?`
	err = db.QueryRow(query, vId).Scan(&vName)
	if err != nil {
		err := fmt.Errorf("Failed to get video name: %v", err)
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	link := uuid.New().String()
	expiry := time.Now().Add(2 * time.Minute)

	query = `INSERT INTO links (video_id, link, expiry) VALUES (?, ?, ?)`
	_, err = db.ExecCmd(query, vId, link, expiry)
	if err != nil {
		err := fmt.Errorf("Failed to insert into database: %v", err)
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	response := api.Response{
		Code:    http.StatusOK,
		Message: fmt.Sprintf("Link generated!! Use this link to access the video: `http://localhost:"+consts.LOCALHOSTPORT+"/api/v1/uploads/%s`", link),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Code)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}
