package handlers

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/KhushPatibandha/vverse/api"
	db "github.com/KhushPatibandha/vverse/internal/DB"
)

func TrimVideo(w http.ResponseWriter, r *http.Request) {
	vIdStr := r.URL.Query().Get("id")
	vId, err := strconv.Atoi(vIdStr)
	if err != nil {
		err := errors.New("Invalid video ID: " + err.Error())
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}
	startStr := r.URL.Query().Get("s")
	endStr := r.URL.Query().Get("e")
	start, err := strconv.Atoi(startStr)
	if err != nil {
		err := errors.New("Invalid start time: " + err.Error())
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}
	end, err := strconv.Atoi(endStr)
	if err != nil {
		err := errors.New("Invalid end time: " + err.Error())
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	var vName string
	query := `SELECT name FROM videos WHERE id = ?`
	err = db.QueryRow(query, vId).Scan(&vName)
	if err != nil {
		err := errors.New("Failed to get first video name: " + err.Error())
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	filePath := "./uploads/" + vName
	duration, err := getDuration(filePath)
	if err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	if start < 0 || end < 0 || start > end || end > int(math.Ceil(duration)) {
		err := errors.New("Invalid start or end time")
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	tempFilePath := "./uploads/trimmed_" + vName + ".mp4"
	defer os.Remove(tempFilePath)

	// ffmpeg -ss 00:01:00 -to 00:02:00 -i input.mp4 -c copy output.mp4
	cmd := exec.Command("ffmpeg", "-ss", "00:00:"+startStr, "-to", "00:00:"+endStr, "-i", filePath, "-c", "copy", tempFilePath)

	err = cmd.Run()
	if err != nil {
		err := errors.New("Failed to trim video: " + err.Error())
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	if err := os.Rename(tempFilePath, filePath); err != nil {
		err := errors.New("Failed to rename trimmed video: " + err.Error())
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	// update the duration in db
	newDuration := end - start
	query = `UPDATE videos SET duration = ? WHERE id = ?`
	_, err = db.ExecCmd(query, newDuration, vId)
	if err != nil {
		err := errors.New("Failed to update duration in db: " + err.Error())
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	response := api.Response{
		Code:    http.StatusOK,
		Message: "Video trimmed successfullt, you can get the trimmed video with the same ID",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Code)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		panic(err)
	}
}
