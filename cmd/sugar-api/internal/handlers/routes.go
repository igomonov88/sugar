package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"

	apiClient "github.com/igomonov88/sugar/internal/food_data_center_api"
	"github.com/igomonov88/sugar/internal/mid"
	"github.com/igomonov88/sugar/internal/platform/web"
)

// API constructs an http.Handler with all application routes defined.
func API(build string, shutdown chan os.Signal, log *log.Logger, db *sqlx.DB, fdcClient *apiClient.Client) http.Handler {

	// Construct the web.App which holds all routes as well as common Middleware.
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics(log))

	// Register health check endpoint. This route is not authenticated.
	check := Check{
		build: build,
		db:    db,
	}

	app.Handle("GET", "/v1/health", check.Health)

	// Register food endpoints.
	f := Food{
		apiClient: fdcClient,
		db:        db,
	}
	app.Handle("POST", "/v1/search", f.Search)
	app.Handle("POST","/v1/details", f.Details)

	return app
}
