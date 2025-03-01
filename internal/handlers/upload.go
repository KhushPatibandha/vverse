package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KhushPatibandha/vverse/api"
)

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
		panic(err)
	}
}
