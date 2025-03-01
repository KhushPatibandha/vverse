package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/KhushPatibandha/vverse/api"
	db "github.com/KhushPatibandha/vverse/internal/DB"
)

func UploadVideo(w http.ResponseWriter, r *http.Request) {
	name, size, duration, err := Helper(w, r)
	if err != nil {
		return
	}
	query := `INSERT INTO videos (name, size, duration, created_at) VALUES (?, ?, ?, ?)`
	res, err := db.ExecCmd(query, name, size, duration, time.Now())
	if err != nil {
		err := fmt.Errorf("Failed to insert into database: %v", err)
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	vId, err := res.LastInsertId()
	if err != nil {
		err := fmt.Errorf("Failed to get last insert id: %v", err)
		log.Error(err)
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
