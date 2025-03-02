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
