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
	query := r.URL.Query()
	dateStr := query.Get("date")
	date, err := time.Parse("2006/01/02", dateStr)
	if err != nil {
		date = time.Now()
	}
	daily, err := cbr.Daily(date)
	if err != nil {
		renderError(w, err.Error())
		return
	}

	renderer.JSON(w, http.StatusOK, daily)
}
