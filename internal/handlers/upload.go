package handlers

import (
	"fmt"
	"net/http"
)

func UploadVideo(w http.ResponseWriter, r *http.Request) {
	name, size, duration := Helper(w, r)
	fmt.Println(name, size, duration)
}
