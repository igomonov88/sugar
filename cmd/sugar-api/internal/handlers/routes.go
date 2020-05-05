package handlers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"

	api "github.com/igomonov88/sugar/internal/fdc_api"
	"github.com/igomonov88/sugar/internal/mid"
	"github.com/igomonov88/sugar/internal/platform/auth"
	"github.com/igomonov88/sugar/internal/platform/cache"
	"github.com/igomonov88/sugar/internal/platform/web"
)

// Food represents the Food Data Central API method handler set.
type Food struct {
	db            *sqlx.DB
	apiClient     *api.Client
	cache         *cache.Cache
	authenticator *auth.Authenticator
}

// API constructs an http.Handler with all application routes defined.
func API(build string, shutdown chan os.Signal, log *log.Logger, db *sqlx.DB, fdcClient *api.Client) http.Handler {

	// Construct the web.App which holds all routes as well as common Middleware.
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics(log))

	// Register health check endpoint. This route is not authenticated.
	check := Check{
		build: build,
		db:    db,
	}

	// Construct cache configuration
	cacheConfig := cache.Config{
		DefaultDuration: 24 * time.Hour,
		Size:            1000,
	}

	c, _ := cache.New(cacheConfig)

	app.Handle("GET", "/v1/health", check.Health)

	// Register food endpoints.
	f := Food{
		apiClient: fdcClient,
		cache:     c,
		db:        db,
	}
	app.Handle("POST", "/v1/search", f.Search)
	app.Handle("POST", "/v1/details", f.Details)

	return app
}
