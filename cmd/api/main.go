package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"

	"github.com/KhushPatibandha/vverse/consts"
	db "github.com/KhushPatibandha/vverse/internal/DB"
	"github.com/KhushPatibandha/vverse/internal/handlers"
)

func main() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetLevel(log.DebugLevel)

	log.Info("Setting up the database...")
	err := db.SetupDB()
	if err != nil {
		log.Error(err)
		return
	}
	defer func() {
		err := db.CloseDB()
		if err != nil {
			log.Error("Failed to close DB: ", err)
		}
	}()
	log.Info("Database setup complete.")

	log.SetReportCaller(true)
	var r *chi.Mux = chi.NewRouter()
	handlers.Handler(r)
	fmt.Println("Starting server on port " + consts.LOCALHOSTPORT)
	err = http.ListenAndServe("localhost:"+consts.LOCALHOSTPORT, r)
	if err != nil {
		log.Error(err)
	}
}
