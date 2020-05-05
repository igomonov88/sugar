package handlers

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	api "github.com/igomonov88/sugar/internal/fdc_api"
	"github.com/igomonov88/sugar/internal/platform/web"
	storage "github.com/igomonov88/sugar/internal/storage"
)

// Details returns info about product with given food detail
func (f *Food) Details(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Food.Details")
	defer span.End()

	var foodDetailsReq api.DetailsRequest
	if err := web.Decode(r, &foodDetailsReq); err != nil {
		return errors.Wrap(err, "")
	}

	//Check for result from given details
	result, err := storage.GetDetails(ctx, f.db, foodDetailsReq.FDCID)
	switch err {
	case sql.ErrNoRows:
		details, err := api.Details(ctx, f.apiClient, foodDetailsReq)
		if err != nil {
			errors.Wrapf(err, "")
		}
		return web.Respond(ctx, w, details, http.StatusOK)
	default:
		var foodDetailsResp api.DetailsResponse
		foodDetailsResp.Description = result.Description
		return web.Respond(ctx, w, foodDetailsResp, http.StatusOK)
	}
}
