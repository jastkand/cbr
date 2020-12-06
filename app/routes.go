package app

import (
	"cbr/cbr"
	"cbr/renderer"
	"net/http"
	"time"
)

func Mount(router *http.ServeMux) {
	router.HandleFunc("/", Index)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func renderError(w http.ResponseWriter, err string) {
	resp := ErrorResponse{Error: err}
	renderer.JSON(w, http.StatusUnprocessableEntity, resp)
}

func Index(w http.ResponseWriter, r *http.Request) {
	daily, err := cbr.Daily(time.Now())
	if err != nil {
		renderError(w, err.Error())
		return
	}

	renderer.JSON(w, http.StatusOK, daily)
}
