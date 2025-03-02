package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"

	_ "github.com/KhushPatibandha/vverse/cmd/api/docs"
	"github.com/KhushPatibandha/vverse/consts"
	db "github.com/KhushPatibandha/vverse/internal/DB"
	"github.com/KhushPatibandha/vverse/internal/handlers"
)

// @title			VideoVerse API
// @version		1.0
// @description	These are the APIs for VideoVerse take home assignment.
// @contact.name	Khush Patibandha
// @contact.email	khush.patibandha@gmail.com
// @BasePath		/api/v1
func main() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetLevel(log.DebugLevel)

	fmt.Println(`
                                                    .__ 
___  _____  __ ___________  ______ ____      _____  |__|
\  \/ /\  \/ // __ \_  __ \/  ___// __ \     \__  \ |  |
 \   /  \   /\  ___/|  | \/\___ \\  ___/      / __ \|  |
  \_/    \_/  \___  >__|  /____  >\___  > /\ (____  /__|
                  \/           \/     \/  \/      \/    
        `)

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
