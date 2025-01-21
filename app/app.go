package app

import (
	"cbr/cbr"
	"cbr/renderer"
	"log/slog"
	"net/http"
	"time"
)

type App struct {
	logger     *slog.Logger
	httpClient *http.Client
	cbr        *cbr.Client
}

func NewApp(logger *slog.Logger) *App {
	return &App{
		logger:     logger,
		httpClient: http.DefaultClient,
		cbr:        cbr.NewClient(http.DefaultClient),
	}
}

func (app *App) Mount(router *http.ServeMux) {
	router.HandleFunc("/", app.Index)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (app *App) renderError(w http.ResponseWriter, err string) {
	resp := ErrorResponse{Error: err}
	renderer.JSON(w, http.StatusUnprocessableEntity, resp)
}

func (app *App) Index(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	dateStr := query.Get("date")
	date, err := time.Parse("2006/01/02", dateStr)
	if err != nil {
		app.logger.Error("Could not parse date", slog.String("error", err.Error()))
		date = time.Now()
	}
	daily, err := app.cbr.Daily(date)
	if err != nil {
		app.logger.Error("Could not get daily response", slog.String("error", err.Error()))
		app.renderError(w, "Could not get daily response")
		return
	}

	renderer.JSON(w, http.StatusOK, daily)
}
