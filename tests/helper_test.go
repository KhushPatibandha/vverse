package tests

import (
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"

	"github.com/KhushPatibandha/vverse/consts"
	db "github.com/KhushPatibandha/vverse/internal/DB"
	"github.com/KhushPatibandha/vverse/internal/handlers"
)

var testServer *http.Server

func TestMain(m *testing.M) {
	log.Info("Setting up test database...")
	err := db.SetupDB()
	if err != nil {
		log.Error("Failed to set up DB:", err)
		os.Exit(1)
	}

	r := chi.NewRouter()
	handlers.Handler(r)

	testServer = &http.Server{
		Addr:    "localhost:" + consts.LOCALHOSTPORT,
		Handler: r,
	}

	go func() {
		log.Info("Starting test server...")
		if err := testServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Test server failed:", err)
		}
	}()

	code := m.Run()

	log.Info("Shutting down test server...")
	testServer.Close()
	_ = db.CloseDB()
	os.Remove("sqlite3.db")
	os.RemoveAll("uploads/")

	os.Exit(code)
}

func resetDatabase() {
	_, err := db.ExecCmd("DELETE FROM videos")
	if err != nil {
		panic("Failed to reset database: " + err.Error())
	}
	_, err = db.ExecCmd("DELETE FROM sqlite_sequence WHERE name='videos'")
	if err != nil {
		panic("Failed to reset autoincrement: " + err.Error())
	}
}

func addMockData() {
	query := `INSERT INTO videos (name, size, duration, created_at) VALUES 
		('e851d391-5ded-40e3-b323-2c44392421d4_1740909320', 1.68618011474609, 10.333333, ?),
		('ff9212f1-981a-479f-99d0-97bc78e83dca_1740909329', 1.1683874130249, 9.373333, ?);`
	_, err := db.ExecCmd(query, "2021-08-01 00:00:00", "2021-08-01 00:00:00")
	if err != nil {
		panic("Failed to add mock data: " + err.Error())
	}

	_ = os.MkdirAll("./uploads", os.ModePerm)
	src1 := "../test_videos/10_work.mp4"
	src2 := "../test_videos/9_work.mp4"
	videoFile1 := "./uploads/e851d391-5ded-40e3-b323-2c44392421d4_1740909320"
	videoFile2 := "./uploads/ff9212f1-981a-479f-99d0-97bc78e83dca_1740909329"

	err = copyFile(src1, videoFile1)
	if err != nil {
		panic("Failed to copy file: " + err.Error())
	}
	err = copyFile(src2, videoFile2)
	if err != nil {
		panic("Failed to copy file: " + err.Error())
	}
}

func cleanMockData() {
	os.Remove("./uploads/e851d391-5ded-40e3-b323-2c44392421d4_1740909320")
	os.Remove("./uploads/ff9212f1-981a-479f-99d0-97bc78e83dca_1740909329")
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return errors.New("failed to open source file: " + err.Error())
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return errors.New("failed to create destination file: " + err.Error())
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return errors.New("failed to copy file: " + err.Error())
	}

	return nil
}
