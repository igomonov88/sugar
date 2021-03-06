package handlers

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"go.opencensus.io/trace"

	"github.com/igomonov88/sugar/internal/platform/database"
	"github.com/igomonov88/sugar/internal/platform/web"
)

// Check provides support for orchestration health checks.
type Check struct {
	build string
	db    *sqlx.DB
}

// Health validates the service is healthy and ready to accept requests.
func (c *Check) Health(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Check.Health")
	defer span.End()

	health := struct {
		Version string `json:"version"`
		Status  string `json:"status"`
	}{
		Version: c.build,
	}

	// Check if the database is ready.
	if err := database.StatusCheck(ctx, c.db); err != nil {

		// If the database is not ready we will tell the client and use a 500
		// status. Do not respond by just returning an error because further up in
		// the call stack will interpret that as an unhandled error.
		health.Status = "db not ready"
		return web.Respond(ctx, w, health, http.StatusInternalServerError)
	}

	health.Status = "ok"
	return web.Respond(ctx, w, health, http.StatusOK)
}
