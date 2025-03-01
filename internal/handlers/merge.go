package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/KhushPatibandha/vverse/api"
	db "github.com/KhushPatibandha/vverse/internal/DB"
)

func MergeVideo(w http.ResponseWriter, r *http.Request) {
	v1Str := r.URL.Query().Get("v1")
	v2Str := r.URL.Query().Get("v2")
	v1Id, err := strconv.Atoi(v1Str)
	if err != nil {
		err := errors.New("Invalid video ID: " + err.Error())
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}
	v2Id, err := strconv.Atoi(v2Str)
	if err != nil {
		err := errors.New("Invalid video ID: " + err.Error())
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	var v1Name string
	var v2Name string
	query := `SELECT name FROM videos WHERE id = ?`
	err = db.QueryRow(query, v1Id).Scan(&v1Name)
	if err != nil {
		err := errors.New("Failed to get first video name: " + err.Error())
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}
	err = db.QueryRow(query, v2Id).Scan(&v2Name)
	if err != nil {
		err := errors.New("Failed to get second video name: " + err.Error())
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	v1Path := "./uploads/" + v1Name
	v2Path := "./uploads/" + v2Name

	tempFile, err := os.Create("concat.txt")
	if err != nil {
		err := errors.New("Failed to create temporary file: " + err.Error())
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}
	defer os.Remove("concat.txt")

	fileContent := fmt.Sprintf("file '%s'\nfile '%s'\n", v1Path, v2Path)
	_, err = tempFile.WriteString(fileContent)
	if err != nil {
		err := errors.New("Failed to write to temporary file: " + err.Error())
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}
	tempFile.Close()

	tempOutputPath := "./uploads/temp_" + v1Name + "_" + v2Name + ".mp4"

	//  ffmpeg -f concat -safe 0 -i concat.txt -c:v libx264 -c:a aac output.mp4
	cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", tempFile.Name(), "-c:v", "libx264", "-c:a", "aac", tempOutputPath)

	var stdoutStderr bytes.Buffer
	cmd.Stdout = &stdoutStderr
	cmd.Stderr = &stdoutStderr

	err = cmd.Run()
	if err != nil {
		log.Errorf("FFmpeg command failed: %v", err)
		log.Errorf("FFmpeg output: %s", stdoutStderr.String())

		err := errors.New("error merging files: " + err.Error())
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}
	defer os.Remove(tempOutputPath)

	mergedFile, err := os.Open(tempOutputPath)
	if err != nil {
		err := errors.New("Failed to open merged file: " + err.Error())
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}
	defer mergedFile.Close()

	fInfo, err := mergedFile.Stat()
	if err != nil {
		err := errors.New("Failed to get file info: " + err.Error())
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	sizeMB, err := checkSize(fInfo.Size())
	if err != nil {
		api.RequestErrorHandler(w, err)
		return
	}

	if !isVideoFile(tempOutputPath) {
		log.Error(err)
		api.RequestErrorHandler(w, errors.New("Uploaded file is not a video"))
		return
	}

	duration, err := getDuration(tempOutputPath)
	if err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	if err := checkDuration(duration); err != nil {
		api.RequestErrorHandler(w, err)
		return
	}

	filename, err := uploadInStorange(tempOutputPath)
	if err != nil {
		api.RequestErrorHandler(w, err)
		return
	}

	vId, err := SaveInDb(filename, sizeMB, duration)
	if err != nil {
		api.RequestErrorHandler(w, err)
		return
	}

	response := api.Response{
		Code:    http.StatusOK,
		Message: fmt.Sprintf("Videos merged!! New merged video Id: %d, use this for further operations", vId),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Code)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		panic(err)
	}
}
