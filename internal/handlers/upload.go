package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/KhushPatibandha/vverse/api"
)

//	@Summary		Upload a video
//	@Description	Uploads a video file and returns a video ID for further operations
//	@Tags			video
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			file	formData	file	true	"The video file to upload"
//	@Success		200		{object}	api.Response
//	@Failure		400		{object}	api.Response
//	@Failure		500		{object}	api.Response
//	@Router			/video [post]
func UploadVideo(w http.ResponseWriter, r *http.Request) {
	name, size, duration, err := Helper(w, r)
	if err != nil {
		return
	}

	vId, err := SaveInDb(name, size, duration)
	if err != nil {
		api.RequestErrorHandler(w, err)
		return
	}

	response := api.Response{
		Code:    http.StatusOK,
		Message: fmt.Sprintf("Video uploaded!! Video Id: %d, use this for further operations", vId),
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
