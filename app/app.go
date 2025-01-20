package app

import (
	"cbr/cbr"
	"cbr/renderer"
	"log/slog"
	"net/http"
	"time"
)

type App struct {
	logger *slog.Logger
}

func NewApp(logger *slog.Logger) *App {
	return &App{logger: logger}
}

func (app *App) Mount(router *http.ServeMux) {
	router.HandleFunc("/", app.Index)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (app *App) renderError(w http.ResponseWriter, err string) {
	app.logger.Error(err)
	resp := ErrorResponse{Error: err}
	renderer.JSON(w, http.StatusUnprocessableEntity, resp)
}

func (app *App) Index(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	dateStr := query.Get("date")
	date, err := time.Parse("2006/01/02", dateStr)
	if err != nil {
		date = time.Now()
	}
	daily, err := cbr.Daily(date)
	if err != nil {
		app.renderError(w, err.Error())
		return
	}

	renderer.JSON(w, http.StatusOK, daily)
}
