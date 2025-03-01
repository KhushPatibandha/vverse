package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
	log "github.com/sirupsen/logrus"

	"github.com/KhushPatibandha/vverse/api"
	db "github.com/KhushPatibandha/vverse/internal/DB"
	"github.com/KhushPatibandha/vverse/internal/middleware"
)

func Handler(r *chi.Mux) {
	r.Use(chiMiddleware.StripSlashes)

	// Route to let user upload a video file
	r.Route("/api/v1/video", func(router chi.Router) {
		router.Use(middleware.Auth)
		router.Post("/", UploadVideo)
	})

	// Router to let user trim a video file
	// r.Route("/api/v1/trim", func(router chi.Router) {
	// 	router.Use(middleware.Auth)
	// 	router.Put("/", TrimVideo)
	// })

	// Router to let user merge two video files
	r.Route("/api/v1/merge", func(router chi.Router) {
		router.Use(middleware.Auth)
		router.Post("/", MergeVideo)
	})

	// Router to let user get a link with time-based expiry
	r.Route("/api/v1/link", func(router chi.Router) {
		router.Use(middleware.Auth)
		router.Get("/", GetLink)
	})

	// Router temp link to uploaded vids
	r.Route("/api/v1/uploads", func(router chi.Router) {
		router.Get("/{link}", RedirectVid)
	})
}

func Helper(w http.ResponseWriter, r *http.Request) (string, float64, float64, error) {
	// create temp to either rename or delete later
	file, err := os.CreateTemp("", "upload_*")
	if err != nil {
		log.Error("Failed to create temp file:", err)
		api.InternalErrorHandler(w)
		return "", 0, 0, err
	}
	tempFilePath := file.Name()
	defer os.Remove(tempFilePath)

	// copy data to temp
	size, err := io.Copy(file, r.Body)
	if err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, errors.New("Failed to save file"))
		return "", 0, 0, err
	}

	// calc size in MB and validate it
	// assumption: 25MB is the max size with no min size
	sizeMB, err := checkSize(size)
	if err != nil {
		api.RequestErrorHandler(w, err)
		return "", 0, 0, err
	}

	// check if the uploaded file is a video or not
	if !isVideoFile(tempFilePath) {
		log.Error(err)
		api.RequestErrorHandler(w, errors.New("Uploaded file is not a video"))
		return "", 0, 0, err
	}

	// if is video, get the duration in seconds and validate it
	duration, err := getDuration(tempFilePath)
	if err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return "", 0, 0, err
	}

	// assumption: 5 sec is the min and 25 sec is the max
	if err := checkDuration(duration); err != nil {
		api.RequestErrorHandler(w, err)
		return "", 0, 0, err
	}

	// if everything works fine, save the file and return the data
	filename, err := uploadInStorange(tempFilePath)
	if err != nil {
		api.RequestErrorHandler(w, err)
		return "", 0, 0, err
	}

	return filename, sizeMB, duration, nil
}

func SaveInDb(name string, size float64, duration float64) (int64, error) {
	query := `INSERT INTO videos (name, size, duration, created_at) VALUES (?, ?, ?, ?)`
	res, err := db.ExecCmd(query, name, size, duration, time.Now())
	if err != nil {
		err := fmt.Errorf("Failed to insert into database: %v", err)
		log.Error(err)
		return 0, err
	}

	vId, err := res.LastInsertId()
	if err != nil {
		err := fmt.Errorf("Failed to get last insert id: %v", err)
		log.Error(err)
		return 0, err
	}
	return vId, nil
}

func isVideoFile(filePath string) bool {
	file, _ := os.Open(filePath)
	defer file.Close()

	head := make([]byte, 261)
	_, err := file.Read(head)
	if err != nil {
		return false
	}
	if filetype.IsVideo(head) {
		return true
	}
	return false
}

func getDuration(filePath string) (float64, error) {
	// ffprobe -v error -select_streams v:0 -show_entries stream=duration -of default=noprint_wrappers=1:nokey=1 input.mp4
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=duration", "-of", "default=noprint_wrappers=1:nokey=1", filePath)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return 0, errors.New("ffprobe error : " + err.Error() + " output: " + out.String())
	}

	durationStr := strings.TrimSpace(out.String())
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, errors.New("Failed to parse duration")
	}
	return duration, nil
}

func checkSize(size int64) (float64, error) {
	sizeMB := float64(size) / (1024 * 1024)
	if sizeMB > 25 {
		err := errors.New("File size should be less than 25MB, got: " + strconv.FormatFloat(sizeMB, 'f', -1, 64))
		log.Error(err)
		return 0, err
	}
	return sizeMB, nil
}

func checkDuration(duration float64) error {
	if duration < 5 || duration > 25 {
		err := errors.New("Video duration should be between 5 and 25 seconds, got: " + strconv.FormatFloat(duration, 'f', -1, 64))
		log.Error(err)
		return err
	}
	return nil
}

func uploadInStorange(tempFilePath string) (string, error) {
	uploadPath := "./uploads/"
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		log.Error("Failed to create uploads directory:", err)
		return "", err
	}

	filename := fmt.Sprintf("%s_%d", uuid.New().String(), time.Now().Unix())
	finalFilePath := filepath.Join(uploadPath, filename)

	if err := os.Rename(tempFilePath, finalFilePath); err != nil {
		log.Error("Failed to move file:", err)
		return "", err
	}
	return filename, nil
}
